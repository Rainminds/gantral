package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Rainminds/gantral/core/activities"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
	"github.com/Rainminds/gantral/core/workflows"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	hostPort := os.Getenv("TEMPORAL_HOST_PORT")
	if hostPort == "" {
		hostPort = "localhost:7233"
	}
	taskQueue := os.Getenv("TASK_QUEUE")
	if taskQueue == "" {
		taskQueue = "gantral-core"
	}

	// 1. Connect to Temporal
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer c.Close()

	// 2. Start Workflow
	workflowID := fmt.Sprintf("demo-inst-%d", time.Now().Unix())
	instanceID := workflowID // Typically generated inside, but we can assume correlation
	// actually, we set it later.

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQueue,
	}

	input := workflows.WorkflowInput{
		WorkflowID: workflowID,
		TriggerContext: map[string]interface{}{
			"amount":   1000000,
			"currency": "USD",
			"vendor":   "Acme Corp",
		},
		Policy: policy.Policy{
			ID:                     "policy-demo-high-materiality",
			Materiality:            policy.MaterialityHigh, // Forces HITL
			RequiresHumanApproval:  true,
			ApprovalTimeoutSeconds: 60,
		},
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.GantralExecutionWorkflow, input)
	if err != nil {
		log.Fatalf("Unable to execute workflow: %v", err)
	}
	fmt.Printf("🚀 Workflow Started: %s (RunID: %s)\n", we.GetID(), we.GetRunID())

	// 3. Poll for Human Approval Requirement (simulating UI waiting)
	fmt.Println("⏳ Waiting for workflow to reach 'WAITING_FOR_HUMAN' state...")
	time.Sleep(3 * time.Second) // Give it a moment to run logic and hit HITL

	// 4. Send Approval Signal
	fmt.Println("👤 Simulate Human Approval (Auditor: Alice)...")
	// Note: We need the InstanceID to Signal correctly?
	// The workflow validates: if decisionInput.InstanceID == inst.ID.
	// But we don't know inst.ID yet (it's generated inside).
	// However, the Workflow logic says:
	// "If decisionInput.InstanceID == empty, set to inst.ID" (Line 144 in execution.go)
	// OR "Invalid signal: Log and continue" (Line 138).
	// Wait, Check Line 128:
	// if decisionInput.ActorID == "SYSTEM" break ...
	// if decisionInput.InstanceID == inst.ID break ...
	// If I send empty InstanceID, it will NOT break loop!
	// This is a Catch-22 if I don't know the ID.
	// But `activities.PersistInstance` generates ID.
	// Usually `PersistInstanceInput.InstanceID` comes from `workflow.GetInfo(ctx).WorkflowExecution.ID` (Line 72).
	// So InstanceID == WorkflowID!
	// instanceID = workflowID // Already set

	signalPayload := activities.RecordDecisionInput{
		InstanceID:    instanceID, // Matches WorkflowID as per Line 72 of execution.go
		DecisionType:  engine.DecisionApprove,
		ActorID:       "auditor-alice@example.com",
		Role:          "AUDITOR",
		Justification: "Demo Approval",
	}

	err = c.SignalWorkflow(context.Background(), workflowID, "", workflows.SignalHumanDecision, signalPayload)
	if err != nil {
		log.Fatalf("Failed to signal workflow: %v", err)
	}
	fmt.Println("✅ Signal Sent: APPROVED")

	// 5. Wait for Result
	fmt.Println("⏳ Waiting for workflow completion...")
	var result workflows.WorkflowResult
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	fmt.Printf("\n🎉 Workflow Completed!\n")
	fmt.Printf("   InstanceID: %s\n", result.InstanceID)
	fmt.Printf("   Final State: %s\n", result.FinalState)

	// 6. Find Generated Artifact
	fmt.Println("\n🔍 Locating Artifact for Verification...")
	artifactPath := findLatestArtifact(result.InstanceID)
	if artifactPath != "" {
		fmt.Printf("   Artifact Found: %s\n", artifactPath)
		fmt.Printf("\n📋 Run Verification Command:\n")
		fmt.Printf("   go run cmd/gantral-verify/main.go file %s\n", artifactPath)
	} else {
		fmt.Println("   ⚠️ No artifact found in ./gantral_artifacts")
	}
}

func findLatestArtifact(instanceID string) string {
	files, err := os.ReadDir("./gantral_artifacts")
	if err != nil {
		return ""
	}

	var match string
	var lastTime time.Time

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		path := filepath.Join("gantral_artifacts", f.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var art models.CommitmentArtifact
		if err := json.Unmarshal(data, &art); err != nil {
			continue
		}

		if art.InstanceID == instanceID {
			// Parse timestamp
			ts, err := time.Parse(time.RFC3339, art.Timestamp)
			if err != nil {
				continue
			}
			if ts.After(lastTime) {
				lastTime = ts
				match = path
			}
		}
	}
	return match
}
