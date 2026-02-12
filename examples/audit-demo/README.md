# Auditor Verification Demo

This demo simulates a "Catastrophic Failure" scenario to prove the **Admissibility** of Gantral's cryptographic evidence.

## The Scenario
1.  **Operation:** A high-value workflow (e.g., "$1M Transfer") is executed and approved by a human.
2.  **Evidence:** Gantral generates an immutable **Commitment Artifact** and stores it in the `artifacts/` volume.
3.  **Disaster:** The Gantral Server and Database are destroyed (simulated by `docker compose stop`).
4.  **The Audit:** An auditor arrives 6 months later. They have no access to the original database or logs. They only have the file system artifacts.

## The Objective
Prove that the auditor can cryptographically verify:
-   **Who** approved the action.
-   **What** was approved.
-   **When** it happened.
-   **Integrity:** That the data has not been tampered with.

## Usage

```bash
./run-demo.sh
```

## Expected Output
The script will succeed and print a "CHAIN VALID" result with a verbose summary of the admissible evidence.
