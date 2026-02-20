---
title: Implementation Guide
---

# **Enterprise-Grade Engineering Implementation Guide**

## **Version: v6.1 — Repository-Aligned Edition**

Status: Authoritative Engineering Implementation Guide  
Audience:  
• Core maintainers  
• AI coding assistants  
• Platform engineers & SREs  
• Security reviewers  
• Enterprise design partners

This guide operationalizes the Gantral TRD and Implementation Paper.

If this guide conflicts with the TRD, the TRD prevails.

---

# **0\. Governing Rules (Non-Negotiable)**

* Fail-closed behavior is mandatory  
* All authority transitions must be deterministic  
* Authority transitions and artifact emission must be atomic  
* Policy engines are advisory only  
* Logs are never admissible evidence  
* Replay must require only artifacts  
* No execution may advance past governed states without artifact persistence  
* If behavior is ambiguous, execution must stop

Gantral is correctness-first infrastructure.

---

# **1\. Implementation Objectives**

The implementation MUST ensure:

1. Policy thresholds are externalized from workflow code  
2. Authority is canonical workflow state  
3. workflow\_version\_id and policy\_version\_id are bound at decision time  
4. Commitment artifacts form a recursive tamper-evident hash chain  
5. Replay is independent of runtime, database, and logs  
6. Human authority produces attributable reasoning  
7. Policy updates do not require workflow redeployment  
8. Authority visibility is unified across instances

Operational efficiency and admissibility are equal goals.

---

# **2\. Technology Stack (Mandatory)**

## **2.1 Control Plane Language**

* Go (required)

Reason:

* Deterministic concurrency  
* Strong typing  
* Mature Temporal SDK  
* Enterprise readiness

---

## **2.2 Workflow Runtime**

* Temporal (required)  
* One workflow per execution instance  
* Deterministic replay enabled  
* No non-deterministic logic inside workflow definitions

---

## **2.3 Policy Engine**

* Open Policy Agent (OPA)  
* Policies authored in Rego  
* Sidecar or service mode  
* Policy bundles versioned independently

Policy is advisory only.

---

## **2.4 APIs**

* REST (server)  
* gRPC (internal where applicable)  
* OpenAPI-compatible external definitions

---

## **2.5 Datastores**

* PostgreSQL → metadata and indices (non-authoritative)  
* Object storage → commitment artifacts (authoritative)  
* Local artifact storage (`local-storage/artifacts`) for development only  
* Redis (optional, non-authoritative)

---

## **2.6 Infrastructure**

* Kubernetes  
* Helm (recommended)  
* CI via GitHub Actions  
* Signed builds recommended

---

# **3\. Repository Structure (Authoritative Layout)**

The current repository layout is compliant and structured as follows:

## **3.1 Core Execution Plane**

/gantral  
  /cmd  
    /server  
    /worker  
    /gantral-verify  
    /gantral-demo  
    /migrate  
  /core  
    /activities  
    /engine  
    /errors  
    /policy  
    /ports  
    /workflows  
  /adapters  
    /primary  
    /secondary  
  /infra  
    /db  
    /migrations  
  /gantral\_artifacts

Core logic MUST remain inside `/core`.

Adapters MUST contain no business logic.

---

## **3.2 Internal Packages**

/internal  
  /artifact  
  /auth  
  /authority  
  /middleware  
  /policy  
  /replay  
  /storage  
  /workflow

All authority semantics must reside inside:

* `/internal/authority`  
* `/internal/artifact`  
* `/internal/replay`

---

## **3.3 Shared Packages**

/pkg  
  /config  
  /constants  
  /logger  
  /models  
  /verifier

Verifier logic MUST remain independent of runtime.

---

## **3.4 Testing**

/tests  
  /unit  
  /statemachine  
  /artifact  
  /replay  
  /integration  
  /golden  
  /helpers

Golden tests must validate artifact chain stability.

---

# **4\. Core Domain Models**

## **4.1 ExecutionInstance**

Fields:

* instance\_id (UUID, immutable)  
* workflow\_id  
* workflow\_version\_id  
* owning\_team\_id  
* current\_state  
* created\_at  
* terminated\_at  
* cost\_metadata

Instances are append-only in state progression.

---

## **4.2 AuthorityDecision**

Fields:

* decision\_id  
* instance\_id  
* decision\_type (APPROVE / REJECT / OVERRIDE)  
* human\_actor\_id  
* role  
* justification  
* context\_snapshot\_hash  
* timestamp

Rules:

* justification MUST be non-empty for APPROVE and OVERRIDE unless configured otherwise  
* identity must be validated via OIDC before persistence

---

## **4.3 CommitmentArtifact**

Fields:

* artifact\_version  
* artifact\_id  
* instance\_id  
* workflow\_version\_id  
* prev\_artifact\_hash  
* authority\_state  
* policy\_version\_id  
* context\_snapshot\_hash  
* human\_actor\_id  
* justification  
* timestamp  
* artifact\_hash

Artifact MUST bind:

* workflow version  
* policy version  
* authority state  
* identity  
* authority-relevant context

---

# **5\. Authority State Machine (Executable Rules)**

Allowed transitions ONLY:

CREATED → RUNNING  
RUNNING → WAITING\_FOR\_HUMAN  
WAITING\_FOR\_HUMAN → APPROVED  
WAITING\_FOR\_HUMAN → REJECTED  
WAITING\_FOR\_HUMAN → OVERRIDDEN  
APPROVED → RESUMED  
OVERRIDDEN → RESUMED  
RESUMED → RUNNING  
RUNNING → COMPLETED  
RUNNING → TERMINATED

Any other transition MUST panic and terminate execution.

No implicit recovery allowed.

---

# **6\. Temporal Workflow Implementation**

* One Temporal workflow per execution instance  
* Workflow history authoritative only for ordering  
* All authority decisions recorded as workflow events  
* No random values  
* No wall-clock calls  
* No external I/O inside workflow logic

Activities handle external side effects.

---

# **7\. Policy Evaluation (OPA)**

## **7.1 Invocation**

* Invoked synchronously during transition guard  
* Input schema must be versioned  
* OPA output: ALLOW / REQUIRE\_HUMAN / DENY

## **7.2 Version Binding**

* policy\_version\_id must be retrieved from bundle  
* policy\_version\_id must be embedded in artifact  
* Policy changes must not require workflow redeployment

Policy never commits authority.

---

# **8\. Commitment Artifact Emission**

## **8.1 Atomicity**

Authority transition and artifact emission MUST be atomic.

If artifact persistence fails:

* State MUST NOT advance  
* Workflow MUST remain in WAITING\_FOR\_HUMAN  
* No retry loops

---

## **8.2 Hash Chain Model**

artifact\_hash\_i \=

* H(payload\_i) if first artifact  
* H(payload\_i || artifact\_hash\_{i-1}) otherwise

Modification of any artifact invalidates subsequent chain.

---

## **8.3 Artifact Storage**

Production:

* Append-only object storage  
* Write-once configuration  
* No mutation APIs

Development:

* local-storage/artifacts (non-authoritative)

`gantral_artifacts` directory must not allow mutation in production mode.

---

# **9\. Replay & Verification**

Verifier location:

/pkg/verifier  
/cmd/gantral-verify

Verifier MUST:

* Operate offline  
* Require no network  
* Accept artifact file(s)  
* Output VALID / INVALID / INCONCLUSIVE

Replay validates:

1. Hash integrity  
2. Transition validity  
3. workflow\_version\_id consistency  
4. policy\_version\_id consistency

Replay reconstructs authority-state projection only.

---

# **10\. Failure Semantics**

Terminate execution on:

* Missing artifact  
* Hash mismatch  
* Illegal transition  
* Identity ambiguity  
* Policy ambiguity  
* Version mismatch  
* Temporal non-determinism

Fail closed always.

---

# **11\. Identity & Security**

* OIDC federation only  
* No local users  
* No password storage  
* No secret persistence  
* Workload identity required  
* Artifact emission only after identity validation

---

# **12\. Operational Efficiency Requirements**

Implementation MUST:

* Prevent governance logic inside workflow code  
* Externalize approval thresholds  
* Avoid workflow forks based solely on thresholds  
* Allow policy bundle updates without redeploying workflows

This reduces:

* Code duplication  
* Redeployment risk  
* Governance drift

---

# **13\. Unified Authority Visibility**

System MUST provide:

* Queryable execution instances  
* Visible WAITING\_FOR\_HUMAN states  
* Instance-level isolation  
* Historical authority progression

No hidden execution state allowed.

---

# **14\. Testing Requirements**

Mandatory tests:

* Transition correctness  
* Atomic artifact emission  
* Hash chain validation  
* Replay determinism  
* Version consistency  
* Fuzz tests for artifact corruption  
* Golden replay stability tests

Replay compatibility cannot break without major version change.

---

# **15\. AI Coding Assistant Rules**

AI assistants MUST:

* Implement only documented transitions  
* Never auto-approve  
* Never infer missing authority  
* Never embed policy logic in workflows  
* Fail on ambiguity  
* Preserve atomicity  
* Preserve version binding

Correctness \> convenience.

---

# **16\. Enterprise Deployment Expectations**

Production deployments should:

* Enforce write-once artifact buckets  
* Separate runtime and control plane  
* Enable justification enforcement  
* Pin policy bundle versions  
* Monitor WAITING\_FOR\_HUMAN backlog

Gantral must be incrementally adoptable.

Agents do not require memory model modification.

---

# **17\. Final Statement**

If implemented according to this guide:

* Policy–code duplication is eliminated  
* Authority becomes canonical state  
* Chain-of-custody becomes cryptographically verifiable  
* Replay is independent of logs  
* Operational fragmentation is structurally removed

Gantral becomes deterministic execution infrastructure by construction.

---

