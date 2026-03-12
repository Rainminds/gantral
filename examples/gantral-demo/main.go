package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	sdk "github.com/Rainminds/gantral/sdk/go"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	gantralURL := os.Getenv("GANTRAL_URL")
	if gantralURL == "" {
		gantralURL = "http://localhost:8080"
	}

	// 1. Initialize SDK Client
	client := sdk.NewClient(gantralURL)
	ctx := context.Background()

	// 2. Start Instance via SDK
	workflowID := fmt.Sprintf("demo-inst-%d", time.Now().Unix())
	fmt.Printf("🚀 Triggering Workflow: %s\n", workflowID)

	triggerCtx := map[string]interface{}{
		"amount":   1000000,
		"currency": "USD",
		"vendor":   "Acme Corp",
	}
	pol := sdk.Policy{
		ID:                     "policy-demo-high-materiality",
		Materiality:            sdk.MaterialityHigh,
		RequiresHumanApproval:  true,
		ApprovalTimeoutSeconds: 60,
	}

	inst, err := client.CreateInstance(ctx, workflowID, triggerCtx, pol)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	fmt.Printf("✅ Instance Created: %s (State: %s)\n", inst.ID, inst.State)

	// 3. Poll for WAITING_FOR_HUMAN state
	fmt.Println("⏳ Polling for 'WAITING_FOR_HUMAN' state...")
	for {
		inst, err = client.GetInstance(ctx, inst.ID)
		if err != nil {
			log.Printf("Poll error: %v", err)
		} else {
			fmt.Printf("   Current State: %s\n", inst.State)
			if inst.State == sdk.StateWaitingForHuman {
				break
			}
		}
		time.Sleep(2 * time.Second)
	}

	// 4. Record Decision via SDK
	fmt.Println("👤 Recording Human Approval (Auditor: Alice)...")
	inst, err = client.RecordDecision(ctx, inst.ID, sdk.DecisionApprove, "auditor-alice@example.com", "Demo Approval")
	if err != nil {
		log.Fatalf("Failed to record decision: %v", err)
	}
	fmt.Println("✅ Decision Sent: APPROVED")

	// 5. Poll for Completion
	fmt.Println("⏳ Polling for 'COMPLETED' state...")
	for {
		inst, err = client.GetInstance(ctx, inst.ID)
		if err != nil {
			log.Printf("Poll error: %v", err)
		} else {
			fmt.Printf("   Current State: %s\n", inst.State)
			if inst.State == sdk.StateCompleted {
				break
			}
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Printf("\n🎉 Demo Workflow Completed Successfully!\n")

	// 6. Find and provide verification command
	artifactPath := findLatestArtifact(inst.ID)
	if artifactPath != "" {
		fmt.Printf("\n🔍 Proof Found: %s\n", artifactPath)
		fmt.Printf("📋 Verify using:\n")
		fmt.Printf("   go run cmd/gantral-verify/main.go file %s\n", artifactPath)
	}
}

func findLatestArtifact(instanceID string) string {
	files, err := os.ReadDir("./gantral_artifacts")
	if err != nil {
		return ""
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		path := filepath.Join("gantral_artifacts", f.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var art sdk.CommitmentArtifact
		if err := json.Unmarshal(data, &art); err != nil {
			continue
		}

		if art.InstanceID == instanceID {
			return path
		}
	}
	return ""
}
