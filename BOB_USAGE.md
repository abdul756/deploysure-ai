# BOB_USAGE.md

## How IBM Bob 2.0 Was Used in DeploySure AI

IBM Bob 2.0 was the core development and workflow-orchestration tool used to
build DeploySure AI. We started with the official IBM hackathon repository
template and used Bob Agent mode to create the project structure while
preserving the supplied `.gitignore`, `.bobignore`, `.env.example`, and
`SECURITY.MD` controls. Bob's document understanding was used to read the
requirements and release runbook and translate them into an implementation plan.

Bob then created the synthetic Go sample application, Go report-serving
backend, and lightweight HTML/CSS/JavaScript dashboard. For release analysis,
Bob coordinated focused parallel reviews across Go code quality, unit-test
coverage, Docker and Kubernetes configuration, and documentation-to-code
consistency. The reviews produced structured JSON findings with severity,
affected file, evidence, and recommended action. Bob consolidated these into
the initial release-readiness report.

IBM watsonx.ai with a Granite model was used to prioritise Bob's structured
findings, explain operational impact, and recommend a remediation order. Bob
remained responsible for inspecting and modifying the repository. Before
editing the sample application, Bob created a remediation plan and waited for
human approval. After approval, Bob implemented the selected fixes, ran
`gofmt`, `go vet`, unit tests and coverage checks, and repeated the release
analysis to produce the final before-and-after evidence.

The public repository includes exported Bob task reports and PNG screenshots
of the relevant Bob session summaries in the `bob_sessions/` folder, covering:
scaffolding, documentation, sample application, parallel analysis, backend,
watsonx integration, dashboard, remediation plan, remediation, and final
validation.

No credentials, client data, personal data, or confidential company
information were provided to Bob or committed to the repository.

> **Security note:** All prompts followed the safe-usage guidelines in
> [SECURITY.MD](SECURITY.MD).

---

## Product Boundary

IBM Bob IDE performed repository analysis, implementation, remediation and
validation.

The web dashboard presents Bob-generated reports and calls watsonx.ai Granite.
It does not invoke an undocumented Bob API.

---

## Human Oversight

Bob did not apply remediation until the generated plan was reviewed and
explicitly approved by a human operator. The remediation plan was created in
Task 08 and held until Task 09 was triggered by explicit approval.

---

## Bob Features Demonstrated

| Feature | Tasks |
|---|---|
| Agent mode | 01, 09 |
| Document understanding | 02 |
| Parallel subagents | 04, 10 |
| Code generation | 03, 05, 07 |
| API integration (watsonx.ai) | 06 |
| Analysis and planning | 08 |
| Human-approved remediation | 09 |
| Final validation | 10 |

---

## Session Log

### Task 01 — Project Scaffolding

**Bob capability:** Agent mode
**Evidence:** `bob_sessions/task01_project_scaffolding_summary.md`
**Screenshot:** `bob_sessions/deploysure_task01_project_scaffolding_summary.png`

Bob scaffolded the entire repository from the official IBM Hackathon GitHub
template. Actions taken:

- Created `backend/`, `sample-app/`, `frontend/`, `docs/`, `reports/`,
  `bob_sessions/`, and `evidence/` directories with `.gitkeep` files.
- Generated `README.md`, `PROBLEM_STATEMENT.md`, `SOLUTION_STATEMENT.md`,
  `DEMO_SCRIPT.md`, `BOB_USAGE.md`, and confirmed `SECURITY.MD` was
  unmodified from the template.
- Produced `docs/architecture.md` and `docs/implementation-plan.md`
  capturing the delivery plan and component architecture.

---

### Task 02 — Document Understanding & Architecture

**Bob capability:** Document understanding
**Evidence:** `bob_sessions/task02_document_understanding_summary.md`
**Screenshots:** `bob_sessions/task02_document_understanding_summary.png`,
`bob_sessions/task02_document_understanding_prompt.png`

Bob read the project requirements (`docs/requirements.md`) and release runbook
(`docs/release-runbook.md`) and produced two design documents:

- **`docs/architecture.md`** — full system overview: Bob Agent (orchestrator),
  four specialised subagents (code quality, test coverage, deployment,
  documentation consistency), watsonx.ai Granite integration, data flow, and
  directory layout.
- **`docs/implementation-plan.md`** — phased delivery plan covering project
  scaffolding through final validation, including seeded defects and the
  human-approval gate.

Key design decisions captured:
- IBM Granite selected as the AI model via watsonx.ai.
- Subagent findings aggregated into `reports/findings-before.json`.
- Human-approved gate required before any remediation is applied.

---

### Task 03 — Sample Application

**Bob capability:** Code generation
**Evidence:** `bob_sessions/task03_sample_application_summary.md`
**Screenshots:** `bob_sessions/task03_sample_application_summary.png`,
`bob_sessions/task03_sample_application_build_prompt.png`

Bob generated the synthetic Go orders-API under `sample-app/` with
intentionally seeded defects for later detection by the analysis subagents.

**Files created:**

- `sample-app/cmd/server/main.go`
- `sample-app/internal/handlers/readiness.go`
- `sample-app/internal/handlers/orders.go`
- `sample-app/internal/handlers/router.go`
- `sample-app/deploy/Dockerfile`
- `sample-app/deploy/deployment.yaml`
- `sample-app/deploy/service.yaml`
- `sample-app/go.mod`
- `docs/requirements.md` and `docs/release-runbook.md` (spec documents)

**Seeded defects (for detection in Task 04):**

| ID | Defect | Location |
|---|---|---|
| SD-01 / DP-002 | Floating `:latest` tag on builder stage | `Dockerfile` |
| SD-02 / DP-001 | Missing `USER` directive — container runs as root | `Dockerfile` |
| SD-03 / DC-001 | Missing `GET /health` endpoint | `handlers/` |
| SD-04 / DP-005 | No `livenessProbe` | `deployment.yaml` |
| SD-05 / DP-006 | No pod-level `securityContext` (`runAsNonRoot`) | `deployment.yaml` |
| SD-06 / DP-003/004 | No resource `requests`/`limits` | `deployment.yaml` |
| SD-07 / DC-006 | No `readOnlyRootFilesystem` | `deployment.yaml` |
| SD-08 / CQ-001 | `json.Encoder` error silently ignored | `handlers/readiness.go` |
| SD-09 / CQ-003 | Invalid Go version in `go.mod` (`go 1.27.0`) | `go.mod` |

---

### Task 04 — Parallel Analysis (Before)

**Bob capability:** Parallel subagents
**Evidence:** `bob_sessions/task04_parallel_analysis_summary.md`
**Screenshot:** `bob_sessions/task04_parallel_analysis_summary.png`

Bob spawned four independent specialised subagents simultaneously to analyse
`sample-app/`:

| Subagent | Findings | Key issues |
|---|---|---|
| Code Quality | 15 | Silent encoder error, dead config code, invalid `go.mod` version, missing `WriteHeader` |
| Test Coverage | 8 | 0% `cmd/server` coverage, uncovered error branches, weak assertions |
| Deployment (Docker + K8s) | 8 | Floating `:latest`, no `USER`, missing probes, no resource limits, no securityContext |
| Documentation Consistency | 10 | 5 runbook gate blockers, undocumented env vars, missing required tests |

**Decision: ❌ BLOCKED — 5 blocker findings, not release-ready.**

Files produced:
- `reports/findings-before.json` (41 total findings)
- `reports/code-findings-before.json`
- `reports/test-findings-before.json`
- `reports/deployment-findings-before.json`
- `reports/document-findings-before.json`
- `reports/release-readiness-before.md`
- `reports/test-results-before.txt`

---

### Task 05 — Go Backend

**Bob capability:** Code generation
**Evidence:** `bob_sessions/task05_backend_summary.md`
**Screenshot:** `bob_sessions/task05_backend_summary.png`

Bob implemented the full DeploySure Go backend under `backend/` using only
the standard library.

**Packages and files:**

| Package | Files | Responsibility |
|---|---|---|
| `config` | `config.go`, `config_test.go` | Env-var config, PORT validation, no credential logging |
| `reports` | `models.go`, `service.go`, `service_test.go` | File-backed report service; `safePath()` prevents path traversal |
| `watsonx` | `client.go`, `client_test.go` | IAM token exchange, Granite inference, credential safety |
| `api` | `handler.go`, `handler_test.go`, `router.go`, `router_test.go` | REST handlers, logging middleware, static file serving |
| `cmd/server` | `main.go` | Graceful shutdown via `signal.Notify` and error channel |
| `cmd/analyze` | `main.go` | One-shot CLI: reads findings file and calls Granite |

**Test results:**
```
backend/internal/api       PASS  coverage: 69.0%
backend/internal/config    PASS  coverage: 90.6%
backend/internal/reports   PASS  coverage: 80.6%
backend/internal/watsonx   PASS  coverage: 95.2%
```

---

### Task 06 — watsonx.ai Granite Integration

**Bob capability:** API integration
**Evidence:** `bob_sessions/task06_watsonx_integration_summary.md`
**Screenshot:** `bob_sessions/task06_watsonx_integration_summary.png`

Bob completed `backend/internal/watsonx/client.go` against the watsonx.ai
REST API:

- IAM token exchange (`POST /identity/token`) with 30-second timeout.
- Granite inference call (`POST /ml/v1/text/generation`) with structured
  prompt construction, non-200 error handling, and empty-result detection.
- Authorization header never logged.
- Config extended with `IBM_CLOUD_API_KEY`, `WATSONX_PROJECT_ID`,
  `WATSONX_URL`, and `WATSONX_MODEL_ID` variables.
- `backend/cmd/analyze/main.go` implemented as a standalone CLI tool that
  writes `reports/granite-risk-assessment.md` and
  `reports/granite-risk-assessment.json`.
- All tests continue to pass; watsonx package coverage at 95.2%.

---

### Task 07 — Frontend Dashboard

**Bob capability:** Code generation
**Evidence:** `bob_sessions/task07_frontend_summary.png.md`
**Screenshot:** `bob_sessions/task07_frontend_summary.png`

Bob built the single-page DeploySure dashboard using vanilla
HTML/CSS/JavaScript — no external libraries or frameworks.

| File | Content |
|---|---|
| `frontend/index.html` | Semantic HTML, all required element IDs |
| `frontend/styles.css` | Responsive layout, severity colour coding, loading states |
| `frontend/app.js` | Fetch calls to all API endpoints; POST to `/api/v1/granite/analyze` |

Post-generation verification confirmed every `getElementById` reference in
`app.js` is present in `index.html`, and all required API endpoints are called.

---

### Task 08 — Remediation Plan (Human Approval Gate)

**Bob capability:** Analysis and planning
**Evidence:** `bob_sessions/task08_remediation_plan_summary.png.md`
**Screenshot:** `bob_sessions/task08_remediation_plan_summary.png`

Bob read `reports/findings-before.json`, `reports/release-readiness-before.md`,
and `reports/granite-risk-assessment.json` and produced
`reports/remediation-plan.md` covering all 41 findings across 36 remediation
entries (REM-001 → REM-036), each containing:

- Finding ID and severity
- Affected file and line
- Proposed change with code or YAML diff
- Reason and validation command
- Potential risk

| Priority | Remediations | Key items |
|---|---|---|
| Blocker (5) | REM-001 – REM-005 | `GET /health`, pod `securityContext`, `livenessProbe`, resource limits, `readOnlyRootFilesystem` |
| High (13) | REM-006 – REM-015 | Pin Dockerfile image, `USER` directive, fix silent JSON error, fix `go.mod` version, 3 missing tests |
| Medium (15) | REM-016 – REM-028 | Buffer-first encoding, `WriteHeader`, PORT validation, configurable timeouts |
| Low (8) | REM-029 – REM-036 | Dead code, log format consistency, context propagation |

**`sample-app/` was not modified. Bob waited for explicit human approval.**

---

### Task 09 — Remediation

**Bob capability:** Agent mode
**Evidence:** `bob_sessions/task09_remediation_summary.png.md`
**Screenshots:** `bob_sessions/task09_remediation_summary.png`,
`bob_sessions/task09_remediation_summary copy.png`

After explicit human approval, Bob applied all approved remediations to
`sample-app/`:

**Blockers resolved (5 / 5):**

| REM | Finding | Change |
|---|---|---|
| REM-001 | DC-001 | Added `GET /health` endpoint (`health.go` + registered in `router.go`) |
| REM-002 | DC-003 / DP-006 | Added `runAsNonRoot: true` + `runAsUser: 1000` to pod `securityContext` |
| REM-003 | DC-004 / DP-005 | Added `livenessProbe` targeting `GET /health:8080` |
| REM-004 | DC-005 / DP-003/004 | Added CPU/memory `requests` and `limits` |
| REM-005 | DC-006 | Added `readOnlyRootFilesystem: true` + `allowPrivilegeEscalation: false` |

**High severity resolved (8):**
- REM-006 / REM-007: Pinned Dockerfile builder image; added non-root `USER` directive
- REM-008: Fixed silent `json.Encoder` error in `ReadinessHandler`
- REM-009: Fixed invalid Go version in `go.mod` (`1.27.0` → `1.22.0`)
- REM-010 – REM-012: Added `TestSeedOrders`, `TestOrdersHandler_EncoderError`, and `TestBuildServer`

**Validation results:**
```
gofmt     PASS — no files need reformatting
go vet    PASS — 0 issues
go test   PASS — 14 tests, 0 failures
coverage: cmd/server 2.1% (was 0.0%), internal/handlers 93.9%
```

Out-of-scope medium/low findings documented in `reports/remediation-summary.md`
but not applied, keeping the diff minimal.

---

### Task 10 — Final Validation

**Bob capability:** Parallel subagents
**Evidence:** `task10_final_validation_summary.md` (repository root)
**Screenshot:** `bob_sessions/task10_final_validation_summary.png`

Bob re-ran the four specialised subagents against the remediated `sample-app/`
and confirmed all blocker and high-severity findings were resolved.

| Metric | Before | After |
|---|---|---|
| Blockers | 5 | **0** |
| High | 13 | **1** (partially resolved) |
| Medium | 15 | 13 |
| Low | 8 | 7 |
| Total findings | 41 | **21** |
| Release decision | ❌ BLOCKED | ⚠️ CONDITIONAL PASS |
| Tests | 7 pass | **14 pass** |

Files produced:
- `reports/findings-after.json`
- `reports/code-findings-after.json`
- `reports/test-findings-after.json`
- `reports/deployment-findings-after.json`
- `reports/document-findings-after.json`
- `reports/release-readiness-after.md`
- `reports/hackathon-impact.md`

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
| `task10_final_validation_summary.md` | 10 | Session transcript (repo root) |
