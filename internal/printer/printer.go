package printer

import (
	"bufio"
	"fmt"
	"io"
	"slices"

	"github.com/danblok/grab/internal/algorithm"
	"github.com/fatih/color"
)

// stores range of a pattern
type pattern struct {
	start, end int
}

// Prints name of the output, then prints lines where patterns were found with line indicators.
// The found patterns are emphasized with red color.
func PrintDefault(reader io.Reader, inputName string, patterns []string) {
	fmt.Printf("-------- %s --------\n", inputName)
	lineCount := 1
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		ps := accumulatePatterns(line, patterns)
		ps = mergePatternsRanges(ps)
		if len(ps) > 0 {
			for i := range ps {
				if i == 0 {
					fmt.Printf("%d: %s%v", lineCount, line[:ps[i].start], color.RedString(line[ps[i].start:ps[i].end]))
				} else {
					fmt.Printf("%s%v", line[ps[i-1].end:ps[i].start], color.RedString(line[ps[i].start:ps[i].end]))
				}
			}
			fmt.Printf("%s\n", line[ps[len(ps)-1].end:])
		}
		lineCount++
	}
}

// Prints lines where the patterns were found without line indicators and input name.
// The found patterns are emphasized with red color.
func PrintQuite(reader io.Reader, patterns []string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		ps := accumulatePatterns(line, patterns)
		ps = mergePatternsRanges(ps)
		if len(ps) > 0 {
			for i := range ps {
				if i == 0 {
					fmt.Printf("%s%v", line[:ps[i].start], color.RedString(line[ps[i].start:ps[i].end]))
				} else {
					fmt.Printf("%s%v", line[ps[i-1].end:ps[i].start], color.RedString(line[ps[i].start:ps[i].end]))
				}
			}
			fmt.Printf("%s\n", line[ps[len(ps)-1].end:])
		}
	}
}

// Prints lines in format: "<line_number> <pattern1_start_idx> <pattern1_end_idx> <pattern2_start_idx> <pattern2_end_idx>"
// The found patterns are emphasized with red color.
func PrintMinimum(reader io.Reader, patterns []string) {
	lineCount := 1
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		ps := accumulatePatterns(line, patterns)
		if len(ps) > 0 {
			for i := range ps {
				if i == 0 {
					fmt.Printf("%d %s", lineCount, color.RedString(fmt.Sprintf("%d %d", ps[i].start, ps[i].end)))
				} else {
					fmt.Printf(" %s", color.RedString(fmt.Sprintf("%d %d", ps[i].start, ps[i].end)))
				}
			}
			fmt.Println()
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
			mappedIdxs[i] = pattern{patternIdxs[i], patternIdxs[i] + len(p)}
		}

		idxs = append(idxs, mappedIdxs...)
	}

	// sort the indicies by start index for better colorization
	slices.SortFunc(idxs, func(a, b pattern) int {
		if a.start < b.start {
			return -1
		} else if a.start > b.start {
			return 1
		} else {
			return 0
		}
	})

	return idxs
}

// Merges the ranges of patterns.
func mergePatternsRanges(patterns []pattern) []pattern {
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
