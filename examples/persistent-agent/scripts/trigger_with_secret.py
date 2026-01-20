
import os
import datetime
import requests
import json
import jwt

# Configuration
GANTRAL_URL = os.environ.get("GANTRAL_URL", "http://localhost:8080")
AUTH_SECRET = os.environ.get("DEV_AUTH_SECRET", "dev-secret-key")
USER_ID = "admin-user"

def generate_dev_token():
    """Generates a self-signed HS256 token for Dev Mode."""
    now = datetime.datetime.utcnow()
    payload = {
        "sub": USER_ID,
        "org_id": "org-dev",
        "roles": ["admin"],
        "iat": now,
        "exp": now + datetime.timedelta(minutes=10) 
    }
    encoded = jwt.encode(payload, AUTH_SECRET, algorithm="HS256")
    return encoded

def trigger_workflow():
    token = generate_dev_token()
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    # Payload with Secret Reference
    payload = {
        "workflow_id": "secret-test-flow", 
        "trigger_context": {
            "tier": "prod",
            "environment": {
                "API_KEY": "gantral+secret://env/MY_SECRET_API_KEY",
                "MOCK_VAL": "gantral+secret://vault/my-mock-secret"
            }
        }, 
        "policy": {"materiality": "LOW"}, # LOW to avoid Human Approval for quick test
    
    try:
        print(f"Triggering workflow with SECRETS at {GANTRAL_URL}...")
        resp = requests.post(f"{GANTRAL_URL}/instances", json=payload, headers=headers, timeout=5)
        print(f"Status: {resp.status_code}")
        print(f"Response: {resp.text}")
        resp.raise_for_status()
    except Exception as e:
        print(f"Failed: {e}")

if __name__ == "__main__":
    trigger_workflow()
