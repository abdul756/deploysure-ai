# Remediation Plan — deploysure-ai / sample-app

**Based on:** reports/findings-before.json · reports/release-readiness-before.md · reports/granite-risk-assessment.md  
**Report Date:** 2026-08-30  
**Status:** AWAITING APPROVAL — no sample-app files have been modified  

---

## Overview

| Severity | Findings | Must fix before release |
|---|---|---|
| Blocker | 5 | ✅ Yes |
| High | 13 | ✅ Yes |
| Medium | 15 | Recommended |
| Low | 8 | Optional |
| **Total** | **41** | |

Remediations are ordered: Blockers → High → Medium → Low. Within each severity, deployment fixes come before code fixes to unblock the Kubernetes gates first.

---

## Part 1 — Blockers (5)

---

### REM-001 · DC-001 — Implement GET /health endpoint

| Field | Detail |
|---|---|
| **Finding ID** | DC-001 |
| **Severity** | Blocker |
| **Affected file** | `sample-app/internal/handlers/router.go` + new `sample-app/internal/handlers/health.go` |

**Proposed change**

1. Create `sample-app/internal/handlers/health.go`:
   ```go
   package handlers

   import (
       "encoding/json"
       "net/http"
   )

   func HealthHandler(w http.ResponseWriter, r *http.Request) {
       if r.Method != http.MethodGet {
           http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
           return
       }
       w.Header().Set("Content-Type", "application/json")
       w.WriteHeader(http.StatusOK)
       if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
           http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
       }
   }
   ```
2. In `sample-app/internal/handlers/router.go`, add:
   ```go
   mux.HandleFunc("/health", HealthHandler)
   ```

**Reason**  
FR-32 mandates `GET /health` returning `{"status":"ok"}` HTTP 200. The route was deliberately omitted (SD-01). Its absence also blocks Kubernetes Gate K-03 (livenessProbe) and K-04 (readinessProbe). Five findings (DC-001, DC-004, TC-003, DP-005, plus the livenessProbe YAML change in REM-003) all depend on this endpoint existing.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestHealthHandler -v
# Also: curl -s http://localhost:8080/health | jq .
```

**Potential risk**  
Low. Adding a new route does not alter existing routes. The only risk is a name collision if `HealthHandler` is already declared elsewhere — verify with `grep -r "HealthHandler" sample-app/`.

---

### REM-002 · DC-003 / DP-006 — Add pod-level securityContext (runAsNonRoot)

| Field | Detail |
|---|---|
| **Finding ID** | DC-003, DP-006 |
| **Severity** | Blocker |
| **Affected file** | `sample-app/deploy/deployment.yaml` |

**Proposed change**

Add a `securityContext` block at the **pod spec** level (under `spec.template.spec`):
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
```

**Reason**  
Kubernetes Gate K-05 will hard-block the release without `runAsNonRoot: true`. Combined with the Dockerfile `USER` directive (REM-007), this ensures the process never runs as UID 0 and satisfies the principle of least privilege.

**Validation command**
```bash
kubectl apply --dry-run=client -f sample-app/deploy/deployment.yaml
# Confirm with: kubectl get pod <pod> -o jsonpath='{.spec.securityContext}'
```

**Potential risk**  
Medium. If the container binary was installed into a root-owned path (e.g., `/app` with mode 700), the non-root user will be unable to execute it. Verify after applying REM-007 (Dockerfile `USER`) that the binary is world-executable.

---

### REM-003 · DC-004 / DP-005 — Add livenessProbe targeting GET /health

| Field | Detail |
|---|---|
| **Finding ID** | DC-004, DP-005 |
| **Severity** | Blocker |
| **Affected file** | `sample-app/deploy/deployment.yaml` |

**Proposed change**

Add under the container spec:
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 3
```

**Reason**  
Kubernetes Gate K-03 requires a livenessProbe. Without it Kubernetes cannot detect a hung application and will never restart it. This fix depends on REM-001 (`/health` must exist first).

**Validation command**
```bash
kubectl apply --dry-run=client -f sample-app/deploy/deployment.yaml
# In-cluster: kubectl describe pod <pod> | grep -A5 Liveness
```

**Potential risk**  
Low-Medium. An `initialDelaySeconds` that is too short will cause restart loops during slow startup. The proposed 5 s is safe for this application but should be tuned if startup latency increases.

---

### REM-004 · DC-005 / DP-003 / DP-004 — Add resource requests and limits

| Field | Detail |
|---|---|
| **Finding ID** | DC-005, DP-003, DP-004 |
| **Severity** | Blocker |
| **Affected file** | `sample-app/deploy/deployment.yaml` |

**Proposed change**

Add under the container spec:
```yaml
resources:
  requests:
    cpu: "100m"
    memory: "64Mi"
  limits:
    cpu: "500m"
    memory: "128Mi"
```

**Reason**  
Kubernetes Gates K-01 (requests) and K-02 (limits) both fail without this block. Without requests the scheduler makes uninformed placement decisions; without limits a single misbehaving pod can exhaust node resources and destabilise other workloads.

**Validation command**
```bash
kubectl apply --dry-run=client -f sample-app/deploy/deployment.yaml
# In-cluster: kubectl describe pod <pod> | grep -A6 Limits
```

**Potential risk**  
Medium. The memory limit of 128 Mi is conservative. If the application's working set grows (e.g., in-memory order cache), it will be OOMKilled. Monitor actual usage for the first release and adjust limits upward if needed.

---

### REM-005 · DC-006 — Add readOnlyRootFilesystem to container securityContext

| Field | Detail |
|---|---|
| **Finding ID** | DC-006 |
| **Severity** | Blocker |
| **Affected file** | `sample-app/deploy/deployment.yaml` |

**Proposed change**

Add a `securityContext` block at the **container** level:
```yaml
securityContext:
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
```

**Reason**  
Kubernetes Gate K-06 requires `readOnlyRootFilesystem: true` on the container security context. A read-only root filesystem prevents an attacker who achieves code execution from persisting malicious binaries or modifying application files.

**Validation command**
```bash
kubectl apply --dry-run=client -f sample-app/deploy/deployment.yaml
# In-cluster: kubectl get pod <pod> -o jsonpath='{.spec.containers[0].securityContext}'
```

**Potential risk**  
Medium. If the application writes to the filesystem at runtime (temp files, logs, PID files), it will crash once this flag is set. Audit all `os.Create`, `os.OpenFile`, and `log` file writes before enabling. Standard output logging (the current pattern) is unaffected.

---

## Part 2 — High Severity (13)

---

### REM-006 · DC-007 / DP-002 — Pin Dockerfile builder base image

| Field | Detail |
|---|---|
| **Finding ID** | DC-007, DP-002 |
| **Severity** | High |
| **Affected file** | `sample-app/deploy/Dockerfile` line 10 |

**Proposed change**
```dockerfile
# Before
FROM golang:latest AS builder

# After
FROM golang:1.22.4-bookworm AS builder
```

**Reason**  
Docker Gate D-02 requires a pinned base image. `golang:latest` is resolved at build time and can silently introduce a newer, unvetted Go toolchain or OS packages, breaking reproducibility and defeating rollback guarantees.

**Validation command**
```bash
docker build -f sample-app/deploy/Dockerfile sample-app/ --no-cache --progress=plain 2>&1 | head -20
```

**Potential risk**  
Low. The only risk is that `golang:1.22.4-bookworm` itself has known CVEs at the time of fix. Run `docker scout cves golang:1.22.4-bookworm` and pin to the latest patch that has no critical CVEs if needed.

---

### REM-007 · DC-008 / DP-001 — Add non-root USER directive to Dockerfile

| Field | Detail |
|---|---|
| **Finding ID** | DC-008, DP-001 |
| **Severity** | High |
| **Affected file** | `sample-app/deploy/Dockerfile` final stage (around line 23) |

**Proposed change**

In the final stage, before the `ENTRYPOINT`:
```dockerfile
RUN addgroup --system appgroup \
    && adduser --system --ingroup appgroup --no-create-home appuser
USER appuser
```

**Reason**  
Docker Gate D-01 requires a non-root USER. Running as UID 0 means a container escape gives an attacker root on the host. The Kubernetes `runAsNonRoot` fix (REM-002) also depends on the image having a valid non-root user.

**Validation command**
```bash
docker build -t orders-api:test -f sample-app/deploy/Dockerfile sample-app/
docker run --rm orders-api:test whoami   # must print "appuser", not "root"
```

**Potential risk**  
Medium. If any path inside the container (`/app`, mounted volumes, config files) is only readable by root, the process will fail to start. Ensure `COPY --chown=appuser:appgroup` is used when copying the binary.

---

### REM-008 · CQ-001 — Handle json.Encoder error in ReadinessHandler

| Field | Detail |
|---|---|
| **Finding ID** | CQ-001 |
| **Severity** | High |
| **Affected file** | `sample-app/internal/handlers/readiness.go` line 22 |

**Proposed change**
```go
// Before
json.NewEncoder(w).Encode(statusResponse{Status: "ready"}) //nolint:errcheck

// After
if err := json.NewEncoder(w).Encode(statusResponse{Status: "ready"}); err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
}
```
Also remove the `//nolint:errcheck` directive.

**Reason**  
The encoding error is silently ignored. If `Encode` fails, the response body is truncated and the client receives malformed JSON with no indication of failure. Suppressing the linter warning with `//nolint` hides the problem from automated review.

**Validation command**
```bash
cd sample-app && go vet ./internal/handlers/... && go test ./internal/handlers/... -v
```

**Potential risk**  
Very low. The `http.Error` call after a partial write may not change the already-sent status code (see CQ-002 / REM-012), but it is still an improvement over silent failure. The real fix for that race is REM-012.

---

### REM-009 · CQ-003 — Fix invalid Go version in go.mod

| Field | Detail |
|---|---|
| **Finding ID** | CQ-003 |
| **Severity** | High |
| **Affected file** | `sample-app/go.mod` line 3 |

**Proposed change**
```
# Before
go 1.27.0

# After
go 1.22.0
```

**Reason**  
Go 1.27.0 does not exist. Toolchains that enforce the `go` directive (Go 1.21+) will refuse to build with a version they cannot resolve, producing `go: updates to go.mod needed; to update it, run: go mod tidy`. This will fail CI builds.

**Validation command**
```bash
cd sample-app && go mod tidy && go build ./...
```

**Potential risk**  
Low. Choosing `go 1.22.0` enables new language features (range-over-integer, etc.). If existing code accidentally relies on pre-1.22 behaviour this could surface a compile error — but that is a desirable catch, not a risk.

---

### REM-010 · TC-001 / DC-010 — Add TestSeedOrders unit test

| Field | Detail |
|---|---|
| **Finding ID** | TC-001, DC-010 |
| **Severity** | High |
| **Affected file** | `sample-app/internal/handlers/handlers_test.go` |

**Proposed change**

Add a new test function:
```go
func TestSeedOrders(t *testing.T) {
    orders := handlers.SeedOrders()
    if got, want := len(orders), 3; got != want {
        t.Fatalf("SeedOrders() len = %d, want %d", got, want)
    }
    if orders[0].ID != "ord-001" {
        t.Errorf("orders[0].ID = %q, want %q", orders[0].ID, "ord-001")
    }
    if orders[1].ID != "ord-002" {
        t.Errorf("orders[1].ID = %q, want %q", orders[1].ID, "ord-002")
    }
    if orders[2].ID != "ord-003" {
        t.Errorf("orders[2].ID = %q, want %q", orders[2].ID, "ord-003")
    }
}
```

**Reason**  
FR-36 requires unit tests for all exported functions. `SeedOrders()` is the fixture source for the entire orders test suite; if it mutates silently the `len > 0` assertion in `TestOrdersHandler_GET` will not catch it.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestSeedOrders -v
```

**Potential risk**  
None. Additive test only.

---

### REM-011 · TC-002 / TC-005 — Test json.Encoder error branch with failing ResponseWriter

| Field | Detail |
|---|---|
| **Finding ID** | TC-002, TC-005 |
| **Severity** | High |
| **Affected file** | `sample-app/internal/handlers/handlers_test.go` |

**Proposed change**

Add a `failingWriter` helper and a new test:
```go
type failingWriter struct{ http.ResponseWriter }

func (f *failingWriter) Write([]byte) (int, error) {
    return 0, fmt.Errorf("disk full")
}

func TestOrdersHandler_EncoderError(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/orders", nil)
    rec := httptest.NewRecorder()
    fw := &failingWriter{ResponseWriter: rec}
    handlers.OrdersHandler(fw, req)
    if rec.Code != http.StatusInternalServerError {
        t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
    }
}
```

**Reason**  
The error path on lines 47–49 of `orders.go` has 0% coverage. Without exercising it, a regression (e.g., accidentally removing the `if err` block) would go undetected. This test drives coverage to ~100% for `OrdersHandler`.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -cover -v
# Expect OrdersHandler coverage to reach 100%
```

**Potential risk**  
Low. The `failingWriter` approach requires that REM-012 (buffer-first encoding) is applied first; otherwise `http.Error` after a partial write still does not set the correct status code and the assertion on `rec.Code` may be unreliable.

---

### REM-012 · TC-004 — Add tests for cmd/server package (0% coverage)

| Field | Detail |
|---|---|
| **Finding ID** | TC-004 |
| **Severity** | High |
| **Affected file** | new `sample-app/cmd/server/main_test.go` |

**Proposed change**

Extract server construction into a testable helper in `main.go`:
```go
func buildServer(port string) *http.Server {
    mux := handlers.NewRouter()
    return &http.Server{
        Addr:         ":" + port,
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
}
```

Then add `main_test.go`:
```go
package main

import (
    "testing"
)

func TestBuildServer(t *testing.T) {
    srv := buildServer("9999")
    if srv.Addr != ":9999" {
        t.Errorf("Addr = %q, want %q", srv.Addr, ":9999")
    }
    if srv.ReadTimeout == 0 {
        t.Error("ReadTimeout must be non-zero")
    }
}
```

**Reason**  
`cmd/server/main.go` has 0.0% coverage. The environment-variable defaulting, server construction, and graceful-shutdown wiring are completely untested. A regression in any of these would only surface in production.

**Validation command**
```bash
cd sample-app && go test ./cmd/server/... -cover -v
# Expect coverage > 0%
```

**Potential risk**  
Low. Refactoring `main()` to call `buildServer()` is a safe, internal restructuring. Care must be taken not to change signal-handling or `log.Fatal` behaviour while extracting.

---

### REM-013 · DP-003 / DP-004 (duplicate gate) — Already covered by REM-004

> See **REM-004**. DP-003 and DP-004 are resolved by adding `resources.requests` and `resources.limits` in `deployment.yaml`.

---

### REM-014 · DP-005 (duplicate gate) — Already covered by REM-003

> See **REM-003**. DP-005 is resolved by adding `livenessProbe` in `deployment.yaml`.

---

### REM-015 · DP-006 (duplicate gate) — Already covered by REM-002

> See **REM-002**. DP-006 is resolved by adding pod-level `securityContext.runAsNonRoot`.

---

## Part 3 — Medium Severity (15)

---

### REM-016 · CQ-002 — Encode to buffer before writing in OrdersHandler

| Field | Detail |
|---|---|
| **Finding ID** | CQ-002 |
| **Severity** | Medium |
| **Affected file** | `sample-app/internal/handlers/orders.go` line 46 |

**Proposed change**
```go
// Before
w.Header().Set("Content-Type", "application/json")
if err := json.NewEncoder(w).Encode(orders); err != nil {
    http.Error(w, "internal server error", http.StatusInternalServerError)
    return
}

// After
b, err := json.Marshal(orders)
if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write(b)
```

**Reason**  
Once `json.NewEncoder(w).Encode()` starts writing to the `ResponseWriter`, the HTTP status line has already been sent. A subsequent `http.Error()` call cannot change the status code; the client receives a 200 with a truncated body followed by error text — the worst of both worlds. Encoding to a buffer first ensures the status code can be set correctly on error.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -v
```

**Potential risk**  
Low. `json.Marshal` loads the entire response into memory before writing, which is safe for the fixed-size seed data. For very large payloads in future, streaming would be preferred — document this trade-off in a code comment.

---

### REM-017 · CQ-004 — Add explicit WriteHeader(200) in OrdersHandler

| Field | Detail |
|---|---|
| **Finding ID** | CQ-004 |
| **Severity** | Medium |
| **Affected file** | `sample-app/internal/handlers/orders.go` line 45 |

**Proposed change**

This is resolved as part of REM-016 — the buffer-first approach adds `w.WriteHeader(http.StatusOK)` explicitly.

**Reason**  
Implicit 200 is fragile; middleware that inspects the written status code before it is set will see 0 rather than 200.

**Validation command**  
Same as REM-016.

**Potential risk**  
None beyond REM-016.

---

### REM-018 · CQ-005 — Add explicit WriteHeader(200) in ReadinessHandler

| Field | Detail |
|---|---|
| **Finding ID** | CQ-005 |
| **Severity** | Medium |
| **Affected file** | `sample-app/internal/handlers/readiness.go` line 21 |

**Proposed change**

Add `w.WriteHeader(http.StatusOK)` after `w.Header().Set(...)` and before `json.NewEncoder`:
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
if err := json.NewEncoder(w).Encode(statusResponse{Status: "ready"}); err != nil {
    // WriteHeader already sent; log the error only
    log.Printf("readiness encode error: %v", err)
}
```

**Reason**  
Consistent with CQ-004. Also resolves the ambiguity between the identical `/health` and `/ready` response contracts.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestReadinessHandler -v
```

**Potential risk**  
Very low.

---

### REM-019 · CQ-008 — Use http.StatusText in OrdersHandler error response

| Field | Detail |
|---|---|
| **Finding ID** | CQ-008 |
| **Severity** | Medium |
| **Affected file** | `sample-app/internal/handlers/orders.go` line 39 |

**Proposed change**
```go
// Before
http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

// After
http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
```

**Reason**  
`http.StatusText()` returns the canonical RFC 9110 phrase (`Method Not Allowed`), which is more interoperable with HTTP clients that parse the reason phrase. It also eliminates an ad-hoc string that could drift from the status code.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestOrdersHandler_MethodNotAllowed -v
```

**Potential risk**  
None. The status code is unchanged; only the message body differs.

---

### REM-020 · CQ-009 — Use http.StatusText in ReadinessHandler error response

| Field | Detail |
|---|---|
| **Finding ID** | CQ-009 |
| **Severity** | Medium |
| **Affected file** | `sample-app/internal/handlers/readiness.go` line 17 |

**Proposed change**
```go
// Before
http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

// After
http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
```

**Reason**  
Identical to CQ-008; consistent error messages across handlers reduce surprises for consumers.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestReadinessHandler_MethodNotAllowed -v
```

**Potential risk**  
None.

---

### REM-021 · CQ-011 — Make shutdown grace period configurable

| Field | Detail |
|---|---|
| **Finding ID** | CQ-011 |
| **Severity** | Medium |
| **Affected file** | `sample-app/cmd/server/main.go` line 62 |

**Proposed change**
```go
shutdownSec := os.Getenv("SHUTDOWN_TIMEOUT_SEC")
if shutdownSec == "" {
    shutdownSec = "30"
}
shutdownTimeout, err := strconv.Atoi(shutdownSec)
if err != nil || shutdownTimeout <= 0 {
    log.Fatalf("invalid SHUTDOWN_TIMEOUT_SEC: %q", shutdownSec)
}
ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
```

**Reason**  
Environments with strict rolling-update windows may need a shorter grace period; high-traffic scenarios may need a longer one. Hard-coding 30 s prevents tuning without a rebuild.

**Validation command**
```bash
cd sample-app && SHUTDOWN_TIMEOUT_SEC=5 go run ./cmd/server/... &
sleep 1 && kill -SIGTERM $! && wait
```

**Potential risk**  
Low. A misconfigured value (0 or negative) must be rejected at startup. The validation above achieves this. An operator setting a very short timeout risks in-flight requests being killed; document acceptable range in `docs/requirements.md`.

---

### REM-022 · CQ-012 — Replace log.Fatalf in goroutine with error channel

| Field | Detail |
|---|---|
| **Finding ID** | CQ-012 |
| **Severity** | Medium |
| **Affected file** | `sample-app/cmd/server/main.go` line 51 |

**Proposed change**
```go
serverErr := make(chan error, 1)
go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        serverErr <- err
    }
}()

select {
case sig := <-quit:
    log.Printf("received signal %v, shutting down", sig)
case err := <-serverErr:
    log.Printf("server error: %v", err)
}
// proceed with graceful shutdown ...
```

**Reason**  
`log.Fatalf` calls `os.Exit(1)` immediately, bypassing all `defer` statements including `cancel()` and the graceful shutdown block. The error channel pattern allows `main()` to handle both signal-initiated and error-initiated shutdown through a single path.

**Validation command**
```bash
cd sample-app && go vet ./cmd/server/... && go build ./cmd/server/...
```

**Potential risk**  
Medium. This is a non-trivial refactor of the startup/shutdown flow. Must be paired with a test (REM-012 / TC-004) to verify the goroutine and select work correctly.

---

### REM-023 · CQ-015 — Validate PORT environment variable

| Field | Detail |
|---|---|
| **Finding ID** | CQ-015 |
| **Severity** | Medium |
| **Affected file** | `sample-app/cmd/server/main.go` line 19 |

**Proposed change**
```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
portNum, err := strconv.Atoi(port)
if err != nil || portNum < 1 || portNum > 65535 {
    log.Fatalf("invalid PORT %q: must be an integer between 1 and 65535", port)
}
```

**Reason**  
An invalid `PORT` value such as `"abc"` or `"99999"` causes `ListenAndServe` to fail with a cryptic `address abc: missing port in address` error at runtime rather than at startup. Failing fast with a descriptive message is much easier to diagnose in a containerised environment.

**Validation command**
```bash
cd sample-app && PORT=abc go run ./cmd/server/... 2>&1 | grep -i "invalid"
# Should print a clear error and exit non-zero
```

**Potential risk**  
Very low. Existing valid deployments using `PORT=8080` are unaffected.

---

### REM-024 · TC-003 — Add test for GET /health route

| Field | Detail |
|---|---|
| **Finding ID** | TC-003 |
| **Severity** | Medium |
| **Affected file** | `sample-app/internal/handlers/handlers_test.go` |

**Proposed change**
```go
func TestHealthHandler_GET(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()
    handlers.HealthHandler(rec, req)
    resp := rec.Result()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
    }
    if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
        t.Errorf("Content-Type = %q, want application/json", ct)
    }
    var body map[string]string
    if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
        t.Fatalf("decode error: %v", err)
    }
    if body["status"] != "ok" {
        t.Errorf("status = %q, want \"ok\"", body["status"])
    }
}
```

**Reason**  
The `/health` route is invisible to the test suite without this test. TC-003 requires explicit route-level coverage of the new `HealthHandler`.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestHealthHandler -v
```

**Potential risk**  
None. Additive test only; depends on REM-001 being applied first.

---

### REM-025 · TC-006 — Strengthen TestOrdersHandler_GET assertion

| Field | Detail |
|---|---|
| **Finding ID** | TC-006 |
| **Severity** | Medium |
| **Affected file** | `sample-app/internal/handlers/handlers_test.go` line 34 |

**Proposed change**
```go
// Before
if len(orders) == 0 {
    t.Fatal("expected at least one order in response")
}

// After
expected := handlers.SeedOrders()
if got, want := len(orders), len(expected); got != want {
    t.Fatalf("order count = %d, want %d", got, want)
}
for i, o := range orders {
    if o.ID != expected[i].ID {
        t.Errorf("orders[%d].ID = %q, want %q", i, o.ID, expected[i].ID)
    }
}
```

**Reason**  
A `len > 0` assertion would pass even if `SeedOrders()` returned a single zero-value `Order{}`. Comparing against the fixture output directly ties the test to the actual contract and catches mutations.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestOrdersHandler_GET -v
```

**Potential risk**  
None. Strengthening an existing assertion cannot introduce false positives on correct code.

---

### REM-026 · DC-002 — Document DB_DSN in requirements.md

| Field | Detail |
|---|---|
| **Finding ID** | DC-002 |
| **Severity** | Medium |
| **Affected file** | `sample-app/docs/requirements.md` section 3 |

**Proposed change**

Add to the environment variable table:
```markdown
| DB_DSN | string | postgres://localhost:5432/orders?sslmode=disable | PostgreSQL connection string. Currently read but not consumed; reserved for future database integration. |
```

**Reason**  
FR-34 requires all configuration to be documented. `DB_DSN` is read in `main.go` but absent from docs, making it invisible to operators who inspect the documentation before deployment.

**Validation command**
```bash
grep -n "DB_DSN" sample-app/docs/requirements.md
# Should return at least one match
```

**Potential risk**  
None. Documentation change only.

---

### REM-027 · DC-009 — Switch Dockerfile final stage to distroless

| Field | Detail |
|---|---|
| **Finding ID** | DC-009 |
| **Severity** | Medium |
| **Affected file** | `sample-app/deploy/Dockerfile` line 23 |

**Proposed change**
```dockerfile
# Before
FROM debian:bookworm-slim

# After
FROM gcr.io/distroless/base-debian12:nonroot
```

**Reason**  
Docker Gate D-04 recommends a minimal final stage. `debian:bookworm-slim` includes a shell, package manager, and system utilities that are unnecessary and increase the attack surface. A distroless image contains only the Go runtime libraries and the CA certificate bundle.

**Validation command**
```bash
docker build -t orders-api:distroless -f sample-app/deploy/Dockerfile sample-app/
docker run --rm orders-api:distroless   # should start, not produce "sh: not found"
docker scout cves orders-api:distroless | grep -i critical
```

**Potential risk**  
Medium. Distroless images have no shell; `docker exec` for debugging is not possible. Use a debug-variant image (`gcr.io/distroless/base-debian12:debug`) in development environments. The switch also requires that `adduser`/`addgroup` commands in REM-007 are executed in the builder stage, not the final stage, since distroless has no shell to run them. Adjust the Dockerfile accordingly.

---

### REM-028 · DP-007 — Pin deployment image tag

| Field | Detail |
|---|---|
| **Finding ID** | DP-007 |
| **Severity** | Medium |
| **Affected file** | `sample-app/deploy/deployment.yaml` line 25 |

**Proposed change**
```yaml
# Before
image: orders-api:latest

# After
image: orders-api:1.0.0   # replace with the actual release tag or SHA digest
```

**Reason**  
`:latest` makes deployments non-reproducible. Kubernetes will not re-pull a `latest` image if it is already cached, silently running stale code. Using a specific version or digest enables reliable rollbacks with `kubectl rollout undo`.

**Validation command**
```bash
kubectl apply --dry-run=client -f sample-app/deploy/deployment.yaml
```

**Potential risk**  
Low. The chosen tag must exist in the registry before applying. Coordinate the image push and manifest update in the same CI pipeline step.

---

## Part 4 — Low Severity (8)

---

### REM-029 · CQ-006 — Remove or connect DB_DSN variable

| Field | Detail |
|---|---|
| **Finding ID** | CQ-006 |
| **Severity** | Low |
| **Affected file** | `sample-app/cmd/server/main.go` lines 27–30 |

**Proposed change**  
Remove the `DB_DSN` block entirely until database integration is implemented:
```go
// Remove:
// dbDSN := os.Getenv("DB_DSN")
// if dbDSN == "" {
//     dbDSN = "postgres://localhost:5432/orders?sslmode=disable"
// }
```

**Reason**  
Dead code misleads readers into thinking a database connection is being established. Combined with REM-026 (doc update), removing it keeps the code honest. Restore when the database client is actually wired up.

**Validation command**
```bash
cd sample-app && go build ./... && go vet ./...
```

**Potential risk**  
None. The variable is not used anywhere.

---

### REM-030 · CQ-007 — Replace LOG_LEVEL with slog or remove it

| Field | Detail |
|---|---|
| **Finding ID** | CQ-007 |
| **Severity** | Low |
| **Affected file** | `sample-app/cmd/server/main.go` lines 32–35 |

**Proposed change**  
Replace the `log` package with `log/slog` (available since Go 1.21):
```go
import "log/slog"

level := slog.LevelInfo
if os.Getenv("LOG_LEVEL") == "debug" {
    level = slog.LevelDebug
}
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
```

**Reason**  
The standard `log` package does not support level filtering; `LOG_LEVEL` is read but has no effect, which is confusing for operators. `slog` is the standard library solution and requires no external dependencies.

**Validation command**
```bash
cd sample-app && LOG_LEVEL=debug go run ./cmd/server/... 2>&1 | head -5
```

**Potential risk**  
Low. Existing log lines will need to be ported from `log.Printf` to `slog.Info`/`slog.Debug`. Output format changes (JSON vs plain text) should be verified against any log-scraping pipelines.

---

### REM-031 · CQ-010 — Expose server timeouts as environment variables

| Field | Detail |
|---|---|
| **Finding ID** | CQ-010 |
| **Severity** | Low |
| **Affected file** | `sample-app/cmd/server/main.go` lines 42–44 |

**Proposed change**
```go
readTimeout  := envDuration("READ_TIMEOUT_SEC",  10)
writeTimeout := envDuration("WRITE_TIMEOUT_SEC", 10)
idleTimeout  := envDuration("IDLE_TIMEOUT_SEC",  60)

// helper:
func envDuration(key string, defaultSec int) time.Duration {
    if v := os.Getenv(key); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            return time.Duration(n) * time.Second
        }
    }
    return time.Duration(defaultSec) * time.Second
}
```

**Reason**  
Hardcoded timeouts require a rebuild and redeployment to tune. In containerised environments, operators expect all behaviour to be configurable via environment variables.

**Validation command**
```bash
cd sample-app && READ_TIMEOUT_SEC=5 go run ./cmd/server/... &
sleep 1 && kill $! && wait
```

**Potential risk**  
Very low. Defaults are unchanged; the change only activates when the variables are explicitly set.

---

### REM-032 · CQ-013 — Establish context propagation pattern

| Field | Detail |
|---|---|
| **Finding ID** | CQ-013 |
| **Severity** | Low |
| **Affected file** | `sample-app/internal/handlers/orders.go` line 37 |

**Proposed change**

No functional change required today. Add a comment establishing the pattern for future callers:
```go
func OrdersHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // propagate to DB/service calls when added
    _ = ctx            // remove when ctx is first used
    // ...
```

**Reason**  
Establishing the pattern now prevents future contributors from ignoring context, which would make cancellation and deadline propagation impossible once a database layer is added.

**Validation command**
```bash
cd sample-app && go vet ./internal/handlers/...
```

**Potential risk**  
None. Comment and no-op assignment only.

---

### REM-033 · CQ-014 — Adopt consistent log format

| Field | Detail |
|---|---|
| **Finding ID** | CQ-014 |
| **Severity** | Low |
| **Affected file** | `sample-app/cmd/server/main.go` line 49 |

**Proposed change**  
This is resolved as a by-product of REM-030 (adopt `slog`). If `slog` is deferred, change all log lines to use the same format (either all plain text or all key=value structured).

**Reason**  
Mixed formats make log parsing fragile when regex-based scrapers or index tools are in use.

**Validation command**
```bash
cd sample-app && go run ./cmd/server/... 2>&1 | head -5
```

**Potential risk**  
None beyond REM-030.

---

### REM-034 · TC-007 — Expand method-not-allowed tests with table-driven tests

| Field | Detail |
|---|---|
| **Finding ID** | TC-007 |
| **Severity** | Low |
| **Affected file** | `sample-app/internal/handlers/handlers_test.go` line 83 |

**Proposed change**
```go
func TestReadinessHandler_MethodNotAllowed(t *testing.T) {
    methods := []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
    for _, method := range methods {
        t.Run(method, func(t *testing.T) {
            req := httptest.NewRequest(method, "/ready", nil)
            rec := httptest.NewRecorder()
            handlers.ReadinessHandler(rec, req)
            if rec.Code != http.StatusMethodNotAllowed {
                t.Errorf("%s: status = %d, want %d", method, rec.Code, http.StatusMethodNotAllowed)
            }
        })
    }
}
```

**Reason**  
Testing only POST leaves PUT/PATCH/DELETE uncovered. A future refactor that accidentally allows PATCH would not be caught.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestReadinessHandler_MethodNotAllowed -v
```

**Potential risk**  
None. Additive test only.

---

### REM-035 · TC-008 — Assert Content-Type in TestReadinessHandler_GET

| Field | Detail |
|---|---|
| **Finding ID** | TC-008 |
| **Severity** | Low |
| **Affected file** | `sample-app/internal/handlers/handlers_test.go` line 63 |

**Proposed change**
```go
// Add after status code assertion:
ct := resp.Header.Get("Content-Type")
if !strings.HasPrefix(ct, "application/json") {
    t.Errorf("Content-Type = %q, want application/json", ct)
}
```

**Reason**  
`TestOrdersHandler_GET` already asserts `Content-Type`; `TestReadinessHandler_GET` inconsistently does not despite the handler setting it. Parity ensures the contract is enforced for both handlers.

**Validation command**
```bash
cd sample-app && go test ./internal/handlers/... -run TestReadinessHandler_GET -v
```

**Potential risk**  
None. Strengthening an existing test.

---

### REM-036 · DP-008 — Migrate hardcoded env vars to a ConfigMap

| Field | Detail |
|---|---|
| **Finding ID** | DP-008 |
| **Severity** | Low |
| **Affected file** | `sample-app/deploy/deployment.yaml` line 28 + new `sample-app/deploy/configmap.yaml` |

**Proposed change**

Create `sample-app/deploy/configmap.yaml`:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: orders-api-config
data:
  PORT: "8080"
  LOG_LEVEL: "info"
```

In `deployment.yaml`, replace:
```yaml
env:
  - name: PORT
    value: '8080'
  - name: LOG_LEVEL
    value: 'info'
```
with:
```yaml
envFrom:
  - configMapRef:
      name: orders-api-config
```

**Reason**  
Embedding configuration in the deployment manifest couples configuration to the deployment lifecycle. A ConfigMap allows values to be changed per-environment without modifying the core deployment definition.

**Validation command**
```bash
kubectl apply --dry-run=client -f sample-app/deploy/configmap.yaml
kubectl apply --dry-run=client -f sample-app/deploy/deployment.yaml
```

**Potential risk**  
Low. The ConfigMap must be applied before the Deployment. If it is missing, pods will fail to start with `configmap "orders-api-config" not found`. Ensure CI applies manifests in dependency order.

---

## Consolidated Execution Order

Apply changes in this sequence to avoid dependency failures:

| Step | Remediation | Finding(s) | Severity | File(s) |
|---|---|---|---|---|
| 1 | REM-009 | CQ-003 | High | `sample-app/go.mod` |
| 2 | REM-006 | DC-007 / DP-002 | High | `Dockerfile` line 10 |
| 3 | REM-007 | DC-008 / DP-001 | High | `Dockerfile` final stage |
| 4 | REM-001 | DC-001 | Blocker | `router.go` + new `health.go` |
| 5 | REM-002 | DC-003 / DP-006 | Blocker | `deployment.yaml` pod spec |
| 6 | REM-003 | DC-004 / DP-005 | Blocker | `deployment.yaml` container spec |
| 7 | REM-004 | DC-005 / DP-003 / DP-004 | Blocker | `deployment.yaml` container spec |
| 8 | REM-005 | DC-006 | Blocker | `deployment.yaml` container spec |
| 9 | REM-008 | CQ-001 | High | `readiness.go` line 22 |
| 10 | REM-016 | CQ-002 | Medium | `orders.go` line 46 |
| 11 | REM-017 | CQ-004 | Medium | `orders.go` (part of REM-016) |
| 12 | REM-018 | CQ-005 | Medium | `readiness.go` line 21 |
| 13 | REM-019 | CQ-008 | Medium | `orders.go` line 39 |
| 14 | REM-020 | CQ-009 | Medium | `readiness.go` line 17 |
| 15 | REM-023 | CQ-015 | Medium | `main.go` line 19 |
| 16 | REM-022 | CQ-012 | Medium | `main.go` line 51 |
| 17 | REM-021 | CQ-011 | Medium | `main.go` line 62 |
| 18 | REM-010 | TC-001 / DC-010 | High | `handlers_test.go` |
| 19 | REM-011 | TC-002 / TC-005 | High | `handlers_test.go` |
| 20 | REM-012 | TC-004 | High | new `main_test.go` |
| 21 | REM-024 | TC-003 | Medium | `handlers_test.go` |
| 22 | REM-025 | TC-006 | Medium | `handlers_test.go` |
| 23 | REM-026 | DC-002 | Medium | `docs/requirements.md` |
| 24 | REM-027 | DC-009 | Medium | `Dockerfile` line 23 |
| 25 | REM-028 | DP-007 | Medium | `deployment.yaml` line 25 |
| 26 | REM-029 | CQ-006 | Low | `main.go` lines 27–30 |
| 27 | REM-030 | CQ-007 | Low | `main.go` lines 32–35 |
| 28 | REM-031 | CQ-010 | Low | `main.go` lines 42–44 |
| 29 | REM-032 | CQ-013 | Low | `orders.go` line 37 |
| 30 | REM-033 | CQ-014 | Low | `main.go` line 49 |
| 31 | REM-034 | TC-007 | Low | `handlers_test.go` line 83 |
| 32 | REM-035 | TC-008 | Low | `handlers_test.go` line 63 |
| 33 | REM-036 | DP-008 | Low | `deployment.yaml` + new `configmap.yaml` |

---

*No changes have been made to any sample-app file. Awaiting explicit approval before implementation begins.*
