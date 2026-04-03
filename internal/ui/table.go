package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	Bold    = color.New(color.Bold).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Dim     = color.New(color.Faint).SprintFunc()
)

// PrintTable prints a formatted table to stdout.
func PrintTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		fmt.Println(Dim("  No results found."))
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Cap description column at 50 chars
	for i := range widths {
		if widths[i] > 60 {
			widths[i] = 60
		}
	}

	// Print header
	headerLine := "  "
	for i, h := range headers {
		headerLine += fmt.Sprintf("%-*s  ", widths[i], Bold(h))
	}
	fmt.Println(headerLine)
	fmt.Println("  " + strings.Repeat("-", len(headerLine)-2))

	// Print rows
	for _, row := range rows {
		line := "  "
		for i, cell := range row {
			if i >= len(widths) {
				break
			}
			display := cell
			if len(display) > widths[i] {
				display = display[:widths[i]-3] + "..."
			}
			line += fmt.Sprintf("%-*s  ", widths[i], display)
		}
		fmt.Println(line)
	}
}

// PrintSuccess prints a success message.
func PrintSuccess(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", Green("✓"), msg)
}

// PrintError prints an error message.
func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", Red("✗"), msg)
}

// PrintInfo prints an info message.
func PrintInfo(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", Cyan("i"), msg)
}

// Truncate shortens a string to maxLen, adding "..." if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
