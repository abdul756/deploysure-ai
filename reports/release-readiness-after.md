# Release Readiness Report — After Remediation

**Application:** deploysure-ai / sample-app
**Report Date:** 2026-08-30
**Report Stage:** POST-REMEDIATION

---

## Decision

> ## ⚠️ CONDITIONAL PASS
>
> All **Blocker** findings and all **High** findings except one have been resolved.
> The single remaining High finding (TC-004) is **partially resolved** — the
> `cmd/server` package gained its first tests (2.1% coverage) but `main()` startup
> and signal-handling paths remain untested.  No new blockers were introduced.
>
> Release may proceed once the team accepts the residual test-coverage gap and
> acknowledges the 13 remaining Medium findings as tracked tech-debt.

---

## Score

| Metric | Value |
|---|---|
| **Starting score** | 100 |
| Blocker findings (×5) | 0 |
| High findings (×1) | −15 |
| Medium findings (×13) | −91 |
| Low findings (×7) | −14 |
| **Raw total** | −20 |
| **Final score** | **0 / 100** (minimum floor applied) |

> **Note:** The scoring formula applies a minimum floor of 0. Despite clearing all
> blockers, the accumulation of 13 medium findings keeps the raw score below zero.
> This reflects that the medium/low technical-debt items were deliberately out of
> scope for the approved remediation plan and remain unresolved.

---

## Finding Summary

| Severity | Before | After | Resolved | New | Net |
|---|---|---|---|---|---|
| Blocker | 5 | **0** | 5 | 0 | −5 |
| High | 13 | **1** | 12 | 0 | −12 |
| Medium | 15 | **13** | 4 | 2 | −2 |
| Low | 8 | **7** | 1 | 0 | −1 |
| **Total** | **41** | **21** | **22** | **2** | **−20** |

---

## Findings by Source

### Code Quality (12 remaining)

| ID | Severity | Status | Title | File | Line |
|---|---|---|---|---|---|
| CQ-002 | Medium | Persists | http.Error after partial write has no effect | orders.go | 46 |
| CQ-004 | Medium | Persists | Missing explicit WriteHeader(200) in OrdersHandler | orders.go | 45 |
| CQ-008 | Medium | Persists | Lowercase HTTP error message in OrdersHandler | orders.go | 39 |
| CQ-009 | Medium | Persists | Lowercase HTTP error message in ReadinessHandler | readiness.go | 17 |
| CQ-011 | Medium | Persists | Shutdown grace period hardcoded at 30s | main.go | 65 |
| CQ-012 | Medium | Persists | log.Fatalf in goroutine bypasses graceful shutdown | main.go | 54 |
| CQ-015 | Medium | Persists | PORT accepted without numeric validation | main.go | 31 |
| CQ-006 | Low | Persists | DB_DSN read but never used | main.go | 38 |
| CQ-007 | Low | Persists | LOG_LEVEL not respected by logger | main.go | 43 |
| CQ-010 | Low | Persists | HTTP server timeouts hardcoded | main.go | 24 |
| CQ-013 | Low | Persists | Request context not propagated | orders.go | 37 |
| CQ-014 | Low | Persists | Inconsistent log format | main.go | 52 |

**Resolved code findings (3):** CQ-001 (ignored encoder error — fixed), CQ-003 (invalid go version — fixed), CQ-005 (missing WriteHeader in ReadinessHandler — fixed).

### Test Coverage (4 remaining)

| ID | Severity | Status | Title | File | Line |
|---|---|---|---|---|---|
| TC-004 | High | Partially Resolved | cmd/server coverage 2.1% — main() untested | main.go | 30 |
| TC-009 | Medium | New | HealthHandler encode error path uncovered | health.go | 17 |
| TC-010 | Medium | New | ReadinessHandler encode error path uncovered | readiness.go | 23 |
| TC-007 | Low | Persists | Method-not-allowed only covers POST and DELETE | handlers_test.go | 40 |

**Resolved test findings (6 of 8):** TC-001 (TestSeedOrders added), TC-002 (failingWriter test added), TC-003 (TestHealthHandler_GET + TestNewRouter_HealthRoute added), TC-005 (OrdersHandler error path now tested), TC-006 (exact count assertion added), TC-008 (Content-Type assertion added to ReadinessHandler test).

**Test run results (measured):**

```
=== RUN   TestBuildServer              --- PASS (0.00s)
=== RUN   TestBuildServer_DefaultPortPattern --- PASS (0.00s)
=== RUN   TestOrdersHandler_GET        --- PASS (0.00s)
=== RUN   TestOrdersHandler_MethodNotAllowed --- PASS (0.00s)
=== RUN   TestOrdersHandler_EncoderError --- PASS (0.00s)
=== RUN   TestSeedOrders               --- PASS (0.00s)
=== RUN   TestReadinessHandler_GET     --- PASS (0.00s)
=== RUN   TestReadinessHandler_MethodNotAllowed --- PASS (0.00s)
=== RUN   TestHealthHandler_GET        --- PASS (0.00s)
=== RUN   TestHealthHandler_MethodNotAllowed --- PASS (0.00s)
=== RUN   TestNewRouter_OrdersRoute    --- PASS (0.00s)
=== RUN   TestNewRouter_ReadyRoute     --- PASS (0.00s)
=== RUN   TestNewRouter_HealthRoute    --- PASS (0.00s)
=== RUN   TestNewRouter_MethodNotAllowedViaRouter --- PASS (0.00s)

OVERALL RESULT: 14 tests PASS, 0 FAIL
coverage: 2.1% (cmd/server)   ← was 0.0%
coverage: 93.9% (internal/handlers)   ← was 94.4% before; encoder error paths now covered
                                          but two new encode-error paths added and uncovered
```

### Deployment (3 remaining)

| ID | Severity | Status | Title | File | Line |
|---|---|---|---|---|---|
| DP-007 | Medium | Persists | Deployment image uses :latest tag | deployment.yaml | 22 |
| DP-009 | Medium | Persists | Final Docker stage uses debian:bookworm-slim | Dockerfile | 19 |
| DP-008 | Low | Persists | Env vars hardcoded instead of ConfigMap | deployment.yaml | 25 |

**Resolved deployment findings (6 of 8):** DP-001 (USER appuser added), DP-002 (golang:1.22.4-bookworm pinned), DP-003 (resource limits added), DP-004 (resource requests added), DP-005 (livenessProbe added), DP-006 (securityContext.runAsNonRoot + readOnlyRootFilesystem added).

### Documentation Consistency (2 remaining)

| ID | Severity | Status | Title | File | Line |
|---|---|---|---|---|---|
| DC-002 | Medium | Persists | DB_DSN undocumented (FR-34) | requirements.md | 102 |
| DC-009 | Medium | Persists | D-04: Final stage not distroless/scratch | release-runbook.md | 95 |

**Resolved documentation findings (8 of 10):** DC-001 (GET /health implemented), DC-003 (runAsNonRoot added), DC-004 (livenessProbe added), DC-005 (resources added), DC-006 (readOnlyRootFilesystem added), DC-007 (builder pinned), DC-008 (USER directive added), DC-010 (TestSeedOrders added).

---

## Runbook Gate Status

### Docker Gate

| Check | ID | Before | After | Finding |
|---|---|---|---|---|
| Non-root execution | D-01 | ❌ FAIL | ✅ PASS | DP-001 resolved |
| Pinned base image (builder) | D-02 | ❌ FAIL | ✅ PASS | DP-002 resolved: golang:1.22.4-bookworm |
| Multi-stage build | D-03 | ✅ PASS | ✅ PASS | Unchanged |
| Minimal final stage | D-04 | ⚠️ WARN | ⚠️ WARN | DC-009 / DP-009: still debian:bookworm-slim |

### Kubernetes Gate

| Check | ID | Before | After | Finding |
|---|---|---|---|---|
| Resource requests | K-01 | ❌ FAIL | ✅ PASS | DC-005 / DP-004 resolved |
| Resource limits | K-02 | ❌ FAIL | ✅ PASS | DC-005 / DP-003 resolved |
| Liveness probe | K-03 | ❌ FAIL | ✅ PASS | DC-004 / DP-005 resolved: GET /health:8080 |
| Readiness probe | K-04 | ⚠️ NOT CHECKED | ✅ PASS | readinessProbe (GET /ready) present |
| Non-root security context | K-05 | ❌ FAIL | ✅ PASS | DC-003 / DP-006 resolved |
| Read-only root filesystem | K-06 | ❌ FAIL | ✅ PASS | DC-006 resolved |

**All mandatory Kubernetes gates now PASS.**
**Docker gates D-01, D-02, D-03 now PASS. D-04 remains at WARN (medium, non-blocking).**

---

## Remaining Blockers

None. All 5 blocker findings from the baseline have been resolved.

---

## Remaining Tech-Debt (Out-of-Scope Medium/Low)

The following items were not in the approved remediation scope. They are tracked
findings and should be prioritised in a follow-up sprint:

| Priority | ID | Title |
|---|---|---|
| 1 | CQ-012 | Fix log.Fatalf in goroutine (bypasses graceful shutdown) |
| 2 | TC-004 | Increase cmd/server test coverage beyond 2.1% |
| 3 | TC-009/010 | Add encoder error tests for HealthHandler and ReadinessHandler |
| 4 | CQ-002 | Buffer-first JSON encoding in OrdersHandler |
| 5 | CQ-015 | PORT validation at startup |
| 6 | DP-007 | Pin deployment image tag |
| 7 | DP-009 / DC-009 | Switch final Docker stage to distroless |
| 8 | DC-002 | Document DB_DSN in requirements.md |

---

## Files Produced

| File | Description |
|---|---|
| `reports/findings-after.json` | All 21 post-remediation findings merged |
| `reports/code-findings-after.json` | 12 code-quality findings |
| `reports/test-findings-after.json` | 4 test findings |
| `reports/deployment-findings-after.json` | 3 deployment findings |
| `reports/document-findings-after.json` | 2 documentation-consistency findings |
| `reports/release-readiness-after.md` | This report |
