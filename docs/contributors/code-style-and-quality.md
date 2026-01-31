---
title: Code Style and Quality
---

# Code Style and Quality

Gantral prioritizes clarity and auditability over cleverness.

Code should be written to be read, reviewed, and reasoned about by humans.

---

## General Guidelines

- Prefer explicit over implicit behavior
- Avoid hidden side effects
- Use descriptive naming
- Document state transitions clearly
- Minimize global state
- Write tests for execution logic

---

## Auditability Considerations

Code should make it easy to determine:

- What state transition occurred
- Why it occurred
- Under what authority
- With what inputs

If this is not clear from reading the code, it should be revised.

---

## Error Handling

Error paths should be explicit and observable.

Silent failure or implicit retries are discouraged unless explicitly justified.

---

## Legal Notice

Code contributions are provided “as is” under the project license.

Contributors are responsible for ensuring that they have the right to submit the code they contribute.
