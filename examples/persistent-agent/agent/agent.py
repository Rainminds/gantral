
import sys
import json
import os
import time

CHECKPOINT_FILE = "checkpoint/state.json"

def load_checkpoint():
    if os.path.exists(CHECKPOINT_FILE):
        with open(CHECKPOINT_FILE, "r") as f:
            return json.load(f)
    return None

def save_checkpoint(state):
    os.makedirs(os.path.dirname(CHECKPOINT_FILE), exist_ok=True)
    with open(CHECKPOINT_FILE, "w") as f:
        json.dump(state, f)

def main():
    # Simulation of receiving status from Runner (infrastructure)
    # The runner will set this env var based on what it got from Gantral
    gantral_status = os.environ.get("GANTRAL_STATUS", "UNKNOWN")
    print(f"[Agent] Started. Status from Gantral: {gantral_status}")

    state = load_checkpoint()

    if not state:
        print("[Agent] No checkpoint found. Starting fresh execution.")
        print("[Agent] Doing Step 1: Gathering Intelligence...")
        time.sleep(1) # Simulate work
        print("[Agent] Step 1 Complete.")
        
        # Now we hit the sensitive step
        print("[Agent] Approaching Sensitive Step: 'Launch Nuclear Missiles'")
        
        # In a real app, here we would call Gantral to *request* permission if we didn't have it.
        # But this demo relies on the Runner to have checked Gantral for us?
        # Re-reading prompt: "Runner: Poll Gantral API... subprocess(agent)... Capture Exit Code... if 3 -> suspend"
        # "Agent: If status != APPROVED... Save state... sys.exit(3)"
        
        # So the agent checks permissions *now*.
        # For this demo, valid statuses are probably "PENDING" (default), "WAITING_FOR_HUMAN", "APPROVED".
        
        if gantral_status != "APPROVED":
            print("[Agent] Stop! Sensitive step requires approval.")
            print("[Agent] Saving state to checkpoint...")
            save_checkpoint({"step": "pre_launch", "data": "ready_to_launch"})
            print("HIBERNATING")
            sys.exit(3)
        
    else:
        print(f"[Agent] Checkpoint found: {state}")
        if state.get("step") == "pre_launch":
            if gantral_status == "APPROVED":
                print("[Agent] Approval Granted! Resuming...")
                print("[Agent] Executing Sensitive Step: LAUNCHING MISSILES.")
                time.sleep(1)
                print("[Agent] Missiles Launched.")
                
                # Cleanup
                if os.path.exists(CHECKPOINT_FILE):
                    os.remove(CHECKPOINT_FILE)
                
                sys.exit(0)
            else:
                # Should not really happen if logic is correct, but safe fallback
                print(f"[Agent] Resumed but status is {gantral_status}. Hibernating again.")
                sys.exit(3)

if __name__ == "__main__":
    main()
