# Solution Statement

## Overview

**DeploySure AI** is an AI-powered release-readiness platform that automates
the full pre-release verification workflow. It combines IBM Bob Agent mode,
specialised subagents, and watsonx.ai Granite to analyse a candidate release
across multiple dimensions simultaneously and produce a prioritised, auditable
risk report — with human-approved remediation built in.

## How It Works

```
Developer pushes code
        │
        ▼
┌───────────────────────┐
│ Bob Agent             |
| (Orchestrator)        │
│ reads requirements,   │
│ runbook, and manifests│
└──────────┬────────────┘
           │  spawns parallel subagents
     ┌─────┼──────┬──────────┬──────────┐
     ▼     ▼      ▼          ▼          ▼
 Code   Test   Docker    Kubernetes   Docs-to-Code
Quality  Gap   Review     Review     Consistency
Agent   Agent  Agent      Agent        Agent
     └─────┴──────┴──────────┴──────────┘
                    │
                    ▼
         watsonx.ai Granite
         risk prioritisation
                    │
                    ▼
         Human review + approval
                    │
                    ▼
         Remediation applied
                    │
                    ▼
         Final validation pass
```

## Capabilities Demonstrated

| Capability | How DeploySure Uses It |
|-----------|----------------------|
| **IBM Bob Agent mode** | Orchestrates the full workflow as a single Bob task |
| **Document understanding** | Bob reads `docs/requirements.md` and `docs/release-runbook.md` to ground every check in documented standards |
| **Parallel tasks** | Code quality, test gap, Docker, Kubernetes, and docs-consistency analyses run concurrently via `spawn_subagent` |
| **Specialised subagents** | Each analysis domain has a focused subagent with a narrow scope and its own context |
| **Code-quality analysis** | Runs `gofmt`, `go vet`, and lint rules; reports violations with file and line references |
| **Unit-test gap analysis** | Compares exported functions/handlers against existing `*_test.go` coverage |
| **Docker review** | Audits `Dockerfile` for root execution, pinned base images, exposed secrets, and layer hygiene |
| **Kubernetes review** | Validates Deployment and Service manifests for resource limits, liveness/readiness probes, and non-root security context |
| **Docs-to-code consistency** | Compares documented API routes and environment variables against actual source code |
| **Human-approved remediation** | Surfaces a consolidated risk report; requires explicit human sign-off before applying fixes |
| **Final validation** | Re-runs all checks after remediation to confirm the release is clean |
| **watsonx.ai Granite risk prioritisation** | Sends the aggregated findings to Granite; receives severity scores and a ranked remediation order |

## The Sample Application

To make the demonstration concrete, DeploySure analyses a synthetic Go
service (`sample-app/`) that intentionally contains seeded defects. The
application exposes `GET /api/orders`, `GET /health`, and `GET /ready` using
the standard `net/http` package, reads configuration from environment variables,
and supports graceful shutdown. It ships with unit tests, a `Dockerfile`, and
Kubernetes Deployment and Service manifests.

## Value Delivered

- **Speed**: A full pre-release check that takes 30–60 minutes manually
  completes in under 2 minutes.
- **Consistency**: Every release is measured against the same documented
  standards with no human variability.
- **Auditability**: Every finding, every approval, and every fix is recorded
  in an exportable Bob session log.
- **Risk focus**: Granite-prioritised findings ensure engineers spend time on
  the highest-impact issues first.

## Out of Scope

Post-release monitoring, CI/CD pipeline integration, and multi-repository
orchestration are not addressed in the current scope.
