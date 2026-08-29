# Architecture

## 1. System Overview

DeploySure AI is composed of three logical layers:

1. **Orchestration layer** — IBM Bob Agent mode drives the entire workflow.
2. **Analysis layer** — five specialised subagents, each responsible for one
   review domain, running in parallel.
3. **Intelligence layer** — watsonx.ai Granite 13B receives the aggregated
   findings and returns severity scores and a ranked remediation order.

```
┌────────────────────────────────────────────────────────────┐
│                   Bob Agent (Orchestrator)                 │
│                                                            │
│  1. Reads docs/requirements.md + docs/release-runbook.md   │
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
