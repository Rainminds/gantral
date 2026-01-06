
import time
import requests
import subprocess
import os
import sys
import logging

# Configuration
GANTRAL_URL = os.environ.get("GANTRAL_URL", "http://gantral-core:8080")
POLL_INTERVAL = 3

# Setup Logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)]
)
logger = logging.getLogger("split-runner")

# Track local state to avoid re-running confirmed steps
# Key: execution_id -> Value: "PRE_DONE" or "POST_DONE"
local_state = {}

def get_instances():
    try:
        resp = requests.get(f"{GANTRAL_URL}/instances", timeout=5)
        resp.raise_for_status()
        data = resp.json()
        return data.get("instances", [])
    except requests.exceptions.RequestException as e:
        logger.error(f"Error polling Gantral: {e}")
        return []
    except Exception as e:
        logger.exception("Unexpected error waiting for tasks")
        return []

def run_script(script_path, execution_id, env_vars):
    logger.info(f"Launching {script_path} for {execution_id}...")
    
    env = os.environ.copy()
    env.update(env_vars)
    
    try:
        # Use python directly
        cmd = ["python", script_path]
        result = subprocess.run(cmd, env=env, capture_output=False)
        return result.returncode
    except Exception as e:
        logger.exception(f"Failed to launch script {script_path}")
        return -1

def process_execution(inst):
    exc_id = inst["id"]
    state = inst.get("state", "UNKNOWN")
    
    # Map state to decision for logic compatibility
    decision = "UNKNOWN"
    if state == "APPROVED":
        decision = "APPROVED"
    elif state == "REJECTED":
        decision = "REJECTED"
    elif state == "WAITING_FOR_HUMAN":
        decision = "WAITING"

    current_prog = local_state.get(exc_id, "START")
    
    if current_prog == "START":
        # Always run PRE first for new instances
        logger.info(f"New Instance {exc_id} found. Running Pre-Agent...")
        
        env = {"GANTRAL_EXECUTION_ID": exc_id}
        # Path assumes docker mapping matches
        exit_code = run_script("/app/agent-pre/agent_pre.py", exc_id, env)
        
        if exit_code == 0:
            logger.info(f"Agent-Pre finished successfully for {exc_id}.")
            local_state[exc_id] = "PRE_DONE"
        else:
            logger.error(f"Agent-Pre failed for {exc_id}.")
            local_state[exc_id] = "FAILED"

    elif current_prog == "PRE_DONE":
        # We finished PRE. Now we wait for Decision.
        if decision == "APPROVED":
            logger.info(f"Instance {exc_id} APPROVED. Launching Post-Agent...")
            
            env = {"GANTRAL_EXECUTION_ID": exc_id}
            exit_code = run_script("/app/agent-post/agent_post.py", exc_id, env)
            
            if exit_code == 0:
                 logger.info(f"Agent-Post finished successfully for {exc_id}.")
                 local_state[exc_id] = "POST_DONE"
            else:
                 logger.error(f"Agent-Post failed for {exc_id}.")
                 local_state[exc_id] = "FAILED"
        
        elif decision == "REJECTED":
            logger.info(f"Instance {exc_id} REJECTED. Cleaning up.")
            local_state[exc_id] = "REJECTED"
        
        else:
            # Still waiting (WAITING_FOR_HUMAN)
            pass

def main():
    logger.info("Starting Split-Agent Runner...")
    logger.info(f"Pointing to Core at: {GANTRAL_URL}")

    while True:
        try:
            instances = get_instances()
            for inst in instances:
                process_execution(inst)
        except Exception as e:
            logger.error(f"Loop error: {e}")
        time.sleep(POLL_INTERVAL)

if __name__ == "__main__":
    main()
