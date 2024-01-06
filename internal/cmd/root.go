package cmd

import (
	"fmt"
	"os"

	"github.com/danblok/grab/internal/printer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&files, "files", "f", []string{}, "Attach files where the pattern will be searched")
	rootCmd.PersistentFlags().BoolVarP(&quite, "quite", "q", false, "Print the output without line indicators")
	rootCmd.PersistentFlags().BoolVarP(&nonhuman, "nonhuman", "n", false, "Print the output in format \"<line_number> <pattern1_start_idx> <pattern1_end_idx> <pattern2_start_idx> <pattern2_end_idx>\"")
}

var rootCmd = &cobra.Command{
	Use:   "grab",
	Short: "Searches for patterns in text",
	Long: `Searches for patterns in stdin and files. Outputs lines where
the pattern was found. The patterns are marked with red color.`,
	Args: cobra.MinimumNArgs(1),
	Run:  run,
}

var (
	files    []string
	quite    bool
	nonhuman bool
)

func run(_ *cobra.Command, args []string) {
	// search in stdin
	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Something went wrong with stdin")
	}
	if stdinStat.Mode()&os.ModeCharDevice == 0 {
		switch {
		case quite:
			printer.PrintQuite(os.Stdin, args)
		case nonhuman:
			printer.PrintNonHuman(os.Stdin, args)
		default:
			printer.PrintDefault(os.Stdin, "stdin", args)

		}
	}

	// search in files
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Couldn't process the file: %s\n", file)
			continue
		}

		switch {
		case quite:
			printer.PrintQuite(f, args)
		case nonhuman:
			printer.PrintNonHuman(f, args)
		default:
			printer.PrintDefault(f, f.Name(), args)
		}
	}
}

func Execute() {
	rootCmd.Execute()
}
