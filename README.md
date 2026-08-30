# DeploySure AI

DeploySure AI is an AI-powered deployment risk analysis tool built for the IBM TechXchange 2026 Pre-conference Dev Day Hackathon.

## Overview

DeploySure analyzes source code, tests, Dockerfiles, Kubernetes manifests and
release documentation to identify release risks and documentation
inconsistencies. It uses IBM Granite (via watsonx.ai) to produce structured
findings and a release-readiness score.

## Project Structure

```
deploysure-ai/
├── frontend/        # Plain HTML/CSS/JS web interface
├── backend/         # Go backend — REST API and watsonx.ai integration
├── sample-app/      # Demo orders-API used as the subject of risk analysis
├── docs/            # Architecture diagrams and design documents
├── reports/         # Generated risk-analysis reports (before & after)
├── bob_sessions/    # Exported Bob AI session logs (required for submission)
└── evidence/        # Demo screenshots and supporting evidence
```

This is a **Go workspace** project (`go.work`). Both `backend/` and `sample-app/`
are separate Go modules managed together in the workspace.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- An IBM Cloud account with a [watsonx.ai](https://www.ibm.com/watsonx) project
- IBM Cloud API key ([create one here](https://cloud.ibm.com/iam/apikeys))

## End-to-End Setup & Run

### 1. Clone and configure environment

```bash
git clone <repo-url>
cd deploysure-ai

cp .env.example .env
# Edit .env and fill in your credentials:
#   IBM_CLOUD_API_KEY  — your IBM Cloud API key
#   WATSONX_PROJECT_ID — your watsonx.ai project ID
#   WATSONX_URL        — e.g. https://us-south.ml.cloud.ibm.com
#   WATSONX_MODEL_ID   — e.g. ibm/granite-13b-instruct-v2
```

### 2. Configure environment variables

Both the **backend** and the **sample app** load a `.env` file from the working
directory automatically if one is present. Export credentials into your shell
**or** place them in `.env` — both work:

```bash
export IBM_CLOUD_API_KEY="<your-api-key>"
export WATSONX_PROJECT_ID="<your-project-id>"
export WATSONX_URL="https://us-south.ml.cloud.ibm.com"
export WATSONX_MODEL_ID="ibm/granite-13b-instruct-v2"
```

> **All optional overrides with defaults:**
>
> | Variable | Default | Description |
> |---|---|---|
> | `PORT` | `8080` | HTTP port the backend server listens on |
> | `REPORTS_DIR` | `reports` | Directory containing JSON/Markdown report files |
> | `READ_TIMEOUT_SEC` | `30` | HTTP read timeout (seconds) |
> | `WRITE_TIMEOUT_SEC` | `30` | HTTP write timeout (seconds) |
> | `IDLE_TIMEOUT_SEC` | `60` | HTTP idle timeout (seconds) |
> | `SHUTDOWN_TIMEOUT_SEC` | `30` | Graceful-shutdown grace period (seconds) |

### 3. Start the backend server

Run from the **repository root**:

```bash
go run ./backend/cmd/server
```

You should see:

```
level=info msg="watsonx integration enabled" url=https://us-south.ml.cloud.ibm.com
level=info msg="server starting" addr=:8080
```

### 4. Open the frontend

Navigate to [http://localhost:8080](http://localhost:8080) in your browser.

The backend serves the frontend directly from the `frontend/` directory. No
separate web server is needed.

### 5. Load and analyze findings

Use the dashboard buttons in order:

| Button | What it does |
|---|---|
| **Load Before Analysis** | Fetches `reports/findings-before.json` and populates the findings table |
| **Analyze with Granite** | Sends findings to IBM watsonx.ai (Granite) and displays the AI risk assessment |
| **Load After Analysis** | Fetches `reports/findings-after.json` (post-remediation findings) |
| **Compare Before vs After** | Side-by-side comparison of readiness scores and finding counts |

### 6. (Optional) Run the one-shot analysis CLI

Generate a fresh Granite risk assessment and write it to `reports/`:

```bash
go run ./backend/cmd/analyze
```

Output files written:
- `reports/granite-risk-assessment.md` — human-readable analysis
- `reports/granite-risk-assessment.json` — compact metadata

### 7. (Optional) Run the sample app

The sample app is a demo orders-API used as the subject of analysis. It loads
a `.env` file automatically if one is present in the working directory
([`github.com/joho/godotenv`](https://github.com/joho/godotenv)), falling back
to real environment variables when no file exists (e.g. in a container).

```bash
# Create a local .env for the sample app (optional)
cat > sample-app/.env <<'EOF'
PORT=9090
LOG_LEVEL=debug
DB_DSN=postgres://localhost:5432/orders?sslmode=disable
EOF

# Run locally (reads sample-app/.env automatically)
cd sample-app && go run ./cmd/server

# Run tests
cd sample-app && go test ./...

# Build and run via Docker (from sample-app/)
docker build -f deploy/Dockerfile -t orders-api .
docker run -p 8080:8080 orders-api
```

> **Note:** The Docker image does not bundle a `.env` file. When running in a
> container or Kubernetes, set environment variables directly — the app falls
> back to them automatically when no `.env` file is found.

## Verify the backend is healthy

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

Both endpoints return `200 OK` when the server is running correctly.

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/findings/before` | Pre-remediation findings (JSON) |
| `GET` | `/api/v1/findings/after` | Post-remediation findings (JSON) |
| `GET` | `/api/v1/reports/before` | Pre-remediation release readiness report (Markdown) |
| `GET` | `/api/v1/reports/after` | Post-remediation release readiness report (Markdown) |
| `GET` | `/api/v1/comparison` | Before vs after comparison summary |
| `POST` | `/api/v1/granite/analyze` | Trigger Granite AI risk analysis via watsonx.ai |

## Environment Variable Reference

### Backend (`backend/`)

Loads `.env` from the working directory if present, then reads environment
variables. Values already set in the environment take precedence over `.env`.

| Variable | Required | Default | Description |
|---|---|---|---|
| `IBM_CLOUD_API_KEY` | Yes* | — | IBM Cloud API key for IAM token exchange |
| `WATSONX_API_KEY` | Yes* | — | Legacy alias for `IBM_CLOUD_API_KEY` |
| `WATSONX_PROJECT_ID` | Yes | — | watsonx.ai project ID |
| `WATSONX_URL` | No | `https://us-south.ml.cloud.ibm.com` | watsonx.ai base URL |
| `WATSONX_MODEL_ID` | No | `ibm/granite-13b-instruct-v2` | Granite model identifier |
| `PORT` | No | `8080` | TCP port the server listens on |
| `REPORTS_DIR` | No | `reports` | Report files directory |
| `READ_TIMEOUT_SEC` | No | `30` | HTTP read timeout (seconds) |
| `WRITE_TIMEOUT_SEC` | No | `30` | HTTP write timeout (seconds) |
| `IDLE_TIMEOUT_SEC` | No | `60` | HTTP idle timeout (seconds) |
| `SHUTDOWN_TIMEOUT_SEC` | No | `30` | Graceful-shutdown period (seconds) |

\* Either `IBM_CLOUD_API_KEY` or `WATSONX_API_KEY` must be set to enable Granite analysis.

### Sample app (`sample-app/`)

Loads `sample-app/.env` if present, otherwise reads from the process environment.

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | TCP port the server listens on |
| `LOG_LEVEL` | `info` | Log level (informational only — printed at startup) |
| `DB_DSN` | `postgres://localhost:5432/orders?sslmode=disable` | Database DSN (read but not used in demo) |

## Security

See [SECURITY.MD](SECURITY.MD) for credential management guidelines. Never
commit `.env` or any credentials to version control — `.env` is listed in
`.gitignore`.

## Documentation

- [BOB_USAGE.md](BOB_USAGE.md) – How Bob AI was used in this project
- [PROBLEM_STATEMENT.md](PROBLEM_STATEMENT.md) – Problem being solved
- [SOLUTION_STATEMENT.md](SOLUTION_STATEMENT.md) – Proposed solution
- [DEMO_SCRIPT.md](DEMO_SCRIPT.md) – Walkthrough for live demo

---
