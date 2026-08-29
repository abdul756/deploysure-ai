# Create the project documentation for DeploySure AI.Create the project documentation for DeploySure AI.

Create:

- docs/requirements.md
- docs/release-runbook.md
- docs/architecture.md
- docs/implementation-plan.md
- PROBLEM_STATEMENT.md
- SOLUTION_STATEMENT.md

DeploySure improves the developer release-readiness workflow.

The solution must demonstrate:

- IBM Bob Agent mode
- Document understanding
- Parallel tasks
- Specialized subagents
- Code-quality analysis
- Unit-test gap analysis
- Docker review
- Kubernetes review
- Documentation-to-code consistency analysis
- Human-approved remediation
- Final validation
- watsonx.ai Granite risk prioritization

The synthetic sample application must:

- use Go and net/http,
- expose GET /api/orders,
- expose GET /health,
- expose GET /ready,
- use environment-based configuration,
- support graceful shutdown,
- include unit tests,
- include a Dockerfile,
- include Kubernetes Deployment and Service manifests.

The release runbook must require:

- gofmt,
- go vet,
- go test,
- health and readiness endpoints,
- non-root container execution,
- liveness and readiness probes,
- resource requests and limits,
- no committed credentials.

Do not implement code yet.

---

**Status:** active  **Date:** 2026-08-29

---

### 👤 User

Create the project documentation for DeploySure AI.Create the project documentation for DeploySure AI.

Create:

- docs/requirements.md
- docs/release-runbook.md
- docs/architecture.md
- docs/implementation-plan.md
- PROBLEM_STATEMENT.md
- SOLUTION_STATEMENT.md

DeploySure improves the developer release-readiness workflow.

The solution must demonstrate:

- IBM Bob Agent mode
- Document understanding
- Parallel tasks
- Specialized subagents
- Code-quality analysis
- Unit-test gap analysis
- Docker review
- Kubernetes review
- Documentation-to-code consistency analysis
- Human-approved remediation
- Final validation
- watsonx.ai Granite risk prioritization

The synthetic sample application must:

- use Go and net/http,
- expose GET /api/orders,
- expose GET /health,
- expose GET /ready,
- use environment-based configuration,
- support graceful shutdown,
- include unit tests,
- include a Dockerfile,
- include Kubernetes Deployment and Service manifests.

The release runbook must require:

- gofmt,
- go vet,
- go test,
- health and readiness endpoints,
- non-root container execution,
- liveness and readiness probes,
- resource requests and limits,
- no committed credentials.

Do not implement code yet.

---

### 🤖 Assistant

> 🔧 called: read_file, update_todo_list



> 🔧 **read_file** — README.md

> 🔧 **update_todo_list** — 
[ ] Create docs/requirements.md
[ ] Create docs/release-runbook.md
[ ] Create docs/architecture.md
[ ] Create docs/implementation-plan.md
[ ] Create PROBLEM_STATEMENT.md
[ ] Create SOLUTION_STATEMENT.md


### 🤖 Assistant

> 🔧 called: write_file, write_file, write_file, write_file

Good — I have the full project context. Now I'll write all six documents in parallel.

> 🔧 **write_file** — PROBLEM_STATEMENT.md

> 🔧 **write_file** — docs/requirements.md

> 🔧 **write_file** — docs/release-runbook.md

> 🔧 **write_file** — SOLUTION_STATEMENT.md

### 🤖 Assistant

> 🔧 called: write_file, write_file

Four files written. Now the architecture and implementation plan:

> 🔧 **write_file** — # Implementation Plan

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


> 🔧 **write_file** — # Architecture

## 1. System Overview

DeploySure AI is composed of three logical layers:

1. **Orchestration layer** — IBM Bob Agent mode drives the entire workflow.
2. **Analysis layer** — five specialised subagents, each responsible for one
   review domain, running in parallel.
3. **Intelligence layer** — watsonx.ai Granite 13B receives the aggregated
   findings and returns severity scores and a ranked remediation order.

```
┌────────────────────────────────────────────────────────────┐
│                   Bob Agent (Orchestrator)                  │
│                                                            │
│  1. Reads docs/requirements.md + docs/release-runbook.md  │
│  2. Spawns five subagents in parallel                      │
│  3. Aggregates findings                                    │
│  4. Calls watsonx.ai Granite                               │
│  5. Presents report + awaits human approval                │
│  6. Applies approved remediations                          │
│  7. Runs final validation pass                             │
└───────────────────────┬────────────────────────────────────┘
                        │ spawn_subagent (parallel)
        ┌───────────────┼────────────────────┐
        │               │                    │
        ▼               ▼                    ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────────┐
│  Code Quality│ │  Test Gap    │ │  Docker Review   │
│  Subagent    │ │  Subagent    │ │  Subagent        │
│              │ │              │ │                  │
│ gofmt        │ │ exported fns │ │ USER directive   │
│ go vet       │ │ vs *_test.go │ │ base image tag   │
│ GoDoc        │ │ coverage %   │ │ secrets in layers│
└──────────────┘ └──────────────┘ └──────────────────┘
        │               │                    │
        ▼               ▼                    ▼
┌──────────────────┐           ┌──────────────────────┐
│  Kubernetes      │           │  Docs-to-Code        │
│  Review Subagent │           │  Consistency Subagent│
│                  │           │                      │
│ resource limits  │           │ routes: docs vs code │
│ probes           │           │ env vars: docs vs code│
│ runAsNonRoot     │           │ discrepancy report   │
│ readOnlyRootFS   │           └──────────────────────┘
└──────────────────┘
        │
        ▼
┌────────────────────────────────┐
│  Aggregated Findings (JSON)    │
└───────────────┬────────────────┘
                │
                ▼
┌────────────────────────────────┐
│  watsonx.ai Granite 13B        │
│  • severity scoring            │
│  • ranked remediation order    │
└───────────────┬────────────────┘
                │
                ▼
┌────────────────────────────────┐
│  Prioritised Risk Report       │
│  reports/latest.md             │
└───────────────┬────────────────┘
                │  human approval required
                ▼
┌────────────────────────────────┐
│  Remediation                   │
│  (Bob applies approved fixes)  │
└───────────────┬────────────────┘
                │
                ▼
┌────────────────────────────────┐
│  Final Validation Pass         │
│  (all five subagents re-run)   │
└────────────────────────────────┘
```

---

## 2. Component Details

### 2.1 Bob Agent (Orchestrator)

- **Mode**: IBM Bob Agent mode.
- **Entry point**: A single Bob task that manages the full lifecycle.
- **Responsibilities**: document ingestion, subagent spawning, result
  aggregation, Granite invocation, human-approval gate, remediation
  coordination, final validation.
- **Persistence**: the complete session is exported to `bob_sessions/`.

### 2.2 Code-Quality Subagent

- **Scope**: `sample-app/**/*.go` (excludes `*_test.go`).
- **Checks**: `gofmt -l`, `go vet`, GoDoc coverage.
- **Output schema**:
  ```json
  { "check": "gofmt", "file": "cmd/server/main.go", "line": 42, "message": "..." }
  ```

### 2.3 Test-Gap Subagent

- **Scope**: `sample-app/**/*.go` and `sample-app/**/*_test.go`.
- **Checks**: enumerate exported functions and HTTP handlers; cross-reference
  against test files; compute per-package coverage.
- **Output schema**:
  ```json
  { "symbol": "OrdersHandler", "file": "handlers/orders.go", "tested": false }
  ```

### 2.4 Docker Review Subagent

- **Scope**: `sample-app/Dockerfile`.
- **Checks**: USER directive, base image tag, credential patterns, layer count.
- **Output schema**:
  ```json
  { "check": "non-root-user", "line": 1, "severity": "CRITICAL", "message": "..." }
  ```

### 2.5 Kubernetes Review Subagent

- **Scope**: `sample-app/k8s/*.yaml`.
- **Checks**: resource requests/limits, liveness/readiness probes,
  `runAsNonRoot`, `readOnlyRootFilesystem`.
- **Output schema**:
  ```json
  { "check": "readinessProbe", "resource": "Deployment/orders-api", "severity": "HIGH", "message": "..." }
  ```

### 2.6 Docs-to-Code Consistency Subagent

- **Scope**: `docs/requirements.md` and `sample-app/**/*.go`.
- **Checks**: route table vs `http.HandleFunc` registrations; env-var table
  vs `os.Getenv` / `os.LookupEnv` calls.
- **Output schema**:
  ```json
  { "check": "undocumented-route", "route": "GET /api/orders", "source_file": "handlers/orders.go:12" }
  ```

### 2.7 watsonx.ai Granite Integration

- **Model**: `ibm/granite-13b-instruct-v2` (or latest available).
- **Input**: concatenated JSON array of all subagent findings.
- **Prompt contract**: structured prompt requesting severity classification
  (CRITICAL / HIGH / MEDIUM / LOW / INFO) and ranked remediation list.
- **Output**: severity-annotated findings array, stored at
  `reports/latest.md`.

---

## 3. Data Flow

```
docs/requirements.md ─────────────────────────────────────────┐
docs/release-runbook.md ──────────────────────────────────────┤
sample-app/**/*.go ───────────────────────────────────────────┤
sample-app/Dockerfile ────────────────────────────────────────┤──► Bob Agent
sample-app/k8s/*.yaml ────────────────────────────────────────┘
                                   │
                   ┌───────────────▼───────────────────────────┐
                   │           Five subagents                   │
                   └───────────────┬───────────────────────────┘
                                   │ findings[]
                                   ▼
                           watsonx.ai Granite
                                   │ scored + ranked findings[]
                                   ▼
                          reports/latest.md
                                   │ human approval
                                   ▼
                           Bob applies fixes
                                   │
                                   ▼
                          Final validation
```

---

## 4. Directory Layout

```
deploysure-ai/
├── frontend/               # Plain HTML/CSS/JS web interface
├── backend/                # Go backend — REST API and Watsonx integration
├── sample-app/             # Synthetic Go service (the subject of analysis)
│   ├── cmd/server/         # main.go — entry point
│   ├── handlers/           # HTTP handlers (orders, health, ready)
│   ├── config/             # Environment-based configuration
│   ├── Dockerfile          # Container build definition
│   └── k8s/                # Kubernetes Deployment and Service manifests
├── docs/
│   ├── requirements.md     # This document's sibling — functional requirements
│   ├── release-runbook.md  # Mandatory pre-release checks
│   ├── architecture.md     # This document
│   └── implementation-plan.md
├── reports/                # Generated risk-analysis reports (gitignored output)
├── bob_sessions/           # Exported Bob AI session logs
└── evidence/               # Demo screenshots and supporting evidence
```

---

## 5. Security Considerations

| Concern | Mitigation |
|---------|-----------|
| Credentials in env vars | `.env` is gitignored; `.env.example` contains only placeholder values |
| Secrets in Docker layers | Docker Review subagent flags credential patterns; multi-stage build minimises layer surface |
| Granite API key | Stored in `.env`; never logged or committed |
| Bob session exports | Reviewed before commit to ensure no secrets are captured in tool outputs |


### 🤖 Assistant

> 🔧 called: execute_command, update_todo_list

All six files are written. Let me do a quick sanity check and mark everything complete.

> 🔧 **execute_command** — ls -1 PROBLEM_STATEMENT.md SOLUTION_STATEMENT.md docs/requirements.md docs/release-runbook.md docs/architecture.md docs/implementation-plan.md && echo "--- line counts ---" && wc -l PROBLEM_STATEMENT.md SOLUTION_STATEMENT.md docs/requirements.md docs/release-runbook.md docs/architecture.md docs/implementation-plan.md

> 🔧 **update_todo_list** — 
[x] Create docs/requirements.md
[x] Create docs/release-runbook.md
[x] Create docs/architecture.md
[x] Create docs/implementation-plan.md
[x] Create PROBLEM_STATEMENT.md
[x] Create SOLUTION_STATEMENT.md


### 🤖 Assistant

All six documents are created and verified (828 lines total). Here's a summary of what was produced:

---

### Files created

| File | Lines | Purpose |
|------|------:|---------|
| [`PROBLEM_STATEMENT.md`](PROBLEM_STATEMENT.md) | 43 | Pain points, affected personas, consequence of inaction, scope boundary |
| [`SOLUTION_STATEMENT.md`](SOLUTION_STATEMENT.md) | 84 | ASCII workflow diagram, full capability table mapped to how DeploySure uses each one, value proposition |
| [`docs/requirements.md`](docs/requirements.md) | 118 | Structured FR/NFR tables — platform (FR-01→FR-29) and sample-app (FR-30→FR-38) requirements |
| [`docs/release-runbook.md`](docs/release-runbook.md) | 162 | Mandatory gates: `gofmt`, `go vet`, `go test`, health/ready smoke test, Docker gate (D-01→D-06), Kubernetes gate (K-01→K-07), security gate, docs-consistency gate, sign-off steps, rollback procedure |
| [`docs/architecture.md`](docs/architecture.md) | 211 | Three-layer architecture, ASCII component diagrams, per-subagent output schemas, data-flow diagram, directory layout, security considerations |
| [`docs/implementation-plan.md`](docs/implementation-plan.md) | 210 | Six phases (0 scaffold → 5 final validation), per-phase deliverable tables, **seeded defect table** (SD-01→SD-08) mapping each defect to its runbook check, timeline, open questions |

### Key design decisions captured

- **Seeded defects table** (Phase 1 → Phase 2): eight intentional defects covering root Dockerfile, floating image tag, missing K8s probes/limits, missing security context, uncovered test, and an undocumented env var — giving the analysis engine concrete targets.
- **Runbook is the authoritative standard**: every automated check in the architecture traces back to a named runbook item (D-xx, K-xx, FR-xx, DC-xx).
- **Human approval is gated**, not optional: the runbook, architecture, and implementation plan all reinforce that no fix is applied without explicit sign-off in the Bob session.