package printer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/danblok/grab/internal/algorithm"
	"github.com/fatih/color"
)

const NAME_PLACEHOLDER = "-------- %s --------\n"

// stores range of a pattern
type pattern struct {
	start, end int
}

// Wrapper function around FprintDefault to make it output to stdout
func PrintDefault(r io.Reader, name string, patterns []string) {
	FprintDefault(r, os.Stdout, name, patterns)
}

// Wrapper function around FprintQuite to make it output to stdout
func PrintQuite(r io.Reader, patterns []string) {
	FprintQuite(r, os.Stdout, patterns)
}

// Wrapper function around FprintMinimum to make it output to stdout
func PrintNonHuman(r io.Reader, patterns []string) {
	FprintNonHuman(r, os.Stdout, patterns)
}

// Prints name of the output, then prints lines where patterns were found with line indicators.
// The found patterns are emphasized with red color.
func FprintDefault(r io.Reader, w io.Writer, inputName string, patterns []string) {
	fmt.Fprintf(w, NAME_PLACEHOLDER, inputName)
	lineCount := 1
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		ps := accumulatePatterns(line, patterns)
		ps = mergePatternRanges(ps)
		if len(ps) > 0 {
			for i := range ps {
				if i == 0 {
					fmt.Fprintf(w, "%d: %s%v", lineCount, line[:ps[i].start], color.RedString(line[ps[i].start:ps[i].end+1]))
				} else {
					fmt.Fprintf(w, "%s%v", line[ps[i-1].end+1:ps[i].start], color.RedString(line[ps[i].start:ps[i].end+1]))
				}
			}
			fmt.Fprintf(w, "%s\n", line[ps[len(ps)-1].end+1:])
		}
		lineCount++
	}
}

// Prints lines where the patterns were found without line indicators and input name.
// The found patterns are emphasized with red color.
func FprintQuite(r io.Reader, w io.Writer, patterns []string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		ps := accumulatePatterns(line, patterns)
		ps = mergePatternRanges(ps)
		if len(ps) > 0 {
			for i := range ps {
				if i == 0 {
					fmt.Fprintf(w, "%s%v", line[:ps[i].start], color.RedString(line[ps[i].start:ps[i].end+1]))
				} else {
					fmt.Fprintf(w, "%s%v", line[ps[i-1].end+1:ps[i].start], color.RedString(line[ps[i].start:ps[i].end+1]))
				}
			}
			fmt.Fprintf(w, "%s\n", line[ps[len(ps)-1].end+1:])
		}
	}
}

// Prints lines in format: "<line_number> <pattern1_start_idx> <pattern1_end_idx> <pattern2_start_idx> <pattern2_end_idx>"
// The found patterns are emphasized with red color.
func FprintNonHuman(r io.Reader, w io.Writer, patterns []string) {
	lineCount := 1
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		ps := accumulatePatterns(line, patterns)
		if len(ps) > 0 {
			for i := range ps {
				if i == 0 {
					fmt.Fprintf(w, "%d %s", lineCount, color.RedString(fmt.Sprintf("%d %d", ps[i].start, ps[i].end)))
				} else {
					fmt.Fprintf(w, " %s", color.RedString(fmt.Sprintf("%d %d", ps[i].start, ps[i].end)))
				}
			}
			fmt.Fprintln(w)
		}
		lineCount++
	}
}

// Searches for and accumulates all patterns on the line.
// Returned slice is sorted.
func accumulatePatterns(line string, patterns []string) []pattern {
	idxs := make([]pattern, 0)
	// accumulate patterns indicies found on the line
	for _, p := range patterns {
		patternIdxs := algorithm.Search(line, p)
		mappedIdxs := make([]pattern, len(patternIdxs))
		for i := range patternIdxs {
			mappedIdxs[i] = pattern{patternIdxs[i], patternIdxs[i] + len(p) - 1}
		}

		idxs = append(idxs, mappedIdxs...)
	}

	// sort the indicies by start index for better colorization
	slices.SortFunc(idxs, func(a, b pattern) int {
		if a.start == b.start {
			return a.end - b.end
		} else {
			return a.start - b.start
		}
	})

	return idxs
}

// Merges the ranges of patterns.
func mergePatternRanges(patterns []pattern) []pattern {
	if len(patterns) == 0 {
		return []pattern{}
	}

	// iteratively swaps ranges
	j := 0
	for i := 1; i < len(patterns); i++ {
		if patterns[i].start <= patterns[j].end {
			patterns[j].end = patterns[i].end
		} else {
			j++
			patterns[j] = patterns[i]
		}
	}

	patterns = patterns[:j+1]
	return patterns
}
