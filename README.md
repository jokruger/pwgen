# pwgen

A flexible, cryptographically secure command‑line password / key / GUID generator written in Go and powered by Cobra.

## Features

- Formats: generic (default), appkey (segmented), guid (UUID v4)
- Toggle character classes: lowercase, uppercase, numbers, symbols
- Enforce minimum counts per character class (e.g. at least 2 symbols, 3 numbers)
- Segmented application key style output (e.g. XXXX-XXXX-XXXX-XXXX)
- Cryptographically secure randomness (crypto/rand)
- Shuffle after satisfying minima to avoid predictable placement
- Simple, scriptable CLI

## Installation

With Go installed (1.25+):

`go install github.com/yourusername/pwgen@latest`

Ensure your GOPATH/bin (or GOBIN) is on PATH.

## Quick Start

Generate a 16‑character password (default):
`pwgen`

Generate a 24‑character password:
`pwgen --length 24`

Generate a password without symbols:
`pwgen --symbol=false`

Require at least 2 numbers and 2 symbols:
`pwgen --min-number 2 --min-symbol 2`

Generate an app style key (4 x 4 chars):
`pwgen --format appkey`

Custom app key (5 segments of 6 = 30 chars):
`pwgen --format appkey --segments 5 --segment-length 6`

Generate a UUID v4:
`pwgen --format guid`

Lowercase + numbers only, length 20:
`pwgen --upper=false --symbol=false --length 20`

Explicit subcommand (same result as calling root):
`pwgen generate --length 32`

## Formats

- generic (default): random string of given length using enabled classes
- appkey: segmented string (segments \* segment-length) separated by dashes
- guid: UUID v4 (length fixed to 36 including dashes, ignores other flags except format)

## Flags

The table below lists flags relevant to the generic and appkey formats. (The guid format ignores all except `--format`.)

| Flag               | Type   | Default | Applies To      | Description                                   |
| ------------------ | ------ | ------- | --------------- | --------------------------------------------- |
| `--length, -l`     | int    | 16      | generic         | Total length (generic only)                   |
| `--format, -f`     | string | generic | all             | Output format: `generic`, `appkey`, or `guid` |
| `--lower`          | bool   | true    | generic, appkey | Include lowercase letters                     |
| `--upper`          | bool   | true    | generic, appkey | Include uppercase letters                     |
| `--number`         | bool   | true    | generic, appkey | Include digits                                |
| `--symbol`         | bool   | true    | generic, appkey | Include symbols                               |
| `--min-lower`      | int    | 0       | generic, appkey | Minimum lowercase letters                     |
| `--min-upper`      | int    | 0       | generic, appkey | Minimum uppercase letters                     |
| `--min-number`     | int    | 0       | generic, appkey | Minimum digits                                |
| `--min-symbol`     | int    | 0       | generic, appkey | Minimum symbols                               |
| `--segments`       | int    | 4       | appkey          | Number of segments                            |
| `--segment-length` | int    | 4       | appkey          | Characters per segment                        |

(Disable booleans with `--flag=false`.)

## Minimum Counts

You may specify any combination of:

- `--min-lower`
- `--min-upper`
- `--min-number`
- `--min-symbol`

Constraints:

- Cannot set a minimum for a disabled class
- Sum of minima must not exceed total length (generic) or derived total (appkey)
- Minima must be >= 0

Example:
`pwgen --length 20 --min-lower 4 --min-upper 4 --min-number 4 --min-symbol 2`

## App Key Details

Total characters = segments \* segment-length (before dashes).
Example:
`pwgen --format appkey --segments 3 --segment-length 5`
Produces 15 random characters output as XXXXX-XXXXX-XXXXX.

All character class rules (toggles and minima) apply to the total character pool (not per segment). Segments are assigned after generation + shuffle.

## GUID Format

`pwgen --format guid`

- Ignores length and character class flags
- Returns RFC 4122 version 4 UUID (e.g. `3d3f3c54-2f47-4e4e-a3ab-e2f9e899b8d2`)

## Exit Codes

- 0 success
- Non-zero: validation or runtime error (e.g. invalid minima, no classes enabled)

## Examples

Strong mixed password with enforced diversity:
`pwgen -l 40 --min-lower 5 --min-upper 5 --min-number 5 --min-symbol 5`

Numeric only one-time code:
`pwgen --lower=false --upper=false --symbol=false --length 10 --min-number 10`

Symbol-heavy:
`pwgen --min-symbol 6 --length 24`

## Security Notes

- Uses crypto/rand (CSPRNG)
- No reuse of math/rand
- Output is not stored or logged
- Shuffles characters after inserting minima

If you need to exclude ambiguous characters (like O/0, l/1) or generate multiple passwords at once, those can be future enhancements.

## Planned / Possible Enhancements

- `--count N` to emit multiple values
- `--exclude-similar` to drop visually ambiguous runes
- `--no-repeat` to prevent repeated characters
- Pronounceable / dictionary-based modes
- Output JSON or structured metadata

## Library Usage (Internal)

You can import the generator package:

`import "github.com/yourusername/pwgen/internal/generator"`

Construct options and call `generator.Generate`.

## License

MIT License. See LICENSE file.

## Help

Run:
`pwgen --help`
or
`pwgen generate --help`
