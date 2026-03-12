---
title: Technical Reference Document
---

# **Technical Reference & Architecture Document (TRD)**

Version: v7.0 — Artifact Specification Edition  
Status: Authoritative Technical Constitution

Audience:

* Core contributors  
* Platform engineers  
* Security reviewers  
* Enterprise architects  
* Design partners

---

# **1\. Purpose**

This document defines the authoritative technical architecture and execution semantics of Gantral.

Gantral is an **AI Execution Control Plane** that enforces deterministic authority transitions in agentic workflows.

This document defines:

* Architectural invariants  
* Canonical execution model  
* Authority semantics  
* Commitment artifact model  
* Hash integrity properties  
* Replay guarantees  
* Responsibility boundaries  
* Enterprise adoption characteristics

If an implementation conflicts with this document, the implementation is incorrect.

The **Commitment Artifact structure is defined in Appendix A — Artifact Specification v1.**

All implementations MUST conform exactly to that specification.

---

# **1.1 Structural Challenges Addressed**

Enterprise adoption of agentic AI introduces structural failures not solved by agent frameworks or orchestration tools alone.

## **Operational Inefficiency**

Governance thresholds are frequently embedded directly within workflow code:

* Monetary limits hardcoded in agents  
* Risk thresholds implemented inside orchestration logic  
* Team-specific approval rules implemented as workflow forks  
* Environment-specific deployments differing only by governance parameters

This creates:

* Policy–code duplication  
* Redeployment risk for policy changes  
* Governance drift between documented and deployed behavior  
* Slower adaptation to regulatory or risk updates

Gantral eliminates policy–code duplication by separating policy thresholds from workflow implementation and binding authority decisions to versioned policy bundles.

---

## **Operational Fragmentation**

Authority decisions are often evaluated in one system but enforced in another:

* Policy engines detached from execution runtime  
* Human approvals recorded outside canonical workflow state  
* Logs capturing events without structural binding to execution progression  
* Multiple agent frameworks lacking shared authority semantics

Gantral represents authority directly as canonical workflow state transitions.

---

## **Broken Chain of Custody**

At AI execution boundaries:

1. An agent proposes an action  
2. A human approves  
3. Execution resumes  
4. Logs attempt reconstruction

Without cryptographic binding, reconstruction depends on:

* Log integrity  
* Runtime availability  
* Policy version drift  
* Operator testimony

Gantral emits **cryptographically chained commitment artifacts** enabling deterministic replay independent of logs or runtime access.

---

# **2\. Scope and Non-Goals**

## **2.1 What Gantral Is**

Gantral is:

* An execution control plane for AI-assisted workflows  
* A deterministic authority state machine  
* An atomic commitment artifact emitter  
* A replay-verifiable execution authority ledger  
* A policy-advisory integration layer  
* Infrastructure positioned above agents and below enterprise accountability systems

Gantral represents authority as canonical workflow state.

---

## **2.2 Explicit Non-Goals**

Gantral does not:

* Build, host, or orchestrate AI agents  
* Persist agent memory or reasoning traces  
* Encode domain-specific business logic  
* Replace workflow runtimes  
* Replace CI/CD or ITSM systems  
* Act as an identity provider  
* Store secrets or credentials  
* Provide autonomous self-approval mechanisms  
* Guarantee regulatory compliance

Gantral enforces **authority transitions only**.

---

# **3\. Architectural Invariants**

The following invariants are non-negotiable.

## **3.1 Instance-First Semantics**

All authority, audit, and replay semantics attach to immutable execution instances.

---

## **3.2 Authority Is State**

Authority exists only as canonical workflow state transitions.

---

## **3.3 HITL Is a Blocking State**

Human-in-the-loop is modeled as `WAITING_FOR_HUMAN`, not as an external notification.

---

## **3.4 Atomic Authority Commitment**

Authority transition and commitment artifact emission occur atomically.

If artifact persistence fails:

* The authority state transition MUST NOT be observable  
* Execution MUST remain in `WAITING_FOR_HUMAN`  
* Partial transitions are forbidden

---

## **3.5 Determinism Over Performance**

Replay correctness supersedes latency optimization.

---

## **3.6 Policy Is Advisory**

Policy engines provide advisory signals only.

Final authority is represented exclusively as workflow state transitions.

---

## **3.7 Version Binding**

Authority decisions must bind:

* `workflow_version_id`  
* `policy_version_id`

Version mismatch during replay yields `INCONCLUSIVE`.

---

## **3.8 Agent State Separation**

Gantral never persists:

* Agent memory  
* Internal plans  
* Tool traces

Agent persistence is the responsibility of agent frameworks.

---

## **3.9 Determinism & Evidence**

Human authority must produce structured reasoning.

Commitment artifacts include:

* `human_actor_id`  
* `justification`

Deployments may enforce minimum reasoning requirements.

---

# **4\. Formal Execution Model**

## **4.1 Canonical State Set**

The canonical state set S consists of:

* CREATED  
* RUNNING  
* WAITING\_FOR\_HUMAN  
* APPROVED  
* REJECTED  
* OVERRIDDEN  
* RESUMED  
* TERMINATED  
* COMPLETED

---

## **4.2 Transition Relation**

Allowed transitions:

* CREATED → RUNNING  
* RUNNING → WAITING\_FOR\_HUMAN  
* WAITING\_FOR\_HUMAN → APPROVED  
* WAITING\_FOR\_HUMAN → REJECTED  
* WAITING\_FOR\_HUMAN → OVERRIDDEN  
* APPROVED → RESUMED  
* OVERRIDDEN → RESUMED  
* RESUMED → RUNNING  
* RUNNING → COMPLETED  
* RUNNING → TERMINATED

Transitions not enumerated above are illegal.

---

## **4.3 Transition Validity**

Let the state sequence for execution instance E be:

σE \= (s₀, s₁, …, sₙ)

For all i:

(sᵢ, sᵢ₊₁) must belong to the allowed transition relation.

Otherwise execution is invalid.

---

# **5\. Commitment Artifacts**

At each authority transition Gantral emits a **Commitment Artifact**.

The artifact structure is defined in **Appendix A — Artifact Specification v1**.

Artifacts provide:

* identity binding  
* policy binding  
* workflow version binding  
* context binding  
* cryptographic ordering

Artifacts form a **tamper-evident hash chain**.

---

## **5.1 Identity Provenance**

Identity is validated using federated OIDC.

Gantral records the verified subject identifier as `human_actor_id`.

Gantral does not maintain its own identity directory.

---

## **5.2 Context Snapshot**

`context_snapshot_hash` binds:

* workflow parameters  
* policy evaluation inputs  
* authority-relevant request payload

It does not persist agent memory.

---

# **6\. Hash Integrity Model**

artifact\_hash computation MUST follow the algorithm defined in Appendix A — Artifact Specification v1.

Integrity property:

Modification of any artifact invalidates the chain from that point forward.

---

# **7\. Replay Determinism**

Replay validation MUST verify:

1. Hash-chain integrity  
2. Valid state transitions  
3. workflow\_version\_id consistency  
4. policy\_version\_id consistency

Replay outcomes:

* VALID  
* INVALID  
* INCONCLUSIVE

Replay verification MUST classify artifact chains whose final artifact represents a non-terminal authority state as INCONCLUSIVE rather than VALID.

Replay requires **no runtime, database, or logs**.

---

# **8\. Reference Architecture**

Enterprise Systems  
↓  
Gantral Control Plane  
↓  
Deterministic Workflow Runtime  
↓  
Execution Systems

Side integrations:

* Policy Engine (OPA)  
* Identity Provider (OIDC)  
* Append-only Artifact Store

---

# **9\. Policy Integration**

Policy evaluation may be implemented using Open Policy Agent.

Policies:

* authored in Rego  
* versioned independently  
* evaluated during authority checkpoints

Gantral records `policy_version_id` within artifacts.

---

# **10\. Storage Model**

Execution indices: PostgreSQL (non-authoritative)

Commitment artifacts: append-only object storage (authoritative)

Artifact store is authoritative for replay.

---

# **11\. Security Architecture**

Identity: OAuth2 / OIDC federation  
Authorization: policy-driven  
Secrets: external secret managers

Gantral never stores raw secrets.

---

# **12\. Unified Authority Visibility**

Gantral exposes:

* running workflows  
* paused workflows awaiting authority  
* authority progression history

Gantral does not provide dashboards or analytics.

---

# **13\. Testing & Constitutional Enforcement**

Testing must enforce:

* transition correctness  
* atomic artifact emission  
* hash integrity  
* replay determinism  
* fail-closed behavior

---

# **14\. Auditor Considerations**

Gantral provides:

* cryptographically bound authority transitions  
* version-bound policy evaluation  
* identity validation at decision time  
* deterministic replay independent of logs

Artifacts are:

* hash-chained  
* digitally signed  
* timestamp-anchored  
* capable of long-term attestations

This enables durable authority evidence.

---

# **15\. Foundational Principle**

Gantral is not about what AI can do.

It is about **what organizations allow AI to do — and how that authority is structurally enforced and provably replayable.**

---

# **Appendix A — Artifact Specification v1 (Canonical)**

This appendix defines the **canonical structure for Gantral Commitment Artifacts**.

All implementations MUST conform exactly.

artifact\_hash uniquely identifies the artifact and MAY be used as the storage key in append-only artifact stores.

---

## **A.1 Artifact Purpose**

Artifacts represent cryptographically verifiable authority decisions.

They provide:

* ordering  
* identity attribution  
* policy binding  
* workflow binding  
* tamper evidence  
* replay verifiability

---

## **A.2 Artifact Structure**

The following fields constitute the canonical artifact schema.  
Implementations MUST include all fields unless explicitly marked optional.

artifact\_version  
artifact\_id

instance\_id  
workflow\_version\_id  
policy\_version\_id

authority\_state

context\_snapshot\_hash

human\_actor\_id  
justification

prev\_artifact\_hash  
artifact\_hash

crypto\_profile

artifact\_signature  
signature\_algorithm

timestamp\_token  
timestamp\_algorithm

attestations\[\] (optional)

---

## **A.3 Hash Model**

The artifact hash commits to the payload and prior artifact hash.

artifact\_hash \=  
H(payload || prev\_artifact\_hash)

Payload excludes:

* artifact\_hash  
* artifact\_signature  
* signature\_algorithm  
* timestamp\_token  
* timestamp\_algorithm  
* attestations

The payload includes all artifact fields except those explicitly excluded above.

---

## **A.4 Signature**

Artifacts are signed by the Gantral control plane.

artifact\_signature \= Sign(private\_key, artifact\_hash)

This proves authorship.

---

## **A.5 Timestamp Token**

Artifacts include a trusted timestamp token proving existence time.

Timestamp tokens MUST bind to `artifact_hash`.

---

## **A.6 Attestations**

Artifacts MAY include additional attestations.

These allow cryptographic upgrades over time without modifying the original artifact.

Example:

* new signature with stronger algorithm  
* later timestamp authority confirmation

Attestation entries are treated as opaque verification extensions.

Replay verification MUST NOT depend on the internal structure of attestations.

---

## **A.7 Artifact Immutability**

Artifacts MUST be stored in append-only storage.

Mutation invalidates replay.

---

## **A.8 Backward Compatibility**

Once Artifact Specification v1 is published:

* fields cannot be renamed  
* semantics cannot change  
* new versions must increment the specification version

---

## **A.9 Canonical Serialization Rule**

Before computing `artifact_hash`, the artifact payload MUST be serialized using **canonical JSON encoding**.

All implementations MUST follow the same serialization rules.

Canonical serialization ensures that independent implementations in Go, Python, Rust, or other languages produce identical byte sequences prior to hashing.

---

## **A.10 Canonical Encoding Requirements**

The canonical encoding MUST satisfy the following rules.

### **1\. UTF-8 Encoding**

Artifacts MUST be serialized as UTF-8.

No alternative encodings are permitted.

---

### **2\. Deterministic Field Ordering**

All JSON object keys MUST be sorted in **lexicographic ascending order**.

Example:

Correct:

\{  
  "authority\_state": "...",  
  "context\_snapshot\_hash": "...",  
  "instance\_id": "...",  
  "workflow\_version\_id": "..."  
\}

Incorrect:

\{  
  "instance\_id": "...",  
  "authority\_state": "...",  
  "workflow\_version\_id": "...",  
  "context\_snapshot\_hash": "..."  
\}

Sorting ensures deterministic serialization across languages.

---

### **3\. No Whitespace Normalization Variance**

Serialization MUST NOT include:

* trailing spaces  
* indentation  
* formatting differences

The payload MUST be serialized as **compact JSON**.

Example:

Correct:

\{"instance\_id":"...","workflow\_version\_id":"..."\}

---

### **4\. Stable Field Inclusion**

Fields MUST NOT be omitted if defined by the schema.

Optional fields MUST appear explicitly as `null` when absent.

Example:

"attestations": null

This prevents differences in hashing between implementations.

---

### **5\. String Normalization**

All strings MUST be UTF-8 encoded without normalization transformations.

Implementations MUST NOT apply locale transformations.

---

### **6\. Boolean Representation**

Booleans MUST use lowercase JSON values:

true  
false

---

### **7\. Numeric Representation**

If numeric fields exist in future versions:

* they MUST be encoded in base-10  
* exponential notation MUST NOT be used

---

## **A.11 Artifact Hash Computation**

After canonical serialization, the hash is computed.

artifact\_hash \= SHA256(canonical\_json\_bytes)

If the artifact has a predecessor:

artifact\_hash \=  
SHA256(canonical\_json\_bytes || prev\_artifact\_hash)

---

## **A.12 Verification Rule**

A verifier MUST:

1. deserialize artifact  
2. reconstruct canonical JSON  
3. recompute artifact\_hash  
4. verify signature  
5. verify hash chain linkage

Verification MUST NOT depend on:

* runtime logs  
* database records  
* workflow runtime state

---

## **Appendix A.13 — Normative Verification Procedure**

An independent verifier MUST validate artifacts using the following procedure.

1\. Deserialize the artifact.

2\. Validate schema completeness.  
   All fields defined in Artifact Specification v1 MUST be present.

3\. Reconstruct canonical JSON payload  
   using the canonical serialization rules defined in Section A.10.

4\. Recompute artifact\_hash:  
      hash \= SHA256(canonical\_payload\_bytes)

5\. If prev\_artifact\_hash is present:  
      hash \= SHA256(canonical\_payload\_bytes || prev\_artifact\_hash)

6\. Verify that recomputed hash \== artifact\_hash.

7\. Verify artifact\_signature using the public key corresponding to the  
   Gantral control plane and the declared signature\_algorithm.

8\. Verify timestamp\_token according to timestamp\_algorithm.

9\. Validate state transition correctness using the canonical  
   transition relation defined in Section 4.2.

10\. Repeat steps 1–9 for each artifact in the chain.

11\. Classify replay result:

   VALID  
     All checks pass and the final artifact represents a terminal  
     authority decision.

   INVALID  
     Any integrity, signature, or transition validation fails.

   INCONCLUSIVE  
     Chain integrity holds but the final artifact represents a  
     non-terminal authority state.

---

