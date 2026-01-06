
import sys
import json
import os
import requests
import time

# Configuration
GANTRAL_URL = os.environ.get("GANTRAL_URL", "http://gantral-core:8080")
EXECUTION_ID = os.environ.get("GANTRAL_EXECUTION_ID")
HANDOFF_FILE = "handoff/context.json"

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
        # First, update context (if supported) or send as part of decision request.
        # Assuming we can trigger a policy check.
        # In this demo, the Runner might handle the "check", OR we explicitly ask.
        # We will try to post a 'request'.
        
        # NOTE: Using the /decisions endpoint to signal we want a decision?
        # Or maybe /eval?
        # Let's assume we update the execution with the context first.
        # (Assuming endpoint exists, based on prev experience)
        # requests.put(f"{GANTRAL_URL}/api/v1/executions/{EXECUTION_ID}", json=context_data) 
        
        # But wait, user said "Call Gantral: 'Request Decision'".
        # I'll simulate this by posting a "REQUEST" decision, or relying on the policy check 
        # that happens when we complete?
        # Prompt: "If Task 1 succeeds -> Report 'COMPLETED' -> Gantral Core state machine will then evaluate policy."
        # Ah! So simply completing this task might trigger the policy check if the Runner reports it.
        # BUT User instruction for 'agent_pre.py' specifically says "4. Call Gantral: 'Request Decision'".
        # So I will add a call here.
        
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
