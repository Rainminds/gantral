# Progress

## Phase 1: Control Foundations (Gantral Core)
- [x] **1.1 Canonical Execution State Machine** (`core/engine`)
- [x] **1.2 Instance Model & Execution History** (`infra/db/schema.sql`, `audit_events`)
- [x] **1.3 HITL State & Decision Capture** (`WAITING_FOR_HUMAN`, `Decision` inputs)

## Phase 2: Governance Hardening
- [x] **2.1 Policy Evaluation Interface** (`core/policy/types.go`)
- [x] **2.2 Policy Evaluation as Transition Guard** (`core/workflows/execution.go`)
- [x] **2.3 Audit Semantics** (Immutable Event Log)

## Phase 3: Enterprise Integration (Temporal-backed)
- [x] **3.1 Workflow Runtime Integration** (Temporal `GantralExecutionWorkflow`)
- [x] **3.2 Adapters Framework** (`adapters/primary/http`)
- [ ] **3.3 SDKs (Thin Wrappers)** (Pending - Scheduled for Next Sprint)
- [ ] **3.4 Observability & Compliance Outputs** (Pending - Scheduled for Next Sprint)

## Phase 4: Developer Experience & Framework Examples
- [ ] **4.1 Reference Agent Proxy**
- [ ] **4.2 Policy Examples Library**
- [x] **4.3 Consumer Guide** (`docs/guides/example-consumer-integration.md`)
- [ ] **4.4 Framework Reference Implementations**

## Phase 5: Federated Execution & Enterprise Identity
- [ ] **5.1 Identity Federation & Claim Mapping**
- [ ] **5.2 Service Identity**
- [x] **5.3 Runner Protocol & Task Queues** (Implemented via `cmd/worker` split & Temporal Activities)
- [x] **5.4 Connection Registry & Secret Resolution** (ADR 006 - Reference Architecture)
