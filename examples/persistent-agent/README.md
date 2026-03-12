# Persistent Agent Example

This example demonstrates how an external system (a "Runner") can integrate with Gantral using the **Persistent Agent** pattern.

## Architecture Let-Down
- **Agent**: A long-running process (e.g., a Python script) that executes a specific task.
- **Runner**: A Go process that bridges the gap between the Agent and Gantral. It uses the **Official Gantral Go SDK** to poll for instances and record human decisions without pulling in internal Gantral kernel dependencies.
- **Gantral Service**: The core workflow engine that dictates policy execution, state changes, and artifact generation.

## Integration Flow
1. A new execution starts in Gantral. It enters state `WAITING_FOR_HUMAN` (if a policy requires human-in-the-loop).
2. The **Runner** (`runner/runner.go`) polls Gantral via the SDK and sees the execution.
3. The Runner launches the **Agent** (`agent/agent.py`) via the local Runtime and passes the Gantral execution ID.
4. If the Agent encounters a condition requiring a human checkbox (e.g., "Hibernation"), the Runner uses `client.RecordDecision()` to send an `APPROVE` or `REJECT` signal back to Gantral.
5. Gantral records the cryptographic decision, resumes the workflow, and updates the task state.

## Local Execution
To spin up the example locally, use the provided `docker-compose.yml` in the `runtime` directory.

```bash
cd runtime
docker-compose up --build
```

The runner will automatically compile and build, pinging the local Gantral core.
