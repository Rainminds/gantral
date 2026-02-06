# Gantral – Enterprise-Grade Engineering Implementation Guide (AI-Coding-Ready)

Version: v6.0 (Implementation-Executable)  
Status: Authoritative engineering implementation guide  
Audience:  
• AI coding assistants acting as primary implementers  
• Core maintainers reviewing AI-generated code  
• Platform engineers and SREs  
• Security, Risk, and Compliance reviewers

This document is written to be directly consumed by AI coding assistants.  
It specifies languages, frameworks, data models, invariants, failure semantics,  
and production standards explicitly.

## 0\. Governing Rules (Non-Negotiable)

• TRD invariants override this document  
• Fail-closed behavior is mandatory  
• No best-effort execution is permitted  
• All authority transitions must be deterministic  
• Logs are never evidence  
• If behavior is ambiguous, execution must stop

## 1\. Technology Stack (Mandatory)

Control Plane Language:  
• Go (primary, required)  
• Reason: deterministic concurrency, strong typing, Temporal SDK maturity

Workflow Runtime:  
• Temporal (required)  
• One workflow per execution instance  
• Deterministic replay enabled

Policy Engine:  
• Open Policy Agent (OPA)  
• Rego policies  
• Sidecar or service mode

APIs:  
• gRPC (internal)  
• REST (external, OpenAPI 3.1)

Datastores:  
• PostgreSQL (metadata, indices only)  
• Object Storage (artifacts, immutable)  
• Redis (optional, non-authoritative caching)

Infrastructure:  
• Kubernetes (required)  
• Helm for deployment  
• GitHub Actions for CI

## 2\. Repository Structure (Required)

/cmd  
  /gantral-api  
  /gantral-worker  
/internal  
  /authority  
  /workflow  
  /policy  
  /artifact  
  /replay  
  /identity  
  /storage  
/pkg  
  /sdk  
  /models  
/docs  
/tests

## 3\. Core Domain Models (Code-Level)

ExecutionInstance:  
• instance\_id (UUID, immutable)  
• workflow\_id  
• workflow\_version  
• owning\_team\_id  
• current\_state  
• created\_at  
• terminated\_at

AuthorityDecision:  
• decision\_id  
• instance\_id  
• decision\_type  
• human\_actor\_id  
• role  
• justification  
• context\_snapshot\_hash  
• timestamp

CommitmentArtifact:  
• artifact\_version  
• artifact\_id  
• instance\_id  
• prev\_artifact\_hash  
• authority\_state  
• policy\_version\_id  
• context\_hash  
• human\_actor\_id  
• timestamp  
• artifact\_hash

## 4\. Authority State Machine (Executable Rules)

Allowed transitions ONLY:  
CREATED → RUNNING  
RUNNING → WAITING\_FOR\_HUMAN  
WAITING\_FOR\_HUMAN → APPROVED | REJECTED | OVERRIDDEN  
APPROVED | OVERRIDDEN → RESUMED  
RESUMED → RUNNING  
RUNNING → COMPLETED | TERMINATED

Any other transition MUST panic and terminate execution.

## 5\. Temporal Workflow Implementation

• One Temporal workflow per execution instance  
• Workflow history is authoritative for ordering only  
• All authority decisions recorded as workflow events  
• Workflow code must be 100% deterministic  
• No random, wall-clock, or external I/O in workflow logic

## 6\. Policy Evaluation (OPA – Advisory Only)

• OPA invoked synchronously as transition guard  
• Input schema fixed and versioned  
• OPA output NEVER persisted as authority  
• ALLOW / REQUIRE\_HUMAN / DENY only

## 7\. Commitment Artifact Emission (Critical)

• Artifact emitted atomically with authority transition  
• Hash chain required  
• Write-once storage  
• Artifact emission failure MUST abort execution  
• No retries

## 8\. Offline Replay & Verification

Verifier Requirements:  
• Standalone binary or library  
• No network calls  
• Input: artifact(s)  
• Output: VALID / INVALID / INCONCLUSIVE  
• Replay authority decisions only

## 9\. Failure Semantics (Strict)

Terminate execution on:  
• Missing artifact  
• Hash mismatch  
• Policy ambiguity  
• Authority ambiguity  
• Temporal non-determinism

## 10\. Security & Identity

• OIDC federation only  
• No local users  
• Workload identity for services  
• Secrets resolved at runtime edge

## 11\. CI / Release / OSS Standards

• Mandatory unit tests for state transitions  
• Replay tests required  
• Fuzz tests for artifact corruption  
• Semantic versioning  
• Backward-compatible artifact schema  
• Signed releases

## 12\. AI Coding Assistant Instructions

AI assistants MUST:  
• Implement only documented behavior  
• Never infer missing authority  
• Never add convenience shortcuts  
• Fail on ambiguity  
• Prefer correctness over performance

## Final Statement

This document is designed to make incorrect implementations impossible  
to justify.

If Gantral is implemented according to this guide and the TRD,  
it will produce execution systems that are verifiable, replayable,  
and enterprise-grade by construction.

