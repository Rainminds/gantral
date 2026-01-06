
import sys
import json
import os
import time

HANDOFF_FILE = "handoff/context.json"
EXECUTION_ID = os.environ.get("GANTRAL_EXECUTION_ID")

def main():
    print(f"[Agent-Post] Started for Execution {EXECUTION_ID}. Resurrected.")
    
    # 1. Load Handoff Context
    if not os.path.exists(HANDOFF_FILE):
        print(f"[Agent-Post] Critical Error: No handoff file found at {HANDOFF_FILE}")
        sys.exit(1)
        
    with open(HANDOFF_FILE, "r") as f:
        context = json.load(f)
    
    print(f"[Agent-Post] Context loaded: {context}")
    
    # Validate context belongs to this execution (simple safety check)
    if context.get("execution_id") != EXECUTION_ID:
        print(f"[Agent-Post] Warning: Context ID {context.get('execution_id')} matches current execution {EXECUTION_ID}?")
        # proceed anyway for demo
    
    # 2. Perform Sensitive Action
    action = context.get("action", "UNKNOWN")
    print(f"[Agent-Post] Authorized to execute: {action}")
    print("[Agent-Post] Executing...")
    time.sleep(2)
    print("[Agent-Post] ACTION COMPLETE. Chaos unleashed.")
    
    # Clean up
    if os.path.exists(HANDOFF_FILE):
        os.remove(HANDOFF_FILE)
    
    # 3. Exit
    sys.exit(0)

if __name__ == "__main__":
    main()
