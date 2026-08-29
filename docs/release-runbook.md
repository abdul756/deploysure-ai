# Release Runbook

## Purpose

This runbook defines the mandatory checks that EVERY release of the sample
application must pass before it is considered release-ready. DeploySure AI
uses this document as its authoritative standard during automated review.

---

## 1. Pre-Flight Checklist

Before running automated analysis, the release engineer must confirm:

- [ ] All feature work is merged to the release branch.
- [ ] The `CHANGELOG` entry for this version is written.
- [ ] Environment-variable documentation in `docs/requirements.md` is up to date.
- [ ] No `.env` or secrets files are staged for commit.

---

## 2. Code-Quality Gate

All items in this section are **mandatory**. A single failure blocks release.

### 2.1 Formatting — `gofmt`

```bash
gofmt -l ./...
```

**Pass condition**: zero files listed (no formatting differences detected).

### 2.2 Static Analysis — `go vet`

```bash
go vet ./...
```

**Pass condition**: exits with code 0; no warnings or errors on stdout/stderr.

### 2.3 Unit Tests — `go test`

```bash
go test -race -count=1 -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**Pass conditions**:
- Exits with code 0 (all tests pass).
- `-race` flag enabled to detect data races.
- Package-level coverage ≥ 80 % for every package.

---

## 3. API Contract Checks

### 3.1 Required Endpoints

Every build artifact MUST serve the following routes:

| Method | Path | Expected Status | Expected Body |
|--------|------|-----------------|---------------|
| GET | `/api/orders` | 200 | JSON array |
| GET | `/health` | 200 | `{"status":"ok"}` |
| GET | `/ready` | 200 | `{"status":"ready"}` |

### 3.2 Smoke Test

```bash
# Start the service locally
PORT=8080 go run ./cmd/server &
sleep 1

curl -sf http://localhost:8080/health  | grep '"status":"ok"'
curl -sf http://localhost:8080/ready   | grep '"status":"ready"'
curl -sf http://localhost:8080/api/orders | jq 'type == "array"'

kill %1
```

**Pass condition**: all three `curl` commands succeed without error.

---

## 4. Docker Gate

All items are **mandatory**.

| # | Check | Pass Condition |
|---|-------|---------------|
| D-01 | Non-root execution | `USER` directive present; value is NOT `root` or `0` |
| D-02 | Pinned base image | Base image tag is a specific version or digest, NOT `:latest` |
| D-03 | No committed secrets | No `ENV`, `ARG`, or `RUN` line matches a credential pattern (password, token, key, secret) |
| D-04 | Minimal final stage | Multi-stage build used; final stage is `scratch` or a distroless image |
| D-05 | Build succeeds | `docker build` exits with code 0 |
| D-06 | Container starts | `docker run --rm` produces a healthy response from `/health` within 10 seconds |

---

## 5. Kubernetes Gate

All items are **mandatory**.

| # | Check | Pass Condition |
|---|-------|---------------|
| K-01 | Resource requests | Every container defines `resources.requests.cpu` and `resources.requests.memory` |
| K-02 | Resource limits | Every container defines `resources.limits.cpu` and `resources.limits.memory` |
| K-03 | Liveness probe | Every container defines `livenessProbe` targeting `GET /health` |
| K-04 | Readiness probe | Every container defines `readinessProbe` targeting `GET /ready` |
| K-05 | Non-root security context | `securityContext.runAsNonRoot: true` on the pod or container spec |
| K-06 | Read-only root filesystem | `securityContext.readOnlyRootFilesystem: true` on the container spec |
| K-07 | Manifest applies cleanly | `kubectl apply --dry-run=client -f k8s/` exits with code 0 |

---

## 6. Security Gate

| # | Check | Tool / Method |
|---|-------|--------------|
| S-01 | No committed credentials | `git grep -i "password\|secret\|token\|apikey"` returns no results in tracked files |
| S-02 | `.env` not tracked | `git ls-files .env` returns empty |
| S-03 | Secrets not in image layers | `docker history --no-trunc` inspected for credential patterns |

---

## 7. Documentation Consistency Gate

| # | Check | Pass Condition |
|---|-------|---------------|
| DC-01 | All documented routes exist in code | Zero undeclared routes found by docs-to-code subagent |
| DC-02 | All documented env vars read in code | Zero undeclared env vars found by docs-to-code subagent |
| DC-03 | All code routes documented | Zero undocumented routes found by docs-to-code subagent |

---

## 8. Release Sign-Off

After all gates above pass, the release engineer must:

1. Review the DeploySure AI risk report at `reports/latest.md`.
2. Confirm Granite-assigned severity — no CRITICAL or HIGH findings open.
3. Approve remediation in the Bob Agent session.
4. Confirm the final validation pass is clean.
5. Tag the release:

```bash
git tag -s v<MAJOR>.<MINOR>.<PATCH> -m "Release v<MAJOR>.<MINOR>.<PATCH>"
git push origin v<MAJOR>.<MINOR>.<PATCH>
```

---

## 9. Rollback Procedure

If a post-release issue is detected:

1. Identify the previous stable image tag in the container registry.
2. Update the Kubernetes Deployment image field to the previous tag.
3. Apply: `kubectl set image deployment/orders-api orders-api=<image>:<prev-tag>`
4. Confirm pods roll to running: `kubectl rollout status deployment/orders-api`
5. Open a post-incident review issue within 24 hours.
