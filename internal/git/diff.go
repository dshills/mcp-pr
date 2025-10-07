package git

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

// FileDiff represents changes to a single file
type FileDiff struct {
	OldPath   string // Path before change (or empty for new files)
	NewPath   string // Path after change (or empty for deleted files)
	Hunks     []Hunk // Individual change hunks
	IsNew     bool   // True if file is newly added
	IsDeleted bool   // True if file is deleted
}

// Hunk represents a contiguous block of changes
type Hunk struct {
	OldStart int      // Starting line in old file
	OldLines int      // Number of lines in old file
	NewStart int      // Starting line in new file
	NewLines int      // Number of lines in new file
	Lines    []string // Actual diff lines (with +/- prefix)
}

// Parse parses unified diff format into structured data
func Parse(diffText string) ([]FileDiff, error) {
	if diffText == "" {
		return []FileDiff{}, nil
	}

	var fileDiffs []FileDiff
	var currentFile *FileDiff
	var currentHunk *Hunk

	scanner := bufio.NewScanner(strings.NewReader(diffText))

	// Regex patterns
	diffHeaderPattern := regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	oldFilePattern := regexp.MustCompile(`^--- (.*)$`)
	newFilePattern := regexp.MustCompile(`^\+\+\+ (.*)$`)
	hunkHeaderPattern := regexp.MustCompile(`^@@ -(\d+),?(\d*) \+(\d+),?(\d*) @@`)

	for scanner.Scan() {
		line := scanner.Text()

		// File header: diff --git a/path b/path
		if matches := diffHeaderPattern.FindStringSubmatch(line); matches != nil {
			// Save previous file if exists
			if currentFile != nil {
				fileDiffs = append(fileDiffs, *currentFile)
			}

			currentFile = &FileDiff{
				OldPath: matches[1],
				NewPath: matches[2],
				Hunks:   []Hunk{},
			}
			currentHunk = nil
			continue
		}

		// Old file marker: --- a/path
		if oldFilePattern.MatchString(line) {
			if currentFile != nil && strings.Contains(line, "/dev/null") {
				currentFile.IsNew = true
			}
			continue
		}

		// New file marker: +++ b/path
		if newFilePattern.MatchString(line) {
			if currentFile != nil && strings.Contains(line, "/dev/null") {
				currentFile.IsDeleted = true
			}
			continue
		}

		// Hunk header: @@ -1,5 +1,7 @@
		if matches := hunkHeaderPattern.FindStringSubmatch(line); matches != nil {
			// Save previous hunk if exists
			if currentHunk != nil && currentFile != nil {
				currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
			}

			oldStart := parseInt(matches[1])
			oldLines := parseInt(matches[2])
			if oldLines == 0 {
				oldLines = 1 // Default to 1 if not specified
			}
			newStart := parseInt(matches[3])
			newLines := parseInt(matches[4])
			if newLines == 0 {
				newLines = 1
			}

			currentHunk = &Hunk{
				OldStart: oldStart,
				OldLines: oldLines,
				NewStart: newStart,
				NewLines: newLines,
				Lines:    []string{},
			}
			continue
		}

		// Diff content lines (start with +, -, or space)
		if currentHunk != nil && (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, " ")) {
			currentHunk.Lines = append(currentHunk.Lines, line)
		}
	}

	// Save last hunk and file
	if currentHunk != nil && currentFile != nil {
		currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
	}
	if currentFile != nil {
		fileDiffs = append(fileDiffs, *currentFile)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning diff: %w", err)
	}

	return fileDiffs, nil
}

// parseInt safely parses an integer from a string, returning 0 if empty
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var n int
	_, _ = fmt.Sscanf(s, "%d", &n) // Ignoring error as we default to 0
	return n
}

// FormatForReview converts parsed diff back to a simplified format for LLM review
func FormatForReview(fileDiffs []FileDiff) string {
	var builder strings.Builder

	for _, file := range fileDiffs {
		builder.WriteString(fmt.Sprintf("File: %s\n", file.NewPath))

		switch {
		case file.IsNew:
			builder.WriteString("Status: New file\n")
		case file.IsDeleted:
			builder.WriteString("Status: Deleted\n")
		default:
			builder.WriteString("Status: Modified\n")
		}

		builder.WriteString("\n")

		for _, hunk := range file.Hunks {
			builder.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
				hunk.OldStart, hunk.OldLines, hunk.NewStart, hunk.NewLines))

			for _, line := range hunk.Lines {
				builder.WriteString(line)
				builder.WriteString("\n")
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
