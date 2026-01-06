
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
logger = logging.getLogger("gantral-runner")

seen_states = {}

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
        logger.exception("Unexpected error during polling")
        return []

def run_agent(instance, state):
    execution_id = instance["id"]
    
    logger.info(f"Processing Instance {execution_id}. State: {state}")
    
    env = os.environ.copy()
    env["GANTRAL_STATUS"] = state
    env["GANTRAL_EXECUTION_ID"] = execution_id
    
    agent_path = "/agent/agent.py" 
    
    try:
        result = subprocess.run(["python", agent_path], env=env, capture_output=False)
        exit_code = result.returncode
        
        logger.info(f"Agent exited with code {exit_code}")
        
        if exit_code == 0:
            logger.info(f"Task Complete for {execution_id}. (No callback available)")
        elif exit_code == 3:
            logger.info(f"Agent requested Hibernation for {execution_id}.")
        else:
            logger.error(f"Agent failed for {execution_id} with code {exit_code}.")
            
    except Exception as e:
        logger.exception(f"Failed to execute agent for {execution_id}")

def main():
    logger.info("Starting Gantral Reference Runner...")
    logger.info(f"Pointing to Core at: {GANTRAL_URL}")

    while True:
        try:
            instances = get_instances()
            for inst in instances:
                inst_id = inst["id"]
                state = inst.get("state", "UNKNOWN")
                
                last_state = seen_states.get(inst_id)
                
                # Only react on state changes or new instances
                if state != last_state:
                    seen_states[inst_id] = state
                    
                    if state == "WAITING_FOR_HUMAN":
                        logger.info(f"Instance {inst_id} is WAITING. Launching Agent to hibernate...")
                        run_agent(inst, state)
                    elif state == "APPROVED":
                        logger.info(f"Instance {inst_id} is APPROVED. Launching Agent to resume...")
                        run_agent(inst, state)
                    elif state == "COMPLETED":
                        logger.info(f"Instance {inst_id} is COMPLETED. Skipping.")
        
        except Exception as e:
            logger.error(f"Loop error: {e}")
            
        time.sleep(POLL_INTERVAL)

if __name__ == "__main__":
    main()
