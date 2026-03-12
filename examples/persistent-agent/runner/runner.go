package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	sdk "github.com/Rainminds/gantral/sdk/go" // Official Gantral Go SDK
)

var (
	gantralURL = getEnv("GANTRAL_URL", "http://localhost:8080")
	pollInt    = 3 * time.Second
	runnerID   = getEnv("RUNNER_ID", "runner-001")
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// seenStates keeps track of instances we've already reacted to
var seenStates = make(map[string]sdk.State)

func main() {
	log.Printf("Starting Gantral Reference Runner (Go implementation)...")
	log.Printf("Pointing to Core at: %s", gantralURL)

	client := sdk.NewClient(gantralURL)
	ctx := context.Background()

	for {
		instances, err := getInstances()
		if err != nil {
			log.Printf("Error polling Gantral API: %v", err)
			time.Sleep(pollInt)
			continue
		}

		for _, inst := range instances {
			lastState, seen := seenStates[inst.ID]

			if !seen || inst.State != lastState {
				seenStates[inst.ID] = inst.State

				switch inst.State {
				case sdk.StateWaitingForHuman:
					log.Printf("Instance %s is WAITING_FOR_HUMAN. Launching Agent to hibernate...", inst.ID)
					go runAgent(inst.ID, string(inst.State), inst.TriggerContext)

					// Automatically approve via SDK for checkpoint proof (Example purposes)
					log.Printf("Runner auto-approving %s via SDK to demonstrate checkpoint...", inst.ID)
					_, err := client.RecordDecision(ctx, inst.ID, sdk.DecisionApprove, runnerID, "Auto-approved by Go Runner via SDK")
					if err != nil {
						log.Printf("Failed to record decision via SDK: %v", err)
					} else {
						log.Printf("Successfully recorded decision for %s", inst.ID)
					}
				case sdk.StateApproved, sdk.StateRunning, sdk.StatePending:
					log.Printf("Instance %s is %s. Launching Agent to resume/start...", inst.ID, inst.State)
					go runAgent(inst.ID, string(inst.State), inst.TriggerContext)
				case sdk.StateCompleted:
					log.Printf("Instance %s is COMPLETED. Skipping.", inst.ID)
				}
			}
		}

		time.Sleep(pollInt)
	}
}

// getInstances fetches instances directly from the API (Long polling or list not in SDK yet)
func getInstances() ([]sdk.Instance, error) {
	resp, err := http.Get(fmt.Sprintf("%s/instances", gantralURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Instances []sdk.Instance `json:"instances"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Instances, nil
}

func runAgent(executionID, state string, triggerCtx map[string]interface{}) {
	log.Printf("Processing Instance %s. State: %s", executionID, state)

	env := os.Environ()
	env = append(env, fmt.Sprintf("GANTRAL_STATUS=%s", state))
	env = append(env, fmt.Sprintf("GANTRAL_EXECUTION_ID=%s", executionID))

	// In a real runner, inject triggerCtx environment vars here

	agentPath := "/agent/agent.py"
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		// Mock execution if the agent doesn't actually exist locally
		log.Printf("Mocking agent execution for %s (no agent found at %s)", executionID, agentPath)
		return
	}

	cmd := exec.Command("python", agentPath)
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			if exitCode == 3 {
				log.Printf("Agent requested Hibernation for %s.", executionID)
			} else {
				log.Printf("Agent failed for %s with code %d. Out: %s", executionID, exitCode, out.String())
			}
		} else {
			log.Printf("Failed to execute agent for %s: %v", executionID, err)
		}
	} else {
		log.Printf("Task Complete for %s. Out: %s", executionID, out.String())
	}
}
