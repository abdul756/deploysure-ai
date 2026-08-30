// Command analyze reads reports/findings-before.json, sends the findings to
// IBM watsonx.ai Granite for risk assessment, validates the response, and
// writes the result to reports/granite-risk-assessment.json and
// reports/granite-risk-assessment.md (the analysis text, human-readable).
//
// Run from the repository root:
//
//	go run ./backend/cmd/analyze
//
// Configuration is read exclusively from environment variables — no .env file.
//
// Environment variables:
//
//	IBM_CLOUD_API_KEY    IBM Cloud API key used for IAM token exchange (required, never logged)
//	WATSONX_API_KEY      Legacy alias for IBM_CLOUD_API_KEY (never logged)
//	WATSONX_PROJECT_ID   IBM watsonx.ai project ID (required, never logged)
//	WATSONX_URL          watsonx.ai base URL (default: https://us-south.ml.cloud.ibm.com)
//	WATSONX_MODEL_ID     Granite model ID (default: ibm/granite-13b-instruct-v2)
//	REPORTS_DIR          Directory containing report files (default: reports)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/config"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/reports"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/watsonx"
)

// graniteRiskAssessment is the structured metadata written to granite-risk-assessment.json.
// The analysis text itself is written to granite-risk-assessment.md so the JSON
// stays compact and the markdown is human-readable without escape sequences.
type graniteRiskAssessment struct {
	GeneratedAt  string `json:"generated_at"`
	ModelID      string `json:"model_id"`
	InputFile    string `json:"input_file"`
	FindingCount int    `json:"finding_count"`
	AnalysisFile string `json:"analysis_file"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("level=fatal msg=\"invalid configuration\" err=%q", err.Error())
	}
	if cfg.WatsonxAPIKey == "" {
		log.Fatal("level=fatal msg=\"IBM_CLOUD_API_KEY (or WATSONX_API_KEY) must be set\"")
	}
	if cfg.WatsonxProjectID == "" {
		log.Fatal("level=fatal msg=\"WATSONX_PROJECT_ID must be set\"")
	}

	// 1. Read reports/findings-before.json.
	svc := reports.NewService(cfg.ReportsDir)
	ctx := context.Background()

	findings, err := svc.FindingsBefore(ctx)
	if err != nil {
		log.Fatalf("level=fatal msg=\"cannot read findings-before.json\" err=%q", err.Error())
	}
	log.Printf("level=info msg=\"loaded findings\" count=%d", len(findings))

	// Serialise findings to structured text for the prompt.
	findingsJSON, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		log.Fatalf("level=fatal msg=\"cannot marshal findings\" err=%q", err.Error())
	}
	inputText := fmt.Sprintf("Deployment findings (JSON):\n\n%s", string(findingsJSON))

	// 2. Send to Granite with a 5-minute deadline for the full round-trip
	//    (IAM exchange + inference).
	client := watsonx.NewClient(cfg.WatsonxAPIKey, cfg.WatsonxProjectID, cfg.WatsonxURL, cfg.WatsonxModelID)
	log.Printf("level=info msg=\"sending findings to watsonx\" model=%s", client.ModelID())

	analyzeCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	analysis, err := client.Analyze(analyzeCtx, inputText)
	if err != nil {
		log.Fatalf("level=fatal msg=\"analysis failed\" err=%q", err.Error())
	}

	// 3. Validate the response — must be non-empty after trimming.
	analysis = strings.TrimSpace(analysis)
	if analysis == "" {
		log.Fatal("level=fatal msg=\"model returned an empty analysis\"")
	}
	log.Printf("level=info msg=\"analysis received\" bytes=%d", len(analysis))

	// 4a. Save the analysis text as markdown — newlines are literal, easy to read.
	mdPath := filepath.Join(cfg.ReportsDir, "granite-risk-assessment.md")
	mdContent := fmt.Sprintf("# Granite Risk Assessment\n\n**Generated:** %s  \n**Model:** %s  \n**Input:** %s  \n**Findings:** %d\n\n---\n\n%s\n",
		time.Now().UTC().Format(time.RFC3339),
		client.ModelID(),
		filepath.Join(cfg.ReportsDir, "findings-before.json"),
		len(findings),
		analysis,
	)
	if err := os.WriteFile(mdPath, []byte(mdContent), 0o644); err != nil {
		log.Fatalf("level=fatal msg=\"cannot write markdown\" path=%q err=%q", mdPath, err.Error())
	}
	log.Printf("level=info msg=\"analysis markdown saved\" path=%q", mdPath)

	// 4b. Save compact metadata JSON — no long strings, stays pretty.
	assessment := graniteRiskAssessment{
		GeneratedAt:  time.Now().UTC().Format(time.RFC3339),
		ModelID:      client.ModelID(),
		InputFile:    filepath.Join(cfg.ReportsDir, "findings-before.json"),
		FindingCount: len(findings),
		AnalysisFile: mdPath,
	}

	out, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		log.Fatalf("level=fatal msg=\"cannot marshal assessment\" err=%q", err.Error())
	}

	outputPath := filepath.Join(cfg.ReportsDir, "granite-risk-assessment.json")
	if err := os.WriteFile(outputPath, out, 0o644); err != nil {
		log.Fatalf("level=fatal msg=\"cannot write assessment\" path=%q err=%q", outputPath, err.Error())
	}

	log.Printf("level=info msg=\"risk assessment saved\" path=%q", outputPath)
}
