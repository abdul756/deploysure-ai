# DeploySure AI

DeploySure AI is an AI-powered deployment risk analysis tool built for the IBM Watsonx Hackathon 2025.

## Overview

DeploySure AI helps engineering teams identify deployment risks before they reach production. It analyses deployment configurations, compares them against historical deployment data, and generates actionable risk reports powered by IBM Watsonx.

## Project Structure

```
deploysure-ai/
├── frontend/        # Plain HTML/CSS/JS web interface
├── backend/         # Go backend — REST API and Watsonx integration
├── sample-app/      # Demo application used to showcase risk analysis
├── docs/            # Architecture diagrams and design documents
├── reports/         # Generated risk-analysis reports
├── bob_sessions/    # Exported Bob AI session logs (required for submission)
└── evidence/        # Demo screenshots and supporting evidence
```

## Getting Started

1. Copy environment variables:

   ```bash
   cp .env.example .env
   # Edit .env with your IBM Cloud and Watsonx credentials
   ```

2. Start the backend:

   ```bash
   # Instructions will be added once the backend is implemented
   ```

3. Open `frontend/index.html` in your browser.

## Security

See [SECURITY.MD](SECURITY.MD) for credential management guidelines. Never commit `.env` or any credentials.

## Documentation

- [BOB_USAGE.md](BOB_USAGE.md) – How Bob AI was used in this project
- [PROBLEM_STATEMENT.md](PROBLEM_STATEMENT.md) – Problem being solved
- [SOLUTION_STATEMENT.md](SOLUTION_STATEMENT.md) – Proposed solution
- [DEMO_SCRIPT.md](DEMO_SCRIPT.md) – Walkthrough for live demo

---

IBM Watsonx Hackathon 2025
