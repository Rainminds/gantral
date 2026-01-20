
import sys
import json
import os
import requests
import time

# Configuration
GANTRAL_URL = os.environ.get("GANTRAL_URL", "http://gantral-core:8080")
EXECUTION_ID = os.environ.get("GANTRAL_EXECUTION_ID")

def get_handoff_path(execution_id):
    safe_id = "".join([c for c in execution_id if c.isalnum() or c in ('-', '_')])
    return f"handoff/context_{safe_id}.json"

HANDOFF_FILE = get_handoff_path(EXECUTION_ID)

def main():
    print(f"[Agent-Pre] Started for Execution {EXECUTION_ID}")
    
    # 1. Do Work (Simulated)
    print("[Agent-Pre] Gathering intelligence...")
    time.sleep(1)
    print("[Agent-Pre] Identified target: 'Production DB'")
    
    # 2. Prepare Handoff Context
    context_data = {
        "execution_id": EXECUTION_ID,
        "target": "Production DB",
        "action": "DROP TABLE users",
        "risk": "HIGH",
        "timestamp": time.time()
    }
    
    # 3. Save to shared volume
    print(f"[Agent-Pre] Saving context to {HANDOFF_FILE}...")
    with open(HANDOFF_FILE, "w") as f:
        json.dump(context_data, f)
        
    # 4. Call Gantral to Request Decision
    # We update the execution context with the risk factor so policy can evaluate it
    print("[Agent-Pre] Requesting Decision from Gantral...")
    
    try:
        # Call Gantral API to request decision / approval.
        
        payload = {
            "execution_id": EXECUTION_ID,
            "decision": "REQUEST_APPROVAL", # Signaling intent
            "context": context_data
        }
        res = requests.post(f"{GANTRAL_URL}/api/v1/decisions", json=payload)
        if res.status_code >= 400:
             print(f"[Agent-Pre] Warning: Request decision failed: {res.text}")
        else:
             print("[Agent-Pre] Decision requested successfully.")
             
    except Exception as e:
        print(f"[Agent-Pre] Failed to contact Gantral: {e}")
    
    # 5. Exit success
    print("[Agent-Pre] Terminating process.")
    sys.exit(0)

if __name__ == "__main__":
    if not EXECUTION_ID:
        print("Error: GANTRAL_EXECUTION_ID not set")
        sys.exit(1)
    main()
