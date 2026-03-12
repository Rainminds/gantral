---
title: Implementation Guide
---

# **Enterprise-Grade Engineering Implementation Guide**

Version: v8.0 — TRD-Aligned Open Infrastructure Edition  
Status: Authoritative Engineering Implementation Guide

Audience

* Core maintainers  
* AI coding assistants  
* Platform engineers & SREs  
* Security reviewers  
* Enterprise architects  
* External contributors and design partners

This guide operationalizes the **Gantral Technical Reference Document (TRD)** and the **Gantral Implementation Paper**.

If any behavior described here conflicts with the TRD, **the TRD prevails**.

---

# **0\. Governing Rules (Non-Negotiable)**

Gantral is correctness-first infrastructure.

All implementations must obey the following rules:

* Fail-closed behavior is mandatory.  
* All authority transitions must be deterministic.  
* Authority transitions and artifact emission must be atomic.  
* Policy engines are advisory only.  
* Logs are never admissible evidence.  
* Replay must require only commitment artifacts.  
* No execution may advance past governed states without artifact persistence.  
* Artifact structure must conform exactly to the TRD Artifact Specification (Appendix A).  
* If behavior is ambiguous, execution must stop.

---

# **1\. Normative Specification References**

The following specifications are defined in the TRD and **MUST NOT be redefined in this document**.

| Topic | Normative Source |
| ----- | ----- |
| Canonical execution state machine | TRD §4 |
| Allowed state transitions | TRD §4.2 |
| Commitment artifact structure | TRD Appendix A |
| Artifact hashing model | TRD Appendix A.3 |
| Artifact signing | TRD Appendix A.4 |
| Timestamp binding | TRD Appendix A.5 |
| Canonical serialization rules | TRD Appendix A.9–A.10 |
| Replay validation rules | TRD §7 |

This implementation guide explains **how to build compliant systems**, not the canonical definitions themselves.

---

# **2\. Implementation Objectives**

The implementation MUST ensure:

1. Policy thresholds are externalized from workflow code.  
2. Authority is canonical workflow state.  
3. `workflow_version_id` and `policy_version_id` are bound at decision time.  
4. Commitment artifacts form a tamper-evident hash chain.  
5. Replay is independent of runtime, database, and logs.  
6. Human authority produces attributable reasoning.  
7. Policy updates do not require workflow redeployment.  
8. Authority visibility is unified across execution instances.

Operational efficiency and admissibility are equal goals.

---

# **3\. Technology Stack (Mandatory)**

## **3.1 Control Plane Language**

Required: **Go**

Reasons:

* Deterministic concurrency  
* Strong typing  
* Mature Temporal SDK  
* Enterprise production suitability

---

## **3.2 Workflow Runtime**

Required: **Temporal**

Requirements:

* One workflow per execution instance  
* Deterministic replay enabled  
* No non-deterministic logic inside workflow definitions

Temporal provides:

* ordering  
* durability  
* deterministic replay

Gantral provides:

* authority enforcement

---

## **3.3 Policy Engine**

Recommended: **Open Policy Agent (OPA)**

Policies:

* written in Rego  
* versioned independently  
* evaluated at authority checkpoints

Policy decisions are **advisory signals only**.

Gantral interprets and enforces them.

---

## **3.4 APIs**

External:

* REST

Internal:

* gRPC where appropriate

External APIs should expose **OpenAPI-compatible schemas**.

---

## **3.5 Datastores**

Non-authoritative stores:

PostgreSQL  
Used for:

* metadata  
* indices  
* instance queries

Authoritative store:

Object storage (append-only)

Used for:

* commitment artifacts

Optional:

Redis (non-authoritative caching)

---

## **3.6 Infrastructure**

Recommended:

* Kubernetes  
* Helm charts  
* GitHub Actions CI  
* Signed builds

---

# **4\. Repository Structure (Authoritative Layout)**

## **4.1 Core Execution Plane**

/gantral  
  /cmd  
    /server  
    /worker  
    /gantral-verify  
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

Core logic MUST reside inside `/core`.

Adapters must contain **no business logic**.

---

## **4.2 Internal Packages**

/internal  
  /artifact  
  /auth  
  /authority  
  /middleware  
  /policy  
  /replay  
  /storage  
  /workflow

Authority semantics must reside in:

* `/internal/authority`  
* `/internal/artifact`  
* `/internal/replay`

---

## **4.3 Shared Packages**

/pkg  
  /config  
  /constants  
  /logger  
  /models  
  /verifier

Verifier logic MUST remain **independent of runtime**.

---

## **4.4 Testing**

/tests  
  /unit  
  /statemachine  
  /artifact  
  /replay  
  /integration  
  /golden  
  /helpers

Golden tests must validate **artifact chain stability across versions**.

---

# **5\. Core Domain Models**

## **5.1 ExecutionInstance**

Fields:

* instance\_id (UUID, immutable)  
* workflow\_id  
* workflow\_version\_id  
* owning\_team\_id  
* current\_state  
* created\_at  
* terminated\_at  
* cost\_metadata

Execution instances are **append-only in state progression**.

---

## **5.2 AuthorityDecision**

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

* justification must be present for APPROVE and OVERRIDE unless explicitly configured otherwise  
* identity must be validated via OIDC before persistence

---

## **5.3 Commitment Artifact**

Artifact structure is defined exclusively in the **TRD Appendix A**.

This implementation guide does **not redefine the schema**.

Implementations must:

* construct artifacts exactly as specified in the TRD  
* follow canonical serialization rules  
* compute artifact hashes exactly as defined in the TRD  
* include all mandatory fields

The artifact specification in the TRD is the **single normative definition**.

---

# **6\. Authority State Machine (Executable Rules)**

Allowed transitions:

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

Any other transition:

→ **terminate execution immediately**

No implicit recovery allowed.

---

# **7\. Temporal Workflow Implementation**

Requirements:

* One workflow per execution instance  
* Workflow history authoritative for ordering only  
* Authority decisions recorded as workflow events

Forbidden inside workflow logic:

* random values  
* wall-clock calls  
* external network calls

External effects must occur in **activities**.

---

# **8\. Policy Evaluation**

Policy evaluation acts as a **transition guard**.

Inputs include:

* instance\_id  
* workflow\_version\_id  
* materiality  
* actor\_id  
* roles  
* policy\_version\_id

Outputs:

* ALLOW  
* REQUIRE\_HUMAN  
* DENY

Gantral interprets the signal.

Policy **does not enforce execution**.

---

# **9\. Commitment Artifact Emission**

## **9.1 Atomicity**

Authority transition and artifact emission must be atomic.

If artifact persistence fails:

* state transition must not occur  
* workflow remains in `WAITING_FOR_HUMAN`

Partial transitions are forbidden.

---

## **9.2 Reference Artifact Emission Algorithm**

The following pseudocode illustrates the required artifact emission process.

The artifact schema and serialization rules are defined in the TRD.

function emitAuthorityArtifact(decision):

  payload \= constructPayload(decision)

  serialized \= canonicalSerialize(payload)

  if prev\_artifact\_hash exists:  
      artifact\_hash \= H(serialized || prev\_artifact\_hash)  
  else:  
      artifact\_hash \= H(serialized)

  signature \= Sign(private\_key, artifact\_hash)

  timestamp\_token \= TSA.timestamp(artifact\_hash)

  artifact \= \{  
      payload fields,  
      prev\_artifact\_hash,  
      artifact\_hash,  
      artifact\_signature,  
      timestamp\_token  
  \}

  persistAppendOnly(artifact)

  return artifact

Implementations must follow the **exact hashing and serialization rules defined in the TRD**.

---

## **9.3 Artifact Storage**

Production requirements:

* append-only object storage  
* write-once configuration  
* no mutation APIs

Development:

* local artifact directory

Artifacts are the **authoritative record for replay**.

---

# **10\. Replay & Verification**

Replay tools must:

* operate offline  
* require no network  
* accept artifact chains

Outputs:

* VALID  
* INVALID  
* INCONCLUSIVE

Replay verifies:

1. hash-chain integrity  
2. transition validity  
3. workflow version consistency  
4. policy version consistency

Replay reconstructs **authority progression only**.

Agent memory is excluded.

---

# **11\. Replay Compatibility Contract**

Gantral guarantees that commitment artifacts remain **replay-verifiable across implementation versions**.

Implementations MUST follow these rules:

1. Artifact fields defined in the TRD MUST NOT be renamed or removed.  
2. Artifact payload ordering must follow canonical serialization rules defined in the TRD.  
3. New artifact versions must increment `artifact_version`.  
4. Replay tools must support verification of all prior artifact versions.  
5. Replay logic must never depend on:  
   * runtime state  
   * databases  
   * logs  
   * external services

Independent verifier implementations must be able to validate artifacts using only:

* the artifact chain  
* the public hashing algorithm  
* the public signature algorithm

---

# **12\. Failure Semantics**

Execution must terminate on:

* missing artifact  
* hash mismatch  
* illegal transition  
* identity ambiguity  
* policy ambiguity  
* version mismatch  
* Temporal nondeterminism

Fail closed always.

---

# **13\. Identity & Security**

Identity requirements:

* OIDC federation only  
* no local user database  
* no password storage  
* workload identity required

Artifact emission occurs **only after identity verification**.

---

# **14\. Operational Efficiency Requirements**

Gantral must:

* prevent governance logic inside workflow code  
* externalize approval thresholds  
* avoid workflow forks caused only by policy differences  
* allow policy updates without redeployment

This eliminates **policy–code duplication**.

---

# **15\. Unified Authority Visibility**

Gantral must expose:

* active execution instances  
* workflows waiting for authority  
* historical authority progression

Gantral intentionally avoids dashboards or analytics layers.

---

# **16\. Implementation Boundary Reminder**

Gantral enforces:

* execution authority  
* deterministic transitions  
* artifact emission  
* replayable evidence

Gantral does **not provide**:

* dashboards  
* policy lifecycle management  
* cross-workflow analytics  
* autonomy orchestration  
* enterprise management UI

---

