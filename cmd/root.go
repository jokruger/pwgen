package cmd

import (
	"fmt"
	"os"

	"github.com/jokruger/pwgen/internal/generator"
	"github.com/spf13/cobra"
)

var (
	flagLength        int
	flagFormat        string
	flagUseLower      bool
	flagUseUpper      bool
	flagUseNumber     bool
	flagUseSymbol     bool
	flagMinLower      int
	flagMinUpper      int
	flagMinNumber     int
	flagMinSymbol     int
	flagSegments      int
	flagSegmentLength int
)

var rootCmd = &cobra.Command{
	Use:   "pwgen",
	Short: "Generate secure passwords, app keys, or GUIDs",
	Long: `pwgen generates cryptographically secure passwords / keys.

Formats:
  generic (default) - random characters according to selected classes
  appkey            - segmented key (e.g. XXXX-XXXX-XXXX)
  guid              - RFC 4122 UUID v4

Character classes can be toggled and minimum counts enforced.`,
	RunE: runGenerate,
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&flagLength, "length", "l", 16, "Total password length (generic format)")
	rootCmd.PersistentFlags().StringVarP(&flagFormat, "format", "f", "generic", "Output format: generic|appkey|guid")
	rootCmd.PersistentFlags().BoolVar(&flagUseLower, "lower", true, "Include lowercase letters")
	rootCmd.PersistentFlags().BoolVar(&flagUseUpper, "upper", true, "Include uppercase letters")
	rootCmd.PersistentFlags().BoolVar(&flagUseNumber, "number", true, "Include numbers")
	rootCmd.PersistentFlags().BoolVar(&flagUseSymbol, "symbol", true, "Include symbols")

	rootCmd.PersistentFlags().IntVar(&flagMinLower, "min-lower", 0, "Minimum lowercase letters")
	rootCmd.PersistentFlags().IntVar(&flagMinUpper, "min-upper", 0, "Minimum uppercase letters")
	rootCmd.PersistentFlags().IntVar(&flagMinNumber, "min-number", 0, "Minimum numbers")
	rootCmd.PersistentFlags().IntVar(&flagMinSymbol, "min-symbol", 0, "Minimum symbols")

	rootCmd.PersistentFlags().IntVar(&flagSegments, "segments", 4, "Number of segments (appkey format)")
	rootCmd.PersistentFlags().IntVar(&flagSegmentLength, "segment-length", 4, "Length of each segment (appkey format)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
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

	pw, err := generator.Generate(opts)
	if err != nil {
		return err
	}
	fmt.Println(pw)
	return nil
}
