---
name: sdd-shared-conventions
description: >
  Shared convention and reference files used by all SDD phase skills.
  This is not an executable skill — it provides persistence contracts,
  engram conventions, openspec conventions, and phase-common rules
  that other SDD skills reference.
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
---

# SDD Shared Conventions

This directory contains reference files shared across all SDD phase skills:

- `persistence-contract.md` — Artifact store contract (engram/openspec/hybrid/none)
- `engram-convention.md` — Engram topic key naming and upsert conventions
- `openspec-convention.md` — OpenSpec file layout and state management
- `sdd-phase-common.md` — Common rules for all SDD executor phase agents
