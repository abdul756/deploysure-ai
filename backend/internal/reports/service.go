package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Service reads report files from disk. No database is used.
type Service struct {
	reportsDir string
}

// NewService creates a Service that reads from reportsDir.
func NewService(reportsDir string) *Service {
	return &Service{reportsDir: reportsDir}
}

// LoadFindings reads and parses a JSON findings file.
// The filename must not escape the reports directory.
func (s *Service) LoadFindings(ctx context.Context, filename string) ([]Finding, error) {
	path, err := s.safePath(filename)
	if err != nil {
		return nil, err
	}
	return s.loadFindingsFromPath(ctx, path)
}

// FindingsBefore returns the main combined findings-before.json.
func (s *Service) FindingsBefore(ctx context.Context) ([]Finding, error) {
	return s.LoadFindings(ctx, "findings-before.json")
}

// FindingsAfter returns findings-after.json if it exists.
func (s *Service) FindingsAfter(ctx context.Context) ([]Finding, error) {
	return s.LoadFindings(ctx, "findings-after.json")
}

// ReportBefore returns the raw content of the before readiness report.
// It returns the file content as a string (markdown or plain text).
func (s *Service) ReportBefore(ctx context.Context) (string, error) {
	return s.loadText(ctx, "release-readiness-before.md")
}

// ReportAfter returns the raw content of the after readiness report.
func (s *Service) ReportAfter(ctx context.Context) (string, error) {
	return s.loadText(ctx, "release-readiness-after.md")
}

// Comparison builds a ComparisonResult from the before and after findings.
func (s *Service) Comparison(ctx context.Context) (*ComparisonResult, error) {
	before, err := s.FindingsBefore(ctx)
	if err != nil {
		return nil, fmt.Errorf("comparison: before: %w", err)
	}
	after, err := s.FindingsAfter(ctx)
	if err != nil {
		// After file may not exist yet — treat as empty rather than an error.
		after = []Finding{}
	}
	return &ComparisonResult{
		Before: before,
		After:  after,
		Summary: ComparisonSummary{
			Before: CountSeverities(before),
			After:  CountSeverities(after),
		},
	}, nil
}

// safePath validates that filename contains no path traversal and returns the
// absolute path within the reports directory.
func (s *Service) safePath(filename string) (string, error) {
	// Reject any component that looks like a traversal attempt.
	if strings.Contains(filename, "..") || strings.ContainsAny(filename, "/\\") {
		return "", fmt.Errorf("reports: invalid filename %q", filename)
	}
	clean := filepath.Join(s.reportsDir, filepath.Base(filename))
	// Confirm the resolved path is still inside the reports directory.
	absReports, err := filepath.Abs(s.reportsDir)
	if err != nil {
		return "", fmt.Errorf("reports: cannot resolve reports dir: %w", err)
	}
	absClean, err := filepath.Abs(clean)
	if err != nil {
		return "", fmt.Errorf("reports: cannot resolve path: %w", err)
	}
	if !strings.HasPrefix(absClean, absReports+string(filepath.Separator)) {
		return "", fmt.Errorf("reports: path %q escapes reports directory", filename)
	}
	return absClean, nil
}

func (s *Service) loadFindingsFromPath(ctx context.Context, path string) ([]Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("reports: file not found: %w", err)
		}
		return nil, fmt.Errorf("reports: read %q: %w", path, err)
	}
	var findings []Finding
	if err := json.Unmarshal(data, &findings); err != nil {
		return nil, fmt.Errorf("reports: parse %q: %w", path, err)
	}
	return findings, nil
}

func (s *Service) loadText(ctx context.Context, filename string) (string, error) {
	path, err := s.safePath(filename)
	if err != nil {
		return "", err
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("reports: file not found: %w", err)
		}
		return "", fmt.Errorf("reports: read %q: %w", path, err)
	}
	return string(data), nil
}
