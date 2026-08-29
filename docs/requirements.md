# Requirements

## 1. Purpose

This document defines the functional and non-functional requirements for the
DeploySure AI release-readiness platform and the synthetic sample application
used to demonstrate it.

---

## 2. Functional Requirements — Platform

### 2.1 Document Understanding

| ID | Requirement |
|----|------------|
| FR-01 | The system MUST read and parse `docs/requirements.md` and `docs/release-runbook.md` before running any analysis. |
| FR-02 | All checks MUST be grounded in the standards expressed in those documents. |

### 2.2 Parallel Analysis

| ID | Requirement |
|----|------------|
| FR-03 | The following five analyses MUST run concurrently: code quality, unit-test gap, Docker review, Kubernetes review, and docs-to-code consistency. |
| FR-04 | Each analysis MUST be executed by a dedicated specialised subagent with its own isolated context. |

### 2.3 Code-Quality Analysis

| ID | Requirement |
|----|------------|
| FR-05 | The code-quality subagent MUST report `gofmt` formatting violations with file and line references. |
| FR-06 | The code-quality subagent MUST report `go vet` warnings and errors. |
| FR-07 | The code-quality subagent MUST flag exported identifiers lacking GoDoc comments. |

### 2.4 Unit-Test Gap Analysis

| ID | Requirement |
|----|------------|
| FR-08 | The test-gap subagent MUST enumerate all exported functions and HTTP handlers in the sample application. |
| FR-09 | The test-gap subagent MUST cross-reference those identifiers against `*_test.go` files and report uncovered items. |
| FR-10 | Coverage percentage MUST be reported at the package level. |

### 2.5 Docker Review

| ID | Requirement |
|----|------------|
| FR-11 | The Docker subagent MUST flag any `USER root` or absent `USER` directive. |
| FR-12 | The Docker subagent MUST flag base images that use a floating `:latest` tag. |
| FR-13 | The Docker subagent MUST flag any `ENV`, `ARG`, or `RUN` instruction that embeds a secret or credential pattern. |
| FR-14 | The Docker subagent MUST flag unnecessary layers that inflate image size. |

### 2.6 Kubernetes Review

| ID | Requirement |
|----|------------|
| FR-15 | The Kubernetes subagent MUST verify that every container spec defines `resources.requests` and `resources.limits`. |
| FR-16 | The Kubernetes subagent MUST verify the presence of `livenessProbe` and `readinessProbe` on every container. |
| FR-17 | The Kubernetes subagent MUST verify that `securityContext.runAsNonRoot: true` is set. |
| FR-18 | The Kubernetes subagent MUST verify that `securityContext.readOnlyRootFilesystem: true` is set. |

### 2.7 Docs-to-Code Consistency

| ID | Requirement |
|----|------------|
| FR-19 | The consistency subagent MUST compare all API routes documented in `docs/requirements.md` against routes registered in source code. |
| FR-20 | The consistency subagent MUST compare all environment variables documented in `docs/requirements.md` against variables read in source code. |
| FR-21 | Discrepancies MUST be reported with the document line number and the source file reference. |

### 2.8 Risk Prioritisation

| ID | Requirement |
|----|------------|
| FR-22 | All findings from the five subagents MUST be aggregated into a single structured report. |
| FR-23 | The aggregated findings MUST be sent to watsonx.ai Granite for severity scoring. |
| FR-24 | The final report MUST present findings in descending priority order as returned by Granite. |

### 2.9 Human-Approved Remediation

| ID | Requirement |
|----|------------|
| FR-25 | The orchestrating Bob Agent MUST pause after presenting the prioritised report and request explicit human approval before applying any fix. |
| FR-26 | Approved fixes MUST be applied in priority order. |
| FR-27 | A remediation log MUST be produced that maps each finding ID to the applied change. |

### 2.10 Final Validation

| ID | Requirement |
|----|------------|
| FR-28 | After remediation, all five analyses MUST be re-run. |
| FR-29 | The final validation pass MUST confirm zero open findings of severity HIGH or CRITICAL before the release is declared ready. |

---

## 3. Functional Requirements — Sample Application

| ID | Requirement |
|----|------------|
| FR-30 | The application MUST be written in Go using only the standard `net/http` package for HTTP handling. |
| FR-31 | The application MUST expose `GET /api/orders` returning a JSON array of order objects. |
| FR-32 | The application MUST expose `GET /health` returning `{"status":"ok"}` with HTTP 200. |
| FR-33 | The application MUST expose `GET /ready` returning `{"status":"ready"}` with HTTP 200 when the service is ready. |
| FR-34 | All configuration (port, log level, database DSN) MUST be read from environment variables with documented defaults. |
| FR-35 | The application MUST handle `SIGTERM` and `SIGINT` with a graceful HTTP server shutdown (drain in-flight requests). |
| FR-36 | The application MUST include unit tests for every exported handler and helper function. |
| FR-37 | The application MUST ship with a `Dockerfile` that produces a minimal, non-root container image. |
| FR-38 | The application MUST ship with Kubernetes `Deployment` and `Service` manifests. |

---

## 4. Non-Functional Requirements

| ID | Requirement |
|----|------------|
| NFR-01 | The full parallel analysis MUST complete within 120 seconds on a standard developer laptop. |
| NFR-02 | Every Bob session used during development MUST be exported and stored in `bob_sessions/`. |
| NFR-03 | No credentials or secrets MAY appear in any committed file. |
| NFR-04 | All generated reports MUST be stored under `reports/` in Markdown format. |
| NFR-05 | The platform MUST run without requiring a Kubernetes cluster (local Docker is sufficient). |
