
import jwt
import requests
import datetime
import os
import sys

# Configuration matches gantral-core defaults
SECRET = os.environ.get("DEV_AUTH_SECRET", "dev-secret-key")
GANTRAL_URL = os.environ.get("GANTRAL_URL", "http://localhost:8080")

def generate_token(sub, role, identity_type):
    now = datetime.datetime.now(datetime.timezone.utc)
    payload = {
        "sub": sub,
        "type": identity_type,
        "roles": [role],
        "iat": now,
        "exp": now + datetime.timedelta(minutes=5)
    }
    return jwt.encode(payload, SECRET, algorithm="HS256")

def test_forbidden(endpoint, role, identity_type, expected_status):
    print(f"Testing {endpoint} with Role={role}, Type={identity_type}...")
    token = generate_token("tester", role, identity_type)
    headers = {"Authorization": f"Bearer {token}"}
    
    # Try to hit the endpoint (POST)
    try:
        resp = requests.post(f"{GANTRAL_URL}{endpoint}", headers=headers, json={"dummy": "data"})
        print(f"Response: {resp.status_code}")
        
        if resp.status_code == expected_status:
            print("✅ PASS: Correct status code received.")
            return True
        else:
            print(f"❌ FAIL: Expected {expected_status}, got {resp.status_code}")
            print(f"Body: {resp.text}")
            return False
    except Exception as e:
        print(f"Error: {e}")
        return False

if __name__ == "__main__":
    success = True
    
    # 1. Machine (Role: runner) trying to Approve (POST /decisions) -> Should be 403 Forbidden
    if not test_forbidden("/decisions", "runner", "machine", 403):
        success = False
        
    # 2. User (Role: admin) trying to Poll (POST /tasks/poll) -> Should be 403 Forbidden
    if not test_forbidden("/tasks/poll", "admin", "human", 403):
        success = False
        
    # 3. User (Role: user) trying to Approve (POST /decisions) -> Should be 200/201 (Authorized)
    # Note: Request body might fail validation, but Auth should pass. 
    # Core returns 400 Bad Request if Auth passes but Body bad. 401/403 means Auth fail.
    # We expect NOT 401/403.
    print(f"Testing /decisions with Role=user, Type=human...")
    token = generate_token("human-tester", "user", "human")
    headers = {"Authorization": f"Bearer {token}"}
    resp = requests.post(f"{GANTRAL_URL}/decisions", headers=headers, json={}) # Empty body -> 400 or 500
    if resp.status_code in [401, 403]:
        print(f"❌ FAIL: User should be authorized, got {resp.status_code}")
        success = False
    else:
        print(f"✅ PASS: User authorized (Status {resp.status_code})")

    if success:
        print("\nAll RBAC tests passed!")
        sys.exit(0)
    else:
        print("\nSome tests failed.")
        sys.exit(1)
