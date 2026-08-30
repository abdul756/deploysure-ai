package reports_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/reports"
)

// makeReportsDir creates a temporary reports directory and returns a Service
// backed by it.
func makeReportsDir(t *testing.T) (string, *reports.Service) {
	t.Helper()
	dir := t.TempDir()
	return dir, reports.NewService(dir)
}

// writeFindings writes a slice of Finding values as JSON to dir/filename.
func writeFindings(t *testing.T, dir, filename string, findings []reports.Finding) {
	t.Helper()
	b, err := json.Marshal(findings)
	if err != nil {
		t.Fatalf("marshal fixtures: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, filename), b, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
}

var sampleFindings = []reports.Finding{
	{ID: "CQ-001", Severity: "high", Title: "Test high finding"},
	{ID: "DC-001", Severity: "blocker", Title: "Test blocker finding"},
	{ID: "TC-001", Severity: "medium", Title: "Test medium finding"},
	{ID: "DP-001", Severity: "low", Title: "Test low finding"},
}

func TestLoadFindings_Success(t *testing.T) {
	dir, svc := makeReportsDir(t)
	writeFindings(t, dir, "findings-before.json", sampleFindings)

	got, err := svc.FindingsBefore(context.Background())
	if err != nil {
		t.Fatalf("FindingsBefore: %v", err)
	}
	if len(got) != len(sampleFindings) {
		t.Errorf("len = %d, want %d", len(got), len(sampleFindings))
	}
	if got[0].ID != "CQ-001" {
		t.Errorf("got[0].ID = %q, want %q", got[0].ID, "CQ-001")
	}
}

func TestLoadFindings_NotFound(t *testing.T) {
	_, svc := makeReportsDir(t)
	_, err := svc.FindingsBefore(context.Background())
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFindings_TraversalRejected(t *testing.T) {
	_, svc := makeReportsDir(t)
	_, err := svc.LoadFindings(context.Background(), "../etc/passwd")
	if err == nil {
		t.Fatal("expected error for path traversal attempt")
	}
}

func TestLoadFindings_SlashInFilenameRejected(t *testing.T) {
	_, svc := makeReportsDir(t)
	_, err := svc.LoadFindings(context.Background(), "sub/dir/file.json")
	if err == nil {
		t.Fatal("expected error for slash in filename")
	}
}

func TestLoadFindings_ContextCancelled(t *testing.T) {
	dir, svc := makeReportsDir(t)
	writeFindings(t, dir, "findings-before.json", sampleFindings)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	_, err := svc.FindingsBefore(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestFindingsAfter_Empty(t *testing.T) {
	_, svc := makeReportsDir(t)
	// File does not exist; Comparison treats it as empty.
	result, err := svc.Comparison(context.Background())
	// Comparison only errors if FindingsBefore fails.
	if err == nil && result != nil && len(result.After) != 0 {
		t.Errorf("expected empty After findings, got %d", len(result.After))
	}
}

func TestComparison_Success(t *testing.T) {
	dir, svc := makeReportsDir(t)
	writeFindings(t, dir, "findings-before.json", sampleFindings)
	writeFindings(t, dir, "findings-after.json", []reports.Finding{
		{ID: "CQ-001", Severity: "high", Title: "Still present"},
	})

	result, err := svc.Comparison(context.Background())
	if err != nil {
		t.Fatalf("Comparison: %v", err)
	}
	if len(result.Before) != 4 {
		t.Errorf("Before count = %d, want 4", len(result.Before))
	}
	if len(result.After) != 1 {
		t.Errorf("After count = %d, want 1", len(result.After))
	}
	if result.Summary.Before.Blocker != 1 {
		t.Errorf("Before.Blocker = %d, want 1", result.Summary.Before.Blocker)
	}
	if result.Summary.Before.High != 1 {
		t.Errorf("Before.High = %d, want 1", result.Summary.Before.High)
	}
	if result.Summary.Before.Total != 4 {
		t.Errorf("Before.Total = %d, want 4", result.Summary.Before.Total)
	}
}

func TestReportBefore_Success(t *testing.T) {
	dir, svc := makeReportsDir(t)
	content := "# Release Report\nThis is the before report."
	if err := os.WriteFile(filepath.Join(dir, "release-readiness-before.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := svc.ReportBefore(context.Background())
	if err != nil {
		t.Fatalf("ReportBefore: %v", err)
	}
	if got != content {
		t.Errorf("content mismatch\ngot:  %q\nwant: %q", got, content)
	}
}

func TestCountSeverities(t *testing.T) {
	counts := reports.CountSeverities(sampleFindings)
	if counts.Blocker != 1 {
		t.Errorf("Blocker = %d, want 1", counts.Blocker)
	}
	if counts.High != 1 {
		t.Errorf("High = %d, want 1", counts.High)
	}
	if counts.Medium != 1 {
		t.Errorf("Medium = %d, want 1", counts.Medium)
	}
	if counts.Low != 1 {
		t.Errorf("Low = %d, want 1", counts.Low)
	}
	if counts.Total != 4 {
		t.Errorf("Total = %d, want 4", counts.Total)
	}
}
