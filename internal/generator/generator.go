package generator

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Format identifies the overall output style.
type Format string

const (
	FormatGeneric Format = "generic" // Random characters per enabled classes
	FormatAppKey  Format = "appkey"  // Segmented groups (e.g. XXXX-XXXX-XXXX)
	FormatGUID    Format = "guid"    // UUID v4
)

// Character class rune slices.
var (
	lowerChars  = []rune("abcdefghijklmnopqrstuvwxyz")
	upperChars  = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numberChars = []rune("0123456789")
	symbolChars = []rune("!@#$%^&*()-_=+[]{};:,.?/<>~")
)

// Options holds generation parameters.
type Options struct {
	// General
	Format Format
	Length int // Used only for generic format (appkey derives total from segments)

	// Toggle character classes
	UseLower  bool
	UseUpper  bool
	UseNumber bool
	UseSymbol bool

	// Minimum counts per class (only applied to generic/appkey)
	MinLower  int
	MinUpper  int
	MinNumber int
	MinSymbol int

	// App key specific
	Segments      int // Number of segments (e.g. 4 -> XXXX-XXXX-XXXX-XXXX)
	SegmentLength int // Characters per segment
}

// Generate produces a password / key string according to the provided Options.
func Generate(o Options) (string, error) {
	switch o.Format {
	case FormatGUID:
		return generateUUIDv4()
	case FormatAppKey:
		return generateAppKey(o)
	case FormatGeneric:
		return generateGeneric(o)
	default:
		return "", fmt.Errorf("unknown format: %s", o.Format)
	}
}

// generateGeneric handles the basic random password generation.
func generateGeneric(o Options) (string, error) {
	if o.Length <= 0 {
		return "", errors.New("length must be > 0")
	}

	classSets, minima, err := collectClassesAndValidate(o)
	if err != nil {
		return "", err
	}

	totalMin := 0
	for _, v := range minima {
		totalMin += v
	}
	if totalMin > o.Length {
		return "", fmt.Errorf("sum of minimum counts (%d) exceeds requested length %d", totalMin, o.Length)
	}

	// Allocate slice for runes
	out := make([]rune, 0, o.Length)

	// Satisfy minima first
	for i, set := range classSets {
		minCount := minima[i]
		for c := 0; c < minCount; c++ {
			r, err := randomRune(set)
			if err != nil {
				return "", err
			}
			out = append(out, r)
		}
	}

	// Fill the remainder
	all := concatClasses(classSets)
	for len(out) < o.Length {
		r, err := randomRune(all)
		if err != nil {
			return "", err
		}
		out = append(out, r)
	}

	// Shuffle to avoid predictable grouping of minima
	shuffleRunes(out)

	return string(out), nil
}

// generateAppKey creates a segmented key of Segments * SegmentLength characters.
func generateAppKey(o Options) (string, error) {
	if o.Segments <= 0 {
		return "", errors.New("segments must be > 0")
	}
	if o.SegmentLength <= 0 {
		return "", errors.New("segment-length must be > 0")
	}

	total := o.Segments * o.SegmentLength
	// Reuse generic logic
	tmp := o
	tmp.Length = total
	tmp.Format = FormatGeneric

	core, err := generateGeneric(tmp)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.Grow(total + (o.Segments - 1))
	for i, r := range core {
		if i > 0 && i%o.SegmentLength == 0 {
			b.WriteRune('-')
		}
		b.WriteRune(r)
	}
	return b.String(), nil
}

// collectClassesAndValidate builds the enabled character sets and validates minima.
func collectClassesAndValidate(o Options) ([][]rune, []int, error) {
	var sets [][]rune
	var mins []int

	if o.UseLower {
		sets = append(sets, lowerChars)
		mins = append(mins, o.MinLower)
	} else if o.MinLower > 0 {
		return nil, nil, errors.New("min-lower specified but lowercase disabled")
	}

	if o.UseUpper {
		sets = append(sets, upperChars)
		mins = append(mins, o.MinUpper)
	} else if o.MinUpper > 0 {
		return nil, nil, errors.New("min-upper specified but uppercase disabled")
	}

	if o.UseNumber {
		sets = append(sets, numberChars)
		mins = append(mins, o.MinNumber)
	} else if o.MinNumber > 0 {
		return nil, nil, errors.New("min-number specified but numbers disabled")
	}

	if o.UseSymbol {
		sets = append(sets, symbolChars)
		mins = append(mins, o.MinSymbol)
	} else if o.MinSymbol > 0 {
		return nil, nil, errors.New("min-symbol specified but symbols disabled")
	}

	if len(sets) == 0 {
		return nil, nil, errors.New("no character classes enabled")
	}

	for _, v := range mins {
		if v < 0 {
			return nil, nil, errors.New("minimum counts cannot be negative")
		}
	}

	return sets, mins, nil
}

// concatClasses merges rune slices into a single slice.
func concatClasses(classes [][]rune) []rune {
	total := 0
	for _, c := range classes {
		total += len(c)
	}
	out := make([]rune, 0, total)
	for _, c := range classes {
		out = append(out, c...)
	}
	return out
}

// randomRune returns a cryptographically random rune from the provided slice.
func randomRune(set []rune) (rune, error) {
	if len(set) == 0 {
		return 0, errors.New("empty character set")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[n.Int64()], nil
}

// shuffleRunes performs an in-place Fisher–Yates shuffle using crypto/rand.
func shuffleRunes(r []rune) {
	for i := len(r) - 1; i > 0; i-- {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		r[i], r[j] = r[j], r[i]
	}
}

// generateUUIDv4 builds a RFC 4122 version 4 UUID.
func generateUUIDv4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Set version (4) and variant (10)
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] & 0x3F) | 0x80

	hex := func(x byte) string {
		const hexdigits = "0123456789abcdef"
		return string([]byte{hexdigits[x>>4], hexdigits[x&0x0F]})
	}

	var sb strings.Builder
	sb.Grow(36)
	for i, v := range b {
		sb.WriteString(hex(v))
		switch i {
		case 3, 5, 7, 9:
			sb.WriteByte('-')
		}
	}
	return sb.String(), nil
}

// DefaultOptions returns a baseline configuration using all character classes.
func DefaultOptions() Options {
	return Options{
		Format:    FormatGeneric,
		Length:    16,
		UseLower:  true,
		UseUpper:  true,
		UseNumber: true,
		UseSymbol: true,
	}
}
