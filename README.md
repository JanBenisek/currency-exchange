# Currency Exchange CLI

Why?
A sim, colorful command-line tool for currency conversion with support for multiple display formats and themes.

## TODO
- add nice short videa, better shorter readme, write portions yuourself
- better release process, then to homebrew!
  - see how they do it https://github.com/xxczaki/cash-cli
- support multiple APIs - there is newer version of the API we are using, also others
- another feat - add swiss exchange currency for tax purpose, the ones I am using

## Features

- Real-time exchange rates from exchangerate-api.com
- 4 beautiful color themes: Ocean, Sunset, Forest, and Neon
- Multiple display formats: Table, Minimal, CSV, and Number
- Smart number formatting: Thousand separators for easy reading
- K suffix support: Use `700k` instead of `700000`
- Multiple currency conversions in a single command

## Installation

### From Source

```bash
git clone https://github.com/janbenisek/currency-exchange.git
cd currency-exchange
go build -o cex .
sudo mv cex /usr/local/bin/  # Optional: add to PATH
```

### From Binary

Download the latest binary from the [Releases](https://github.com/janbenisek/currency-exchange/releases) page.

## Usage

### Basic Usage

```bash
cex <amount> <from_currency> to <to_currency1> [to_currency2] ...
```

### Examples

```bash
# Convert 100 USD to EUR and GBP
cex 100 usd to eur gbp

# Convert 700,000 CHF to CZK (using k suffix)
cex 700k chf to czk

# Convert to multiple currencies
cex 50 usd to gbp jpy cad aud

# Use sunset theme (with shorthand)
cex -t sunset 100 eur to usd gbp

# Use minimal format
cex -f minimal 100 usd to eur jpy

# Use CSV format
cex -f csv 100 usd to eur gbp

# Use number format (for scripts)
cex -f number 100 usd to eur gbp

# Combine theme and format
cex -t neon -f csv 200 gbp to usd eur chf

# Flags work anywhere in the command
cex 100 -t forest usd to eur

# Show version
cex -v
```

## Options

### `-v, --version` - Show Version

Display the current version of cex.

```bash
$ cex -v
cex 1.0.0
```

### `-t, --theme` - Color Theme

Choose from 4 beautiful color themes:

| Theme    | Description                    | Colors              |
|----------|--------------------------------|---------------------|
| `ocean`  | Blues and cyans - professional | Cyan, Blue         |
| `sunset` | Warm tones - yellows and oranges | Orange, Yellow    |
| `forest` | Greens and earthy tones        | Green, Lime        |
| `neon`   | Bright high-contrast colors    | Magenta, Green     |

**Default:** `ocean`

### `-f, --format` - Output Format

Choose from 4 output formats:

| Format   | Description                              | Example                          |
|----------|------------------------------------------|-----------------------------------|
| `table`  | Formatted table with headers             | See examples below                |
| `minimal`| Compact format with currency codes       | Quick reference                   |
| `csv`    | Semicolon-delimited values with header  | For data processing              |
| `json`    | `.json` formatted  | For data processing              |
| `number` | Only converted values (space-separated) | Script integration               |

**Default:** `table`

## Output Examples

### Table Format (Default)

```bash
$ cex 700k chf to czk
```

```
╭──────────┬───────────────┬─────────╮
│ CURRENCY │   CONVERTED  │   RATE  │
├──────────┼───────────────┼─────────┤
│ CZK      │ 18,046,000.00 │ 25.7800 │
╰──────────┴───────────────┴─────────╯
```

### Multiple Currencies

```bash
$ cex 100 usd to eur gbp jpy
```

```
╭──────────┬───────────┬──────────╮
│ CURRENCY │ CONVERTED │     RATE │
├──────────┼───────────┼──────────┤
│ EUR      │ 85.60     │   0.8560 │
│ GBP      │ 73.30     │   0.7330 │
│ JPY      │ 15,895.00 │ 158.9500 │
╰──────────┴───────────┴──────────╯
```

### Sunset Theme

```bash
$ cex -t sunset 100 usd to eur gbp
```

```
╭──────────┬───────────┬────────╮
│ CURRENCY │ CONVERTED │  RATE  │
├──────────┼───────────┼────────┤
│ EUR      │ 85.60     │ 0.8560 │
│ GBP      │ 73.30     │ 0.7330 │
╰──────────┴───────────┴────────╯
```

### Minimal Format

```bash
$ cex -f minimal 100 usd to eur gbp
```

```
100.00 USD

  EUR        85.60  @  0.8560
  GBP        73.30  @  0.7330
```

### CSV Format

```bash
$ cex -f csv 100 usd to eur gbp
```

```
amount;currency_from;currency_to;rate;converted;date
100.00;USD;EUR;0.8560;85.60;2026-08-23
100.00;USD;GBP;0.7330;73.30;2026-08-23
```

### Number Format

```bash
$ cex -f number 100 usd to eur gbp
```

```
85.60 73.30
```

Perfect for scripting and automation.

## Features in Detail

### K Suffix Support

Use `k` suffix for thousands to make large numbers more readable:

```bash
cex 700k chf to czk    # Same as: cex 700000 chf to czk
cex 1.5k usd to eur   # Supports decimals: 1.5k = 1,500
```

### Number Formatting

All numbers are automatically formatted with thousand separators:

- `18046000.00` → `18,046,000.00`
- `1234567.89` → `1,234,567.89`

### Color Themes

Each theme provides a unique color palette:

- **Ocean**: Professional blues and cyans, perfect for work environments
- **Sunset**: Warm oranges and yellows for a friendly feel
- **Forest**: Calming greens and earthy tones
- **Neon**: High-contrast bright colors for dark terminals

## Development

### Running Tests

```bash
go test -v
```

### Building from Source

```bash
go build -o cex .
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Release Process

To create a new release:

1. **Update version in Makefile** (if using) or ensure the build process includes version info
2. **Create a git tag**:
   ```bash
   git tag -a v1.x.x -m "Release v1.x.x"
   git push origin v1.x.x
   ```
3. **GitHub Actions will automatically**:
   - Build binaries for multiple platforms (darwin, linux, windows)
   - Create a GitHub Release with the binaries
   - Generate a checksums file

4. **Manual steps (if needed)**:
   - Update Homebrew formula (if maintaining one)
   - Update documentation with new features
   - Test the release binaries

### Development

To build locally with version info:

```bash
go build -ldflags "-X main.Version=dev" -o cex .
```

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Exchange rates provided by [exchangerate-api.com](https://www.exchangerate-api.com/)
- Built with Go and beautiful terminal output using ANSI colors