# Gantral – Commitment Artifact & Offline Verifier Implementation Specification

Version: v1.0 (Verifier-Executable)  
Status: Normative implementation specification

Audience:  
• AI coding assistants implementing Gantral  
• Core maintainers reviewing correctness  
• Independent auditors and verifiers  
• Enterprise platform and security teams

This document specifies the exact structure, semantics, and verification logic  
for Gantral commitment artifacts and offline verification tooling.

## **0\. Scope and Authority**

This document is a normative specification.

If an implementation deviates from this document, it is incorrect.  
If this document conflicts with the TRD, the TRD is authoritative.

## **1\. Purpose of the Commitment Artifact**

The commitment artifact is the sole authoritative proof of execution authority.

It allows any third party to verify:  
• What authority decision occurred  
• For which execution instance  
• Under which policy and execution context  
• Without trusting Gantral operators, databases, logs, or dashboards

## **2\. Artifact Design Principles**

• Immutable  
• Append-only  
• Hash-chained  
• Deterministic  
• Log-independent  
• Verifiable offline  
• Fail-closed by default

## **3\. Artifact Serialization Format**

Primary format: Canonical JSON  
Alternate format: CBOR (binary equivalent)

Canonical JSON rules:  
• UTF-8 encoding  
• Sorted keys  
• No insignificant whitespace  
• Deterministic field ordering

## **4\. Commitment Artifact Schema (Normative)**

artifact\_version: string (semver)  
artifact\_id: UUID  
instance\_id: UUID  
sequence\_number: integer (monotonic)  
prev\_artifact\_hash: string (hex, nullable for first artifact)  
authority\_state: enum {APPROVED, REJECTED, OVERRIDDEN}  
policy\_version\_id: string  
execution\_context\_hash: string (hex)  
human\_actor\_id: string  
human\_role: string  
decision\_timestamp: RFC3339 timestamp  
justification\_hash: string (hex)  
artifact\_hash: string (hex)

## **5\. Hash Construction (Deterministic)**

artifact\_hash is computed as a hash over the concatenation of:

artifact\_version |  
artifact\_id |  
instance\_id |  
sequence\_number |  
prev\_artifact\_hash |  
authority\_state |  
policy\_version\_id |  
execution\_context\_hash |  
human\_actor\_id |  
human\_role |  
decision\_timestamp |  
justification\_hash

Any mutation invalidates the artifact.

## **6\. Atomic Emission Requirements**

Artifact emission MUST be atomic with authority transition.

If any of the following fail, execution MUST terminate:  
• Authority state transition  
• Artifact serialization  
• Artifact persistence  
• Hash computation

No retries. No partial success.

## **7\. Artifact Storage Requirements**

• Write-once storage  
• Append-only semantics  
• Database records are non-authoritative indices only  
• Artifacts must not be reconstructible from logs

## **8\. Offline Verifier Specification**

Verifier characteristics:  
• Standalone binary or library  
• No network access  
• Input: artifact(s)  
• Output: VALID / INVALID / INCONCLUSIVE

Verifier responsibilities:  
• Validate schema  
• Recompute hashes  
• Validate hash chain  
• Validate sequence ordering

## **9\. Verifier Outcome Semantics**

VALID:  
• All hashes match  
• Chain intact  
• Schema valid

INVALID:  
• Hash mismatch  
• Broken chain  
• Illegal authority state  
• Schema violation

INCONCLUSIVE:  
• Missing artifact(s)  
• Partial chain  
• Unsupported artifact\_version

## **10\. Replay Semantics**

Replay is authority-only.

Verifier MUST NOT:  
• Replay agent logic  
• Inspect logs  
• Infer missing authority

Replay input \= artifact chain only.

## **11\. Failure Handling Guarantees**

Any ambiguity results in INVALID or INCONCLUSIVE.

The verifier MUST NOT guess or infer intent.

## **12\. Versioning and Compatibility**

• artifact\_version follows semver  
• Unknown major versions MUST be rejected  
• Minor versions must be backward compatible

## **13\. Security Considerations**

• Artifacts contain no secrets  
• PII must be hashed  
• Verifier assumes hostile environment

## **14\. Non-Goals**

This specification does not:  
• Select cryptographic algorithms  
• Make legal admissibility claims  
• Define policy semantics

## **15\. Final Guarantee**

If a valid artifact exists per this specification,  
then the authority decision is provable independent of Gantral.

If no valid artifact exists, authority did not occur.  
