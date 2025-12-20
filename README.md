# Gantral

**The AI Execution Control Plane.**

Gantral is an open-source protocol for governing AI agents in the enterprise. It solves the "Who authorized this?" problem by providing a deterministic execution engine, immutable audit logs, and a first-class Human-in-the-Loop (HITL) state machine.

> **Status:** Initialization / Pre-Alpha.

## 📚 Documentation

The technical constitution of Gantral lives in the `specs/` directory. These documents are the Single Source of Truth.

- **[Technical Specifications](specs/README.md)**: The complete technical reference.
- **[Architecture](specs/01-architecture.md)**: High-level design and invariants.
- **[Domain Model](specs/02-domain-model.md)**: Workflows, Instances, and Decisions.
- **[State Machine](specs/03-state-machine.md)**: The lifecycle of an execution.

## 🏛️ Governance & Community

Gantral is a "Maintainer-Led" project committed to transparency and community collaboration.

- **[Governance](GOVERNANCE.md)**: How decisions are made.
- **[Contributing](CONTRIBUTING.md)**: How to get involved.
- **[Code of Conduct](CODE_OF_CONDUCT.md)**: Our pledge for a safe community.
- **[Security](SECURITY.md)**: Reporting vulnerabilities.

## 🛠️ Usage

To build the core engine:

```bash
make build
```

To run tests:

```bash
make test
```

## ⚖️ License

Apache 2.0. See [LICENSE](LICENSE) for details.
