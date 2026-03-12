# Split Agent Example

This example demonstrates how an external system integrates with Gantral using the **Split Agent** pattern.

## Architecture
- **Agent Pre/Post**: The workflow is split into two distinct ephemeral stages.
- **Runner**: A generic execution bot that runs the pre-tasks, suspends itself until Gantral signals approval, and then wakes up to run post-tasks.
- **Gantral Integration**: Like the persistent agent, the Runner interacts entirely via the **Official Gantral SDK**. It does not leak internal dependencies.

## Local Execution
Spin up the example locally using `docker-compose`:

```bash
docker-compose up --build
```
