# DEMO_SCRIPT.md

## Demo Script – DeploySure AI

This document provides a step-by-step walkthrough for demonstrating DeploySure AI to judges.

### Duration

Approximately 5 minutes.

### Prerequisites

- Backend running locally or on IBM Cloud
- Frontend open in browser
- Sample deployment configurations ready in `sample-app/deploy/`

---

### Step 1 – Introduce the Problem (30 s)

Briefly explain the problem: teams deploy without knowing the risk, and incidents follow.

### Step 2 – Show the Interface (30 s)

Open `frontend/index.html`. Walk through the upload form and explain what it does.

### Step 3 – Submit a Low-Risk Configuration (1 min)

Paste a simple, low-risk deployment config from `sample-app/deploy/`. Click **Analyze Risk**. Show the resulting report and highlight the low risk score.

### Step 4 – Submit a High-Risk Configuration (1 min)

Replace the config with a high-risk example. Show that DeploySure AI correctly identifies multiple risks and provides actionable mitigations.

### Step 5 – Show the Generated Report (1 min)

Walk through the exported risk report. Point out the risk score, identified patterns, and recommended actions.

### Step 6 – Summarise (1 min)

Recap the value: faster, consistent, AI-driven risk assessment at the point of change.

---

_Script will be refined once the full implementation is complete._
