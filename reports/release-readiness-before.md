# Release Readiness Report — Before Fixes

**Application:** deploysure-ai / sample-app
**Report Date:** 2026-08-30
**Report Stage:** PRE-FIX (baseline)

---

## Decision

> ## ❌ BLOCKED
>
> Blocker and High severity findings exist. This application **must not be released** until all Blocker and High findings are resolved.

---

## Score

| Metric | Value |
|---|---|
| **Starting score** | 100 |
| Blocker findings (×5) | −125 |
| High findings (×13) | −195 |
| Medium findings (×15) | −105 |
| Low findings (×8) | −16 |
| **Final score** | **0 / 100** (minimum floor applied) |

---

## Finding Summary

| Severity | Count | Deduction Each | Total Deducted |
|---|---|---|---|
| Blocker | 5 | −25 | −125 |
| High | 13 | −15 | −195 |
| Medium | 15 | −7 | −105 |
| Low | 8 | −2 | −16 |
| **Total** | **41** | | |

---

## Findings by Source

### Code Quality (15 findings)

| ID | Severity | Title | File | Line |
|---|---|---|---|---|
| CQ-001 | High | Ignored error in ReadinessHandler JSON encoding | sample-app/internal/handlers/readiness.go | 22 |
| CQ-003 | High | Invalid Go version in go.mod | sample-app/go.mod | 3 |
| CQ-002 | Medium | http.Error after partial write has no effect | sample-app/internal/handlers/orders.go | 46 |
| CQ-004 | Medium | Missing explicit WriteHeader(200) in OrdersHandler | sample-app/internal/handlers/orders.go | 45 |
| CQ-005 | Medium | Missing explicit WriteHeader(200) in ReadinessHandler | sample-app/internal/handlers/readiness.go | 21 |
| CQ-008 | Medium | Lowercase HTTP error message in OrdersHandler | sample-app/internal/handlers/orders.go | 39 |
| CQ-009 | Medium | Lowercase HTTP error message in ReadinessHandler | sample-app/internal/handlers/readiness.go | 17 |
| CQ-011 | Medium | Shutdown grace period hardcoded at 30s | sample-app/cmd/server/main.go | 62 |
| CQ-012 | Medium | log.Fatalf in goroutine bypasses graceful shutdown | sample-app/cmd/server/main.go | 51 |
| CQ-015 | Medium | PORT accepted without validation | sample-app/cmd/server/main.go | 19 |
| CQ-006 | Low | DB_DSN read but never used | sample-app/cmd/server/main.go | 27 |
| CQ-007 | Low | LOG_LEVEL read but logging doesn't respect it | sample-app/cmd/server/main.go | 32 |
| CQ-010 | Low | HTTP server timeouts are hardcoded | sample-app/cmd/server/main.go | 42 |
| CQ-013 | Low | Request context not propagated | sample-app/internal/handlers/orders.go | 37 |
| CQ-014 | Low | Inconsistent log format | sample-app/cmd/server/main.go | 49 |

### Test Coverage (8 findings)

| ID | Severity | Title | File | Line |
|---|---|---|---|---|
| TC-001 | High | No unit test for exported SeedOrders() | sample-app/internal/handlers/orders.go | 27 |
| TC-002 | High | json.Encoder error branch 0% coverage | sample-app/internal/handlers/orders.go | 47 |
| TC-004 | High | cmd/server/main.go 0.0% coverage | sample-app/cmd/server/main.go | 18 |
| TC-003 | Medium | No test for GET /health route | sample-app/internal/handlers/router.go | 22 |
| TC-005 | Medium | OrdersHandler function coverage 81.8% | sample-app/internal/handlers/orders.go | 37 |
| TC-006 | Medium | TestOrdersHandler_GET weak assertion | sample-app/internal/handlers/handlers_test.go | 34 |
| TC-007 | Low | Method-not-allowed tests only cover POST | sample-app/internal/handlers/handlers_test.go | 83 |
| TC-008 | Low | TestReadinessHandler_GET missing Content-Type assertion | sample-app/internal/handlers/handlers_test.go | 63 |

**Test Run Results:**

```
$ go test ./...
?   github.com/abdul756/deploysure-ai/sample-app/cmd/server   [no test files]
ok  github.com/abdul756/deploysure-ai/sample-app/internal/handlers   (cached)

$ go test -cover ./...
    github.com/abdul756/deploysure-ai/sample-app/cmd/server        coverage: 0.0%
ok  github.com/abdul756/deploysure-ai/sample-app/internal/handlers  coverage: 94.4%
```

- Tests run: 7 | Passed: 7 | Failed: 0 | Skipped: 0
- `cmd/server` package: **0.0%** (no test files)
- `internal/handlers` package: **94.4%** (json.Encode error branch uncovered)

### Deployment (8 findings)

| ID | Severity | Title | File | Line |
|---|---|---|---|---|
| DP-001 | High | Container runs as root — no USER directive | sample-app/deploy/Dockerfile | 43 |
| DP-002 | High | Builder base image uses :latest tag | sample-app/deploy/Dockerfile | 10 |
| DP-003 | High | No CPU/memory resource limits | sample-app/deploy/deployment.yaml | 24 |
| DP-004 | High | No CPU/memory resource requests | sample-app/deploy/deployment.yaml | 24 |
| DP-005 | High | Missing livenessProbe | sample-app/deploy/deployment.yaml | 24 |
| DP-006 | High | Missing securityContext.runAsNonRoot | sample-app/deploy/deployment.yaml | 16 |
| DP-007 | Medium | Deployment image uses :latest tag | sample-app/deploy/deployment.yaml | 25 |
| DP-008 | Low | Env vars hardcoded instead of ConfigMap | sample-app/deploy/deployment.yaml | 28 |

### Documentation Consistency (10 findings)

| ID | Severity | Title | File | Line |
|---|---|---|---|---|
| DC-001 | **Blocker** | GET /health not implemented (FR-32, runbook 3.1/K-03) | docs/requirements.md | 100 |
| DC-003 | **Blocker** | K-05: securityContext.runAsNonRoot absent | docs/release-runbook.md | 111 |
| DC-004 | **Blocker** | K-03: livenessProbe absent | docs/release-runbook.md | 109 |
| DC-005 | **Blocker** | K-01/K-02: resource requests/limits absent | docs/release-runbook.md | 107 |
| DC-006 | **Blocker** | K-06: readOnlyRootFilesystem absent | docs/release-runbook.md | 112 |
| DC-007 | High | D-02: Dockerfile base image not pinned | docs/release-runbook.md | 93 |
| DC-008 | High | D-01: USER directive missing | docs/release-runbook.md | 92 |
| DC-002 | Medium | DB_DSN not documented (FR-34) | docs/requirements.md | 102 |
| DC-009 | Medium | D-04: Final stage not distroless/scratch | docs/release-runbook.md | 95 |
| DC-010 | Medium | FR-36: SeedOrders() has no unit test | docs/requirements.md | 104 |

---

## Runbook Gate Status

### Docker Gate

| Check | ID | Status | Finding |
|---|---|---|---|
| Non-root execution | D-01 | ❌ FAIL | DC-008 / DP-001 |
| Pinned base image | D-02 | ❌ FAIL | DC-007 / DP-002 |
| Multi-stage build | D-03 | ✅ PASS | Present in Dockerfile |
| Minimal final stage | D-04 | ⚠️ WARN | DC-009 (debian:bookworm-slim, not distroless) |

### Kubernetes Gate

| Check | ID | Status | Finding |
|---|---|---|---|
| Resource requests | K-01 | ❌ FAIL | DC-005 / DP-004 |
| Resource limits | K-02 | ❌ FAIL | DC-005 / DP-003 |
| Liveness probe | K-03 | ❌ FAIL | DC-004 / DP-005 |
| Readiness probe | K-04 | ⚠️ NOT CHECKED | No readinessProbe found in deployment.yaml |
| Non-root security context | K-05 | ❌ FAIL | DC-003 / DP-006 |
| Read-only root filesystem | K-06 | ❌ FAIL | DC-006 |

---

## Critical Path to Release

The following items **must** be resolved before re-evaluation:

1. **[DC-001]** Implement `GET /health` endpoint returning `{"status":"ok"}` HTTP 200
2. **[DC-003 / DP-006]** Add `securityContext.runAsNonRoot: true` to pod spec in `deployment.yaml`
3. **[DC-004 / DP-005]** Add `livenessProbe` targeting `GET /health:8080` in `deployment.yaml`
4. **[DC-005 / DP-003 / DP-004]** Add `resources.requests` and `resources.limits` in `deployment.yaml`
5. **[DC-006]** Add `securityContext.readOnlyRootFilesystem: true` to container spec
6. **[DC-007 / DP-002]** Pin Dockerfile builder image from `golang:latest` to a specific version
7. **[DC-008 / DP-001]** Add non-root `USER` directive to Dockerfile final stage
8. **[CQ-001]** Fix ignored JSON encode error in `ReadinessHandler`
9. **[CQ-003]** Fix invalid Go version `1.27.0` in `go.mod`
10. **[TC-001]** Add `TestSeedOrders` unit test
11. **[TC-002]** Test json.Encode error branch (inject failing ResponseWriter)
12. **[TC-004]** Add tests for `cmd/server` package (0% coverage)

---

## Files Produced

| File | Description |
|---|---|
| `reports/findings-before.json` | All 41 findings merged |
| `reports/code-findings-before.json` | 15 code-quality findings |
| `reports/test-findings-before.json` | 8 test findings |
| `reports/deployment-findings-before.json` | 8 deployment findings |
| `reports/document-findings-before.json` | 10 documentation-consistency findings |
| `reports/test-results-before.txt` | Full `go test` and `go test -cover` output |
| `reports/release-readiness-before.md` | This report |
