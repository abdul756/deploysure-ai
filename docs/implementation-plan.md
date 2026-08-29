# Implementation Plan

## 1. Overview

This document tracks the phased implementation of DeploySure AI. Each phase
has a clear deliverable, acceptance criteria, and the Bob Agent capabilities
it exercises.

---

## 2. Phases

### Phase 0 — Project Scaffold ✅

**Goal**: Establish the repository structure, documentation baseline, and
contribution conventions.

**Deliverables**:
- `README.md` with project overview and directory layout.
- `PROBLEM_STATEMENT.md` and `SOLUTION_STATEMENT.md`.
- `docs/requirements.md`, `docs/release-runbook.md`, `docs/architecture.md`,
  `docs/implementation-plan.md`.
- `.env.example` with placeholder values.
- `SECURITY.md` with credential management guidelines.

**Bob Capabilities Exercised**:
- Document understanding (reading and interpreting requirements).
- Agent mode (orchestrating the scaffold via a single Bob task).

**Acceptance Criteria**:
- All documents exist and are internally consistent.
- No credentials committed.
- `git status` is clean after initial commit.

---

### Phase 1 — Sample Application

**Goal**: Build the synthetic Go service that will be the subject of all
subsequent analysis.

**Deliverables**:

| File | Description |
|------|-------------|
| `sample-app/cmd/server/main.go` | Entry point; wires router, config, graceful shutdown |
| `sample-app/handlers/orders.go` | `GET /api/orders` handler returning a JSON array |
| `sample-app/handlers/health.go` | `GET /health` and `GET /ready` handlers |
| `sample-app/config/config.go` | Reads `PORT`, `LOG_LEVEL`, `DB_DSN` from environment with defaults |
| `sample-app/handlers/orders_test.go` | Unit tests for orders handler |
| `sample-app/handlers/health_test.go` | Unit tests for health and ready handlers |
| `sample-app/Dockerfile` | Multi-stage build; distroless final image; non-root user |
| `sample-app/k8s/deployment.yaml` | Kubernetes Deployment with probes, resource limits, security context |
| `sample-app/k8s/service.yaml` | Kubernetes Service (ClusterIP) |

**Bob Capabilities Exercised**:
- Code generation (Go `net/http`, Dockerfile, Kubernetes YAML).
- Docs-to-code consistency (generated code must match FR-30 through FR-38).

**Acceptance Criteria**:
- `go build ./...` succeeds.
- `go test ./...` passes with ≥ 80 % coverage.
- `gofmt -l ./...` returns empty output.
- `go vet ./...` exits 0.
- `docker build` succeeds.
- Container serves `/health` as non-root user.

---

### Phase 2 — DeploySure Analysis Engine

**Goal**: Implement the five specialised subagents and the orchestrating Bob
Agent task.

**Deliverables**:

| Component | Description |
|-----------|-------------|
| Bob task definition | Prompt/task that reads docs and spawns all subagents |
| Code-quality subagent | Runs `gofmt`, `go vet`; reports findings as structured JSON |
| Test-gap subagent | Enumerates exported symbols; cross-references `*_test.go` files |
| Docker review subagent | Audits `Dockerfile` against runbook checks D-01 through D-06 |
| Kubernetes review subagent | Audits `k8s/*.yaml` against runbook checks K-01 through K-07 |
| Docs-to-code subagent | Compares routes and env vars in docs vs source |
| Aggregator | Merges all subagent outputs into a single `findings[]` JSON array |

**Bob Capabilities Exercised**:
- Parallel tasks (`spawn_subagent` with independent contexts).
- Specialised subagents (each with a narrow, focused scope).
- Code-quality analysis, unit-test gap analysis, Docker review,
  Kubernetes review, docs-to-code consistency analysis.

**Acceptance Criteria**:
- All five subagents complete successfully and return structured findings.
- Aggregated `findings[]` array contains at least the seeded defects from
  Phase 1's intentionally misconfigured artefacts.

---

### Phase 3 — watsonx.ai Granite Integration

**Goal**: Send aggregated findings to Granite for risk prioritisation.

**Deliverables**:

| File | Description |
|------|-------------|
| `backend/watsonx/client.go` | Go client for the watsonx.ai Inference API |
| `backend/watsonx/prioritise.go` | Constructs the Granite prompt; parses the scored response |
| `reports/latest.md` | Template for the generated risk report |

**Bob Capabilities Exercised**:
- watsonx.ai Granite risk prioritisation.
- Document understanding (interpreting Granite's natural-language response).

**Acceptance Criteria**:
- Each finding in `reports/latest.md` has a `severity` field
  (CRITICAL / HIGH / MEDIUM / LOW / INFO).
- Findings are sorted in descending severity order.
- The Granite API call uses credentials from `.env`, not hardcoded.

---

### Phase 4 — Human-Approved Remediation

**Goal**: Implement the approval gate and automated fix application.

**Deliverables**:

| Component | Description |
|-----------|-------------|
| Bob approval prompt | Presents `reports/latest.md`; pauses for human `yes/no` per finding |
| Remediation applicator | Bob Agent applies approved fixes (edit files, update manifests) |
| Remediation log | `reports/remediation-<timestamp>.md` mapping finding IDs to applied changes |

**Bob Capabilities Exercised**:
- Human-approved remediation (`ask_followup_question` gate before every fix).
- Agent mode (applying `apply_diff` and `search_and_replace` edits).

**Acceptance Criteria**:
- No fix is applied without an explicit human approval in the Bob session.
- The remediation log records every approved and rejected finding.

---

### Phase 5 — Final Validation & Reporting

**Goal**: Re-run all checks post-remediation and produce the final evidence
package.

**Deliverables**:

| Component | Description |
|-----------|-------------|
| Final validation pass | All five subagents re-run; zero CRITICAL/HIGH findings required |
| `reports/final-<timestamp>.md` | Clean validation report |
| `bob_sessions/<session>.json` | Exported Bob session covering Phases 1–5 |
| `evidence/` | Screenshots of key workflow steps |

**Bob Capabilities Exercised**:
- Final validation (structured re-run of all checks).
- Full end-to-end demonstration of all twelve listed capabilities.

**Acceptance Criteria**:
- Final report contains zero CRITICAL or HIGH severity findings.
- Bob session export is present and readable.
- All evidence files are committed.

---

## 3. Seeded Defects (Phases 1 → 2)

To make the analysis demonstration meaningful, Phase 1 will intentionally
introduce the following defects into the sample application. Phase 2 must
detect all of them.

| ID | Artefact | Defect | Maps to Runbook Check |
|----|----------|--------|----------------------|
| SD-01 | `Dockerfile` | `FROM golang:latest` — floating tag | D-02 |
| SD-02 | `Dockerfile` | No `USER` directive — runs as root | D-01 |
| SD-03 | `k8s/deployment.yaml` | Missing `resources.limits` | K-02 |
| SD-04 | `k8s/deployment.yaml` | Missing `livenessProbe` | K-03 |
| SD-05 | `k8s/deployment.yaml` | Missing `readinessProbe` | K-04 |
| SD-06 | `k8s/deployment.yaml` | Missing `runAsNonRoot` | K-05 |
| SD-07 | `handlers/orders.go` | Exported `Helper` function without unit test | FR-09 |
| SD-08 | `config/config.go` | `DB_DSN` env var read in code but not in `docs/requirements.md` | DC-02 |

---

## 4. Timeline

| Phase | Estimated Effort | Dependency |
|-------|-----------------|------------|
| 0 — Scaffold | 1 hour | None |
| 1 — Sample App | 2 hours | Phase 0 |
| 2 — Analysis Engine | 3 hours | Phase 1 |
| 3 — Granite Integration | 2 hours | Phase 2 |
| 4 — Remediation | 1 hour | Phase 3 |
| 5 — Final Validation | 1 hour | Phase 4 |
| **Total** | **~10 hours** | |

---

## 5. Open Questions

| # | Question | Owner | Status |
|---|----------|-------|--------|
| OQ-01 | Which Granite model version is available in the hackathon environment? | Team | Open |
| OQ-02 | Is `kubectl` available in the demo environment for dry-run validation? | Team | Open |
| OQ-03 | Should the web frontend (Phase 6+) be in scope for the hackathon submission? | Team | Open |
