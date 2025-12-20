# 1. Use Architectural Decision Records

*   **Status**: Accepted
*   **Date**: 2025-12-18
*   **Deciders**: Rainminds Engineering

## Context
We need to track architectural decisions for Gantral. As an open-core project with potential for distributed contributors, relying on Slack history or verbal agreements leads to lost context and "tribal knowledge."

## Decision
We will use Architectural Decision Records (ADRs), specifically the Nygard format, stored in `specs/adr/`.

## Consequences
*   **Positive:** We have a version-controlled history of "why" decisions were made. New team members (and AI agents) can read the history to understand the evolution.
*   **Negative:** Requires discipline to write them.
