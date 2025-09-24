package cmd

import (
	"github.com/jokruger/pwgen/internal/generator"
	"github.com/spf13/cobra"
)

// generateCmd provides an explicit "generate" subcommand. Invoking `pwgen generate`
// is equivalent to calling the root command directly (e.g. `pwgen`), but can
// improve clarity in scripts.
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate a password / key (same behavior as root command)",
	Long: `Explicit generation subcommand.

Examples:
  pwgen generate
  pwgen generate --length 32 --min-number 2 --min-symbol 2
  pwgen generate --format appkey --segments 5 --segment-length 6
  pwgen generate --format guid
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Reuse existing flag variables defined in root.go
		opts := generator.Options{
			Format:        generator.Format(flagFormat),
			Length:        flagLength,
			UseLower:      flagUseLower,
			UseUpper:      flagUseUpper,
			UseNumber:     flagUseNumber,
			UseSymbol:     flagUseSymbol,
			MinLower:      flagMinLower,
			MinUpper:      flagMinUpper,
			MinNumber:     flagMinNumber,
			MinSymbol:     flagMinSymbol,
			Segments:      flagSegments,
			SegmentLength: flagSegmentLength,
		}

		out, err := generator.Generate(opts)
		if err != nil {
			return err
		}
		cmd.Println(out)
		return nil
	},
}

func init() {
	// Flags are inherited from rootCmd; no need to redefine here.
	rootCmd.AddCommand(generateCmd)
}
