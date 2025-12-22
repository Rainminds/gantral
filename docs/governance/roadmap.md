# Roadmap

This roadmap outlines the high-level sequencing for Gantral's development.

## Phase 1: The Core (Months 0-3)
*Focus: Execution State Machine & Basic HITL*

*   [x] **Repository Scaffold:** Monorepo structure, Go modules, initial governance.
*   [ ] **Execution Engine:** Core state machine implementation.
*   [ ] **API & SDK:** Basic gRPC/REST API for creating/managing instances.
*   [ ] **HITL State:** Ability to pause for human input and resume.

## Phase 2: Intelligence (Months 3-6)
*Focus: Policy & Replay*

*   [ ] **Policy Engine:** Declarative rules for when to trigger HITL.
*   [ ] **Deterministic Replay:** Re-run executions from immutable logs.
*   [ ] **SDKs:** Python and TypeScript SDKs.
*   [ ] **CLI:** Developer tools for local management.

## Phase 3: Enterprise Ready (Months 6-9)
*Focus: Audit & Integration*

*   [ ] **Adapter Framework:** Standard integrations (Github, Jira).
*   [ ] **Audit Logging:** Tamper-evident, exportable logs.
*   [ ] **Hardening:** Security reviews, performance optimization.
