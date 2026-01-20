
import sys
import json
import os
import time


def get_checkpoint_path(execution_id):
    # Sanitize execution_id to avoid path traversal
    safe_id = "".join([c for c in execution_id if c.isalnum() or c in ('-', '_')])
    return f"checkpoint/state_{safe_id}.json"

def load_checkpoint(execution_id):
    path = get_checkpoint_path(execution_id)
    if os.path.exists(path):
        with open(path, "r") as f:
            return json.load(f)
    return None

def save_checkpoint(execution_id, state):
    path = get_checkpoint_path(execution_id)
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w") as f:
        json.dump(state, f)

def main():
    # Simulation of receiving status from Runner (infrastructure)
    # The runner will set this env var based on what it got from Gantral
    gantral_status = os.environ.get("GANTRAL_STATUS", "UNKNOWN")
    execution_id = os.environ.get("GANTRAL_EXECUTION_ID", "UNKNOWN")
    
    print(f"[Agent] Started. Execution: {execution_id}. Status from Gantral: {gantral_status}")
    
    # Verify JIT Secrets
    api_key = os.environ.get("API_KEY")
    mock_val = os.environ.get("MOCK_VAL")
    
    if api_key:
        print(f"[Agent] SECRET FOUND: API_KEY={api_key}")
    else:
        print("[Agent] API_KEY not found.")
        
    if mock_val:
        print(f"[Agent] SECRET FOUND: MOCK_VAL={mock_val}")

    state = load_checkpoint(execution_id)

    if not state:
        print("[Agent] No checkpoint found. Starting fresh execution.")
        print("[Agent] Doing Step 1: Gathering Intelligence...")
        time.sleep(1) # Simulate work
        print("[Agent] Step 1 Complete.")
        
        # Now we hit the sensitive step
        print("[Agent] Approaching Sensitive Step: 'Launch Nuclear Missiles'")
        
        # In this persistent-agent model, the Agent checks the status injected by the Runner.
        # If not APPROVED, it hibernates (exit 3).
        
        if gantral_status != "APPROVED":
            print("[Agent] Stop! Sensitive step requires approval.")
            print("[Agent] Saving state to checkpoint...")
            save_checkpoint(execution_id, {"step": "pre_launch", "data": "ready_to_launch"})
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
                path = get_checkpoint_path(execution_id)
                if os.path.exists(path):
                    os.remove(path)
                
                sys.exit(0)
            else:
                # Should not really happen if logic is correct, but safe fallback
                print(f"[Agent] Resumed but status is {gantral_status}. Hibernating again.")
                sys.exit(3)

if __name__ == "__main__":
    main()
