# BOB_USAGE.md

## How IBM Bob 2.0 Was Used in DeploySure AI

IBM Bob 2.0 was the core development and workflow-orchestration component of
DeploySure AI. This document records every session in which Bob was used,
covering all ten tasks from project scaffolding through final validation.

> **Security note:** No credentials or secrets were shared with Bob at any
> point. All prompts followed the safe-usage guidelines in
> [SECURITY.MD](SECURITY.MD).

---

## Product Boundary

IBM Bob IDE performs repository analysis, implementation, remediation and
validation.

The web dashboard presents Bob-generated reports and calls watsonx.ai Granite.
It does not invoke an undocumented Bob API.

---

## Human Oversight

Bob did not apply remediation until the generated plan was reviewed and
explicitly approved by a human operator. The remediation plan was created in
Task 08 and held until Task 09 was triggered by explicit approval.

---

## Features Demonstrated

| Feature | Tasks |
|---|---|
| Agent mode | 01, 09 |
| Document understanding | 02 |
| Parallel subagents | 04, 10 |
| Code generation | 03, 05, 07 |
| API integration (watsonx.ai) | 06 |
| Release analysis & planning | 08 |
| Human-approved remediation | 09 |
| Final validation | 10 |

---

## Session Log

### Task 01 — Project Scaffolding

**Bob capability:** Agent mode  
**Evidence:** `bob_sessions/task01_project_scaffolding_summary.md`  
**Screenshots:** `bob_sessions/deploysure_task01_project_scaffolding_summary.png`

Bob scaffolded the entire repository from the official IBM Hackathon GitHub
template. Actions taken:

- Created `backend/`, `sample-app/`, `frontend/`, and empty top-level
  directories with `.gitkeep` files.
- Generated `README.md`, `PROBLEM_STATEMENT.md`, `SOLUTION_STATEMENT.md`,
  `DEMO_SCRIPT.md`, `BOB_USAGE.md`, and `SECURITY.MD`.
- Produced `docs/architecture.md` and `docs/project-plan.md` capturing the
  six-phase delivery plan and full component architecture (Bob orchestrator,
  four specialized subagents, watsonx.ai Granite integration).
- Confirmed all security files were unmodified from the template.

---

### Task 02 — Document Understanding & Architecture

**Bob capability:** Document understanding  
**Evidence:** `bob_sessions/task02_document_understanding_summary.md`  
**Screenshots:** `bob_sessions/task02_document_understanding_summary.png`,
`bob_sessions/task02_document_understanding_prompt.png`

Bob read the project brief and produced two design documents:

- **`docs/project-plan.md`** — six-phase delivery plan (Phase 0 scaffold
  through Phase 5 final validation), including seeded defects and timeline.
- **`docs/architecture.md`** — full system overview covering all components:
  Bob Agent (orchestrator), Code-Quality subagent, Test-Gap subagent, Docker
  Review subagent, Kubernetes Review subagent, Docs-to-Code Consistency
  subagent, and watsonx.ai Granite Integration; data flow and directory layout.

Key design decisions captured:
- watsonx.ai Granite-3-8B-Instruct selected as the AI model.
- Subagent findings aggregated into `reports/findings-before.json`.
- Human-approved gate required before any remediation is applied.

---

### Task 03 — Sample Application Build

**Bob capability:** Code generation  
**Evidence:** `bob_sessions/task03_sample_application_summary.md`  
**Screenshots:** `bob_sessions/task03_sample_application_summary.png`,
`bob_sessions/task03_sample_application_build_prompt.png`

Bob generated the synthetic Go microservice under `sample-app/` with
intentionally seeded defects for later detection:

**Files created:**
- `sample-app/cmd/server/main.go`
- `sample-app/internal/handlers/readiness.go`
- `sample-app/internal/handlers/orders.go`
- `sample-app/internal/store/orders.go`
- `sample-app/go.mod`
- `sample-app/Dockerfile`
- `sample-app/k8s/deployment.yaml`
- `sample-app/k8s/service.yaml`

**Seeded defects (for Phase 2 detection):**

| ID | Defect | Location |
|---|---|---|
| SD-01 / DP-002 | Floating `:latest` tag on builder stage | `Dockerfile` |
| SD-02 / DP-001 | Missing `USER` directive — container runs as root | `Dockerfile` |
| SD-03 / DC-001 | Missing `GET /health` endpoint | `handlers/` |
| SD-04 / DC-003 | No pod-level `securityContext` (`runAsNonRoot`) | `k8s/deployment.yaml` |
| SD-05 / DC-004 | No `livenessProbe` | `k8s/deployment.yaml` |
| SD-06 / DC-005 | No resource `requests`/`limits` | `k8s/deployment.yaml` |
| SD-07 / DC-006 | No `readOnlyRootFilesystem` | `k8s/deployment.yaml` |
| SD-08 / CQ-001 | `json.Encoder` error silently ignored | `handlers/readiness.go` |
| SD-09 / CQ-003 | Invalid Go version in `go.mod` (`go 1.21.0` instead of `1.21`) | `go.mod` |

All validation tests passed after generation.

---

### Task 04 — Parallel Initial Analysis

**Bob capability:** Parallel subagents  
**Evidence:** `bob_sessions/task04_parallel_analysis_summary.md`  
**Screenshots:** `bob_sessions/task04_parallel_analysis_summary.png`

Bob spawned four independent specialized subagents simultaneously to analyse
`sample-app/`:

| Subagent | Findings | Key issues |
|---|---|---|
| Code Quality | 15 | Silent error, dead code, missing test coverage, invalid go.mod version |
| Test Coverage | 8 | 0% `cmd/server` coverage, missing error-branch tests, no store tests |
| Deployment (Docker + K8s) | 8 | Floating `:latest`, no `USER`, missing probes, no resource limits, no securityContext |
| Documentation Consistency | 10 | Undocumented endpoints, `README` gaps, missing runbook cross-references |

**Decision: ❌ BLOCKED — 5 blocker findings, not release-ready.**

Files produced:
- `reports/findings-before.json` (41 total findings)
- `reports/release-readiness-before.md`
- `reports/test-results-before.txt`

---

### Task 05 — Go Backend Implementation

**Bob capability:** Code generation and testing  
**Evidence:** `bob_sessions/task05_backend_summary.md`  
**Screenshots:** `bob_sessions/task05_backend_summary.png`

Bob implemented the full DeploySure Go backend under `backend/` using only
`net/http` and standard-library packages.

**Packages and files:**

| Package | Files | Responsibility |
|---|---|---|
| `config` | `config.go`, `config_test.go` | Env-var config, PORT validation, no `.env`, no credential logging |
| `reports` | `models.go`, `service.go`, `service_test.go` | File-backed service; `safePath()` prevents `..` traversal |
| `watsonx` | `client.go`, `client_test.go` | 30-second HTTP client, Granite prompt builder, credential safety |
| `api` | `handler.go`, `handler_test.go`, `router.go`, `router_test.go` | All 8 REST handlers, logging middleware, static file serving |
| `cmd/server` | `main.go` | Graceful shutdown via `signal.Notify` |
| `cmd/analyze` | `main.go` | CLI tool: reads file arg or stdin |

**Routes:**
```
GET  /health                    → {"status":"ok"}
GET  /ready                     → {"status":"ready"}
GET  /api/v1/findings/before    → []Finding (JSON)
GET  /api/v1/findings/after     → []Finding (JSON)
GET  /api/v1/reports/before     → {"content":"..."}
GET  /api/v1/reports/after      → {"content":"..."}
GET  /api/v1/comparison         → ComparisonResult
POST /api/v1/granite/analyze    → {"analysis":"..."}
GET  /                          → frontend/index.html
GET  /styles.css                → frontend/styles.css
GET  /app.js                    → frontend/app.js
```

**Test results:**
```
backend/internal/api       37 tests  PASS  coverage: 69.0%
backend/internal/config     6 tests  PASS  coverage: 90.6%
backend/internal/reports    9 tests  PASS  coverage: 80.6%
backend/internal/watsonx    5 tests  PASS  coverage: 95.2%
```

---

### Task 06 — watsonx.ai Granite Integration

**Bob capability:** API integration  
**Evidence:** `bob_sessions/task06_watsonx_integration_summary.md`  
**Screenshots:** `bob_sessions/task06_watsonx_integration_summary.png`

Bob completed `backend/internal/watsonx/client.go` against the current
watsonx.ai REST API:

- IAM token exchange (`POST /identity/token`) with 30-second timeout.
- Granite inference call (`POST /ml/v1/text/generation`) with prompt
  construction, non-200 handling, empty-result detection, and JSON error
  handling.
- Authorization header never logged.
- `backend/internal/config/config.go` extended with `WATSONX_API_KEY`,
  `WATSONX_PROJECT_ID`, `WATSONX_URL`, and `WATSONX_MODEL_ID` vars.
- `backend/cmd/analyze/main.go` updated as a standalone CLI tool.
- All existing tests continue to pass; watsonx coverage at 95.2%.

---

### Task 07 — Frontend Dashboard

**Bob capability:** Frontend development  
**Evidence:** `bob_sessions/task07_frontend_summary.png.md`  
**Screenshots:** `bob_sessions/task07_frontend_summary.png`

Bob built the single-page DeploySure dashboard using only vanilla
HTML/CSS/JavaScript — no external libraries or frameworks.

**Files:**

| File | Content |
|---|---|
| `frontend/index.html` | Semantic HTML, all required element IDs, no external assets |
| `frontend/styles.css` | Responsive layout, severity colour coding, loading states |
| `frontend/app.js` | Fetch calls to all four API endpoints; POST to `/api/v1/granite/analyze` |

Post-generation verification confirmed:
- Every `getElementById` reference in `app.js` present in `index.html`.
- No external libraries referenced.
- All four required API endpoints called.
- POST used for the Granite endpoint.

---

### Task 08 — Remediation Plan (Human Approval Gate)

**Bob capability:** Analysis and planning  
**Evidence:** `bob_sessions/task08_remediation_plan_summary.png.md`  
**Screenshots:** `bob_sessions/task08_remediation_plan_summary.png`

Bob read the three analysis artefacts (`findings-before.json`,
`release-readiness-before.md`, `granite-risk-assessment.json`) and produced
`reports/remediation-plan.md` covering all 41 findings.

**33 remediation entries (REM-001 → REM-036)** — each containing:
- Finding ID and severity
- Affected file and line
- Proposed change (with code diff)
- Reason
- Validation command
- Potential risk

**Breakdown by severity:**

| Priority | Remediations | Key items |
|---|---|---|
| Blocker (5) | REM-001 – REM-005 | `GET /health`, pod `securityContext`, `livenessProbe`, resource limits, `readOnlyRootFilesystem` |
| High (13) | REM-006 – REM-015 | Pin Dockerfile image, `USER` directive, fix silent JSON error, fix `go.mod` version, 3 missing tests |
| Medium (15) | REM-016 – REM-028 | Buffer-first encoding, `WriteHeader`, PORT validation, configurable timeouts, distroless image |
| Low (8) | REM-029 – REM-036 | Dead code removal, `slog` adoption, context propagation, expanded test coverage |

A consolidated execution order table (33 steps) sequences changes so no
dependency is applied before its prerequisite.

**`sample-app/` was not modified. Bob waited for explicit human approval.**

---

### Task 09 — Remediation (Agent Mode)

**Bob capability:** Agent mode  
**Evidence:** `bob_sessions/task09_remediation_summary.png.md`  
**Screenshots:** `bob_sessions/task09_remediation_summary.png`,
`bob_sessions/task09_remediation_summary copy.png`

After explicit human approval, Bob applied all approved remediations to
`sample-app/`:

**Blockers resolved (5/5):**

| REM | Finding | Change |
|---|---|---|
| REM-001 | DC-001 | Added `GET /health` endpoint (`handlers/health.go`) |
| REM-002 | DC-003 / DP-006 | Added `runAsNonRoot: true` pod securityContext to `k8s/deployment.yaml` |
| REM-003 | DC-004 / DP-005 | Added `livenessProbe` targeting `GET /health` |
| REM-004 | DC-005 / DP-003/004 | Added CPU/memory `requests` and `limits` |
| REM-005 | DC-006 | Added `readOnlyRootFilesystem: true` to container securityContext |

**High severity resolved (5 unique + 3 duplicates):**
- REM-006: Pinned Dockerfile builder base image (removed floating `:latest`)
- REM-007: Added non-root `USER` directive to Dockerfile
- REM-008: Fixed silent `json.Encoder` error in `ReadinessHandler`
- REM-009: Fixed invalid Go version in `go.mod`
- REM-010 – REM-012: Added missing unit tests

**Validation results post-remediation:**
```
go build ./sample-app/...     PASS
go test ./sample-app/...      PASS  (all new tests green)
go vet ./sample-app/...       PASS  (0 issues)
```

Readiness probe preserved unchanged. Out-of-scope findings (medium/low)
documented but not applied to keep diff minimal.

---

### Task 10 — Final Validation

**Bob capability:** Parallel subagents  
**Evidence:** `task10_final_validation_summary.md` (repository root)  
**Screenshots:** `bob_sessions/task10_final_validation_summary.png`

Bob re-ran the four specialized subagents against the remediated `sample-app/`
and confirmed all blocker and high-severity findings were resolved.

See `task10_final_validation_summary.md` for the full before/after comparison
and release-readiness decision.

---

## Session Files Index

| File | Task | Type |
|---|---|---|
| `bob_sessions/deploysure_task01_project_scaffolding_summary.png` | 01 | Screenshot |
| `bob_sessions/task01_project_scaffolding_summary.md` | 01 | Session transcript |
| `bob_sessions/task02_document_understanding_prompt.png` | 02 | Screenshot |
| `bob_sessions/task02_document_understanding_summary.png` | 02 | Screenshot |
| `bob_sessions/task02_document_understanding_summary.md` | 02 | Session transcript |
| `bob_sessions/task03_sample_application_build_prompt.png` | 03 | Screenshot |
| `bob_sessions/task03_sample_application_summary.png` | 03 | Screenshot |
| `bob_sessions/task03_sample_application_summary.md` | 03 | Session transcript |
| `bob_sessions/task04_parallel_analysis_summary.png` | 04 | Screenshot |
| `bob_sessions/task04_parallel_analysis_summary.md` | 04 | Session transcript |
| `bob_sessions/task05_backend_summary.png` | 05 | Screenshot |
| `bob_sessions/task05_backend_summary.md` | 05 | Session transcript |
| `bob_sessions/task06_watsonx_integration_summary.png` | 06 | Screenshot |
| `bob_sessions/task06_watsonx_integration_summary.md` | 06 | Session transcript |
| `bob_sessions/task07_frontend_summary.png` | 07 | Screenshot |
| `bob_sessions/task07_frontend_summary.png.md` | 07 | Session transcript |
| `bob_sessions/task08_remediation_plan_summary.png` | 08 | Screenshot |
| `bob_sessions/task08_remediation_plan_summary.png.md` | 08 | Session transcript |
| `bob_sessions/task09_remediation_summary.png` | 09 | Screenshot |
| `bob_sessions/task09_remediation_summary copy.png` | 09 | Screenshot (alternate) |
| `bob_sessions/task09_remediation_summary.png.md` | 09 | Session transcript |
| `bob_sessions/task10_final_validation_summary.png` | 10 | Screenshot |
