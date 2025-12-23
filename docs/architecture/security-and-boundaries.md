# Security and Trust Boundaries

This document outlines security-related design intent and trust boundaries.

---

## Trust Model

Gantral assumes:

- Agents may be untrusted or partially trusted
- Human actors are authenticated externally
- Enterprise systems enforce their own access controls

Gantral does not replace identity or access management systems.

---

## Authority Boundaries

Gantral is designed to:

- Record who made a decision
- Record under what role or context
- Avoid implicit authority escalation

Authorization mechanisms are policy-driven and configurable.

---

## Security Responsibilities

Gantral is designed to provide execution control mechanisms.

Organizations are responsible for:

- Deployment security
- Identity management
- Network controls
- Regulatory compliance

---

## Security Limitations

Gantral does not:

- Guarantee security outcomes
- Prevent misuse by authorized actors
- Detect all malicious behavior

It provides infrastructure intended to support secure operation when deployed appropriately.
