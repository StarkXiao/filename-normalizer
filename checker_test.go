package normalizer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanFindsIssues(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"bad  name?.txt", "report.txt.txt", "Photo.JPG", "photo.jpg", strings.Repeat("a", 12) + ".txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	c, err := New(Config{MaxLength: 15})
	if err != nil {
		t.Fatal(err)
	}
	findings, err := c.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) < 4 {
		t.Fatalf("got %d findings, want at least 4", len(findings))
	}
}

func TestWriteCSVAndPlan(t *testing.T) {
	findings := []Finding{{Path: "/tmp/a  b.txt", Name: "a  b.txt", Issue: ConsecutiveSpaces, Suggestion: "a b.txt"}}
	var b strings.Builder
	if err := WriteCSV(&b, findings); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "path,name,issue,suggestion") || !strings.Contains(b.String(), "consecutive_spaces") {
		t.Fatal(b.String())
	}
	if got := RenamePlan(findings); len(got) != 1 {
		t.Fatalf("got %d plan items", len(got))
	}
}

func TestNewRejectsInvalidRegex(t *testing.T) {
	if _, err := New(Config{AllowedRegex: "["}); err == nil {
		t.Fatal("expected regex error")
	}
}
