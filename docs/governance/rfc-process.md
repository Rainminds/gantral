---
title: RFC Process
---

# RFC Process

Gantral evolves through an explicit Request for Comments (RFC) process.

The RFC process exists to balance:

- Community input
- Architectural coherence
- Execution velocity

Not all ideas are equal.  
Not all changes belong in the core.

---

## When an RFC Is Required

An RFC is required for changes that affect:

- Execution semantics
- HITL behavior
- Policy evaluation
- Audit guarantees
- Public APIs or specifications
- Architectural invariants

Minor bug fixes and documentation changes do not require RFCs.

---

## RFC Lifecycle

1. **Proposal**
   - Author submits RFC as a Markdown document
   - Motivation, design, and trade-offs are explicit

2. **Discussion**
   - Maintainers and community review
   - Feedback is incorporated or rejected with rationale

3. **Decision**
   - Maintainers accept, request revision, or reject
   - Decision is recorded publicly

4. **Implementation**
   - Accepted RFCs guide implementation
   - Specs precede code

5. **Stability**
   - Once implemented, semantics are considered stable
   - Breaking changes require new RFCs

---

## Decision Authority

Final decision authority currently rests with the core maintainers, subject to future governance evolution.

This is intentional.

Gantral is infrastructure; consensus-by-committee is not viable for execution semantics.

Community input is critical.  
Architectural coherence is non-negotiable.

---

## RFC Design Principles

RFCs should:

- Be explicit about non-goals
- Respect architectural invariants
- Prefer simplicity over completeness
- Avoid coupling to specific vendors or tools
- Consider regulatory and audit implications

An RFC that improves convenience but weakens control will not be accepted.
