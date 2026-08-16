package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	normalizer "filename-normalizer"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan")
	recursive := flag.Bool("recursive", true, "scan nested directories")
	max := flag.Int("max-length", normalizer.DefaultMaxLength, "maximum filename length in runes")
	allowed := flag.String("allowed-regex", "", "regular expression filenames must match")
	csvPath := flag.String("csv", "", "write findings CSV to this file")
	dryRun := flag.Bool("dry-run", false, "print proposed renames without changing files")
	apply := flag.Bool("rename", false, "apply proposed renames")
	flag.Parse()
	if *dryRun && *apply {
		fmt.Fprintln(os.Stderr, "--dry-run and --rename cannot be used together")
		os.Exit(2)
	}
	c, err := normalizer.New(normalizer.Config{MaxLength: *max, AllowedRegex: *allowed, Recursive: *recursive})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	findings, err := c.Scan(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("scanned %s: %d finding(s)\n", *dir, len(findings))
	if *csvPath != "" {
		f, e := os.Create(*csvPath)
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		e = normalizer.WriteCSV(f, findings)
		closeErr := f.Close()
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		if closeErr != nil {
			fmt.Fprintln(os.Stderr, closeErr)
			os.Exit(1)
		}
	}
	plan := normalizer.RenamePlan(findings)
	if *dryRun || *apply {
		for _, item := range plan {
			target := filepath.Join(filepath.Dir(item.Path), item.Suggestion)
			fmt.Printf("rename %s -> %s\n", item.Path, target)
			if *apply {
				if _, e := os.Stat(target); e == nil {
					fmt.Fprintf(os.Stderr, "skip existing target: %s\n", target)
					continue
				}
				if e := os.Rename(item.Path, target); e != nil {
					fmt.Fprintf(os.Stderr, "rename %s: %v\n", item.Path, e)
				}
			}
		}
	}
	if len(findings) > 0 {
		os.Exit(1)
	}
}
