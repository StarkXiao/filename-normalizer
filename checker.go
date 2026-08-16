package normalizer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const DefaultMaxLength = 255

var validNamePattern = regexp.MustCompile(`^[\pL\pN._ -]+$`)

type Config struct {
	MaxLength    int
	AllowedRegex string
	Recursive    bool
}

type Issue string

const (
	IllegalCharacters  Issue = "illegal_characters"
	ConsecutiveSpaces  Issue = "consecutive_spaces"
	DuplicateExtension Issue = "duplicate_extension"
	CaseConflict       Issue = "case_conflict"
	NameTooLong        Issue = "name_too_long"
)

type Finding struct {
	Path       string
	Name       string
	Issue      Issue
	Suggestion string
}

type Checker struct {
	config  Config
	allowed *regexp.Regexp
}

func New(config Config) (*Checker, error) {
	if config.MaxLength <= 0 {
		config.MaxLength = DefaultMaxLength
	}
	var allowed *regexp.Regexp
	if config.AllowedRegex != "" {
		var err error
		allowed, err = regexp.Compile(config.AllowedRegex)
		if err != nil {
			return nil, fmt.Errorf("compile allowed regex: %w", err)
		}
	}
	return &Checker{config: config, allowed: allowed}, nil
}

func (c *Checker) Scan(root string) ([]Finding, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat scan path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("scan path is not a directory: %s", root)
	}
	var paths []string
	walkErr := filepath.Walk(root, func(path string, entry os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root && entry.IsDir() {
			return nil
		}
		if entry.IsDir() {
			if !c.config.Recursive {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Mode().IsRegular() {
			paths = append(paths, path)
		}
		return nil
	})
	if walkErr != nil {
		return nil, fmt.Errorf("walk scan path: %w", walkErr)
	}
	sort.Strings(paths)
	caseNames := make(map[string][]string)
	for _, path := range paths {
		caseNames[strings.ToLower(filepath.Base(path))] = append(caseNames[strings.ToLower(filepath.Base(path))], path)
	}
	var findings []Finding
	for _, path := range paths {
		name := filepath.Base(path)
		for _, issue := range c.issues(name) {
			findings = append(findings, Finding{Path: path, Name: name, Issue: issue, Suggestion: Suggestion(name, issue, c.config.MaxLength)})
		}
		if len(caseNames[strings.ToLower(name)]) > 1 {
			findings = append(findings, Finding{Path: path, Name: name, Issue: CaseConflict, Suggestion: Suggestion(name, CaseConflict, c.config.MaxLength)})
		}
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Path == findings[j].Path {
			return findings[i].Issue < findings[j].Issue
		}
		return findings[i].Path < findings[j].Path
	})
	return findings, nil
}

func (c *Checker) issues(name string) []Issue {
	var issues []Issue
	if !utf8.ValidString(name) || !validNamePattern.MatchString(name) || strings.ContainsAny(name, `/\\:*?"<>|`) {
		issues = append(issues, IllegalCharacters)
	}
	if strings.Contains(name, "  ") {
		issues = append(issues, ConsecutiveSpaces)
	}
	ext := filepath.Ext(name)
	if ext != "" && strings.HasSuffix(strings.TrimSuffix(name, ext), ext) {
		issues = append(issues, DuplicateExtension)
	}
	if c.allowed != nil && !c.allowed.MatchString(name) {
		issues = append(issues, IllegalCharacters)
	}
	if len(name) > c.config.MaxLength {
		issues = append(issues, NameTooLong)
	}
	return uniqueIssues(issues)
}

func uniqueIssues(in []Issue) []Issue {
	seen := map[Issue]bool{}
	out := make([]Issue, 0, len(in))
	for _, v := range in {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func Suggestion(name string, issue Issue, max int) string {
	result := name
	switch issue {
	case IllegalCharacters:
		result = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune("._ -", r) {
				return r
			}
			return '_'
		}, result)
	case ConsecutiveSpaces:
		result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")
	case DuplicateExtension:
		ext := filepath.Ext(result)
		result = strings.TrimSuffix(result, ext)
	case CaseConflict:
		result = strings.ToLower(result)
	}
	if max > 0 && utf8.RuneCountInString(result) > max {
		ext := filepath.Ext(result)
		base := strings.TrimSuffix(result, ext)
		keep := max - utf8.RuneCountInString(ext)
		if keep < 1 {
			keep = 1
		}
		runes := []rune(base)
		if len(runes) > keep {
			result = string(runes[:keep]) + ext
		}
	}
	return result
}

func WriteCSV(w io.Writer, findings []Finding) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"path", "name", "issue", "suggestion"}); err != nil {
		return err
	}
	for _, f := range findings {
		if err := cw.Write([]string{f.Path, f.Name, string(f.Issue), f.Suggestion}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func RenamePlan(findings []Finding) []Finding {
	seen := map[string]bool{}
	out := make([]Finding, 0, len(findings))
	for _, f := range findings {
		if f.Suggestion == f.Name {
			continue
		}
		target := filepath.Join(filepath.Dir(f.Path), f.Suggestion)
		if seen[target] {
			continue
		}
		seen[target] = true
		out = append(out, f)
	}
	return out
}
