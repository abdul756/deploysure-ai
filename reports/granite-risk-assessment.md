# Granite Risk Assessment

**Generated:** 2026-08-30T03:58:03Z  
**Model:** ibm/granite-4-h-small  
**Input:** reports/findings-before.json  
**Findings:** 41

---

### Risk Assessment Summary

| **Category** | **Risk ID** | **Title** | **Severity** | **Key Findings** | **Recommendation** |
|--------------|-------------|-----------|--------------|------------------|--------------------|
| **Error Handling** | CQ-001 | Ignored error in ReadinessHandler JSON encoding | High | JSON encoding error ignored → incomplete response. | Check error and return HTTP 500 if encoding fails. |
| **Error Handling** | CQ-002 | http.Error after partial write in OrdersHandler has no effect | Medium | Error handling ineffective after partial write. | Encode to buffer first, then write. |
| **Correctness** | CQ-003 | Invalid Go version in go.mod | High | Specified Go version does not exist. | Change to a real released version (e.g., `go 1.22`). |
| **HTTP Response Handling** | CQ-004 | Missing explicit WriteHeader(200) before body write in OrdersHandler | Medium | Implicit 200 OK is fragile. | Add `w.WriteHeader(http.StatusOK)` after headers. |
| **HTTP Response Handling** | CQ-005 | Missing explicit WriteHeader(200) in ReadinessHandler | Medium | Same as CQ-004. | Add `w.WriteHeader(http.StatusOK)`. |
| **Configuration Validation** | CQ-006 | DB_DSN environment variable read but never used | Low | DB_DSN parsed but not used. | Remove or connect to DB client. |
| **Configuration Validation** | CQ-007 | LOG_LEVEL read but logging does not respect it | Low | LOG_LEVEL read but standard log does not support levels. | Use leveled logger or remove variable. |
| **Error Handling** | CQ-008 | Inconsistent lowercase HTTP error message in OrdersHandler | Medium | Error message uses lowercase text. | Use `http.StatusText()` for consistency. |
| **Error Handling** | CQ-009 | Inconsistent lowercase HTTP error message in ReadinessHandler | Medium | Same as CQ-008. | Use `http.StatusText()`. |
| **Maintainability** | CQ-010 | HTTP server timeouts are hardcoded constants | Low | Timeouts not configurable. | Expose as environment variables. |
| **Graceful Shutdown** | CQ-011 | Shutdown grace period hardcoded at 30s | Medium | Grace period not configurable. | Make configurable via env var. |
| **Error Handling** | CQ-012 | log.Fatalf in goroutine bypasses graceful shutdown | Medium | `log.Fatalf()` exits immediately. | Use error channel to handle shutdown. |
| **Context Usage** | CQ-013 | Request context not propagated to downstream operations | Low | No context propagation. | Pass `r.Context()` to downstream calls. |
| **Maintainability** | CQ-014 | Inconsistent log format across main.go | Low | Mixed log formats. | Adopt single log format. |
| **Correctness** | CQ-015 | PORT environment variable accepted without validation | Medium | Invalid PORT causes runtime failure. | Validate port range and exit with error if invalid. |
| **Missing Test** | TC-001 | No unit test for exported SeedOrders() function | High | Exported function untested. | Add test asserting exact count and values. |
| **Missing Error Path Test** | TC-002 | json.Encoder error branch in OrdersHandler has 0% coverage | High | Error path never executed. | Inject failing writer and assert HTTP 500. |
| **Missing Test** | TC-003 | No test for GET /health route | Medium | Missing route not tested. | Add test for `/health` returning 200. |
| **Coverage Gap** | TC-004 | cmd/server/main.go has 0.0% test coverage | High | Entire entry point untested. | Extract config and server into testable helpers. |
| **Coverage Gap** | TC-005 | OrdersHandler function coverage is 81.8% | High | Error path uncovered. | Inject failing writer and assert HTTP 500. |
| **Missing Test** | TC-006 | TestOrdersHandler_GET does not assert exact order count or field values | Medium | Weak assertion. | Assert len == 3 and compare field values. |
| **Missing Test** | TC-007 | Method-not-allowed tests only cover POST | Low | Non-GET methods not tested. | Use table-driven tests for all methods. |
| **Missing Test** | TC-008 | TestReadinessHandler_GET does not assert Content-Type header | Low | Header not asserted. | Add assertion for `Content-Type: application/json`. |
| **Security** | DP-001 | Container runs as root — no USER directive | High | Runs as root → violates least privilege. | Add USER directive for non-root user. |
| **Security** | DP-002 | Builder base image uses floating :latest tag |
