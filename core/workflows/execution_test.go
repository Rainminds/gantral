package workflows

import (
	"errors"
	"testing"
	"time"

	"github.com/Rainminds/gantral/core/activities"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_HappyPath_AutoApprove() {
	// Input with Low policy -> Auto Approve
	input := WorkflowInput{
		WorkflowID:     "wf-123",
		TriggerContext: map[string]interface{}{"foo": "bar"},
		Policy: policy.Policy{
			ID:                    "pol-low",
			Materiality:           policy.MaterialityLow,
			RequiresHumanApproval: false,
		},
	}

	// Mock Activity: PersistInstance
	var a *activities.ExecutionActivities
	s.env.OnActivity(
		a.PersistInstance,
		mock.Anything,
		mock.MatchedBy(func(arg activities.PersistInstanceInput) bool {
			return arg.WorkflowID == "wf-123" && arg.InitialState == engine.StateRunning
		}),
	).Return(&engine.Instance{
		ID:    "inst-mock-1",
		State: engine.StateRunning,
	}, nil)

	s.env.ExecuteWorkflow(GantralExecutionWorkflow, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result WorkflowResult
	s.env.GetWorkflowResult(&result)
	s.Equal("inst-mock-1", result.InstanceID)
	s.Equal(engine.StateRunning, result.FinalState)
}

func (s *UnitTestSuite) Test_HITL_Flow_Approved() {
	// Input with High policy -> HITL
	input := WorkflowInput{
		WorkflowID: "wf-hitl",
		Policy: policy.Policy{
			ID:          "pol-high",
			Materiality: policy.MaterialityHigh,
		},
	}

	var a *activities.ExecutionActivities
	// Mock Activity: PersistInstance (Initial State: WaitingForHuman)
	s.env.OnActivity(
		a.PersistInstance,
		mock.Anything,
		mock.MatchedBy(func(arg activities.PersistInstanceInput) bool {
			return arg.InitialState == engine.StateWaitingForHuman
		}),
	).Return(&engine.Instance{
		ID:    "inst-hitl-1",
		State: engine.StateWaitingForHuman,
	}, nil)

	// Mock Activity: RecordDecision
	s.env.OnActivity(
		a.RecordDecision,
		mock.Anything,
		mock.MatchedBy(func(arg activities.RecordDecisionInput) bool {
			return arg.InstanceID == "inst-hitl-1" && arg.DecisionType == engine.DecisionApprove
		}),
	).Return(nil) // Success

	// Register Delayed Signal
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalHumanDecision, activities.RecordDecisionInput{
			InstanceID:    "inst-hitl-1",
			DecisionType:  engine.DecisionApprove,
			ActorID:       "human-1",
			Justification: "Looks good",
		})
	}, 1*time.Second)

	s.env.ExecuteWorkflow(GantralExecutionWorkflow, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result WorkflowResult
	s.env.GetWorkflowResult(&result)
	s.Equal("inst-hitl-1", result.InstanceID)
	// Simplified workflow logic updates local struct but returns the "FinalState"
	// based on the switch case if we implemented strict state tracking.
	// In the implementation I assumed StateApproved.
	s.Equal(engine.StateApproved, result.FinalState)
}

func (s *UnitTestSuite) Test_ActivityFailure_Retry() {
	// Test that activity failure bubbles up (or is retried, but in test env we see error if max retries hit)
	// We can just verify proper error handling.
	input := WorkflowInput{WorkflowID: "wf-fail"}

	var a *activities.ExecutionActivities
	s.env.OnActivity(a.PersistInstance, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	s.env.ExecuteWorkflow(GantralExecutionWorkflow, input)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError()) // Should fail eventually
}

func (s *UnitTestSuite) Test_HITL_Timeout() {
	// Verify 24h timeout triggers Rejection
	input := WorkflowInput{
		WorkflowID: "wf-timeout",
		Policy: policy.Policy{
			ID:          "pol-high",
			Materiality: policy.MaterialityHigh,
		},
	}

	var a *activities.ExecutionActivities
	// 1. Initial Persist (Waiting for Human)
	s.env.OnActivity(
		a.PersistInstance,
		mock.Anything,
		mock.Anything,
	).Return(&engine.Instance{
		ID:    "inst-timeout-1",
		State: engine.StateWaitingForHuman,
	}, nil)

	// 2. Record Decision (SYSTEM REJECT) - expected after timeout
	s.env.OnActivity(
		a.RecordDecision,
		mock.Anything,
		mock.MatchedBy(func(arg activities.RecordDecisionInput) bool {
			return arg.InstanceID == "inst-timeout-1" &&
				arg.DecisionType == engine.DecisionReject &&
				arg.ActorID == "SYSTEM"
		}),
	).Return(nil)

	// Note: No signal sent. We rely on time skipping.
	// Temporal TestEnv automatically skips time if workflow is blocked on timer.

	s.env.ExecuteWorkflow(GantralExecutionWorkflow, input)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result WorkflowResult
	s.env.GetWorkflowResult(&result)
	s.Equal("inst-timeout-1", result.InstanceID)
	s.Equal(engine.StateRejected, result.FinalState)
}

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
