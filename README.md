# Currency Exchange CLI

A beautiful, colorful command-line tool for currency conversion with support for multiple display formats and themes.

## Features

- Real-time exchange rates from exchangerate-api.com
- 4 beautiful color themes: Ocean, Sunset, Forest, and Neon
- Multiple display formats: Table, Cards, List, and Minimal
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

### Homebrew (macOS/Linux)

```bash
brew tap janbenisek/currency-exchange
brew install cex
```

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

# Use sunset theme
cex -theme sunset 100 eur to usd gbp

# Use list format
cex -format list 100 usd to eur jpy

# Combine theme and format
cex -theme neon -format cards 200 gbp to usd eur chf
```

## Options

### `-theme` - Color Theme

Choose from 4 beautiful color themes:

| Theme    | Description                    | Colors              |
|----------|--------------------------------|---------------------|
| `ocean`  | Blues and cyans - professional | Cyan, Blue         |
| `sunset` | Warm tones - yellows and oranges | Orange, Yellow    |
| `forest` | Greens and earthy tones        | Green, Lime        |
| `neon`   | Bright high-contrast colors    | Magenta, Green     |

**Default:** `ocean`

### `-format` - Display Format

Choose from 4 display formats:

| Format   | Description                              | Example                          |
|----------|------------------------------------------|-----------------------------------|
| `table`  | Clean two-column table with borders     | See examples below                |
| `cards`  | Individual cards for each currency        | Good for single conversions      |
| `list`   | Compact list format                      | Great for multiple currencies     |
| `minimal`| Ultra-compact single line                | Quick reference                   |

**Default:** `table`

## Output Examples

### Table Format (Default)

```bash
$ cex 700k chf to czk
```

```
Currency Exchange converting 700000.00 CHF

┌─────────────────────┬────────────────────┐
│ Converted Value     │ Rate (1 CHF = X)   │
├─────────────────────┼────────────────────┤
│ 18,046,000.00 CZK  │ 1 CHF = 25.7800 CZK│
└─────────────────────┴────────────────────┘
```

### Multiple Currencies

```bash
$ cex 100 usd to eur gbp jpy
```

```
Currency Exchange converting 100.00 USD

┌─────────────────────┬────────────────────┐
│ Converted Value     │ Rate (1 USD = X)   │
├─────────────────────┼────────────────────┤
│ 85.60 EUR          │ 1 USD = 0.8560 EUR │
│ 73.30 GBP          │ 1 USD = 0.7330 GBP │
│ 15,895.00 JPY      │ 1 USD = 158.9500 JPY│
└─────────────────────┴────────────────────┘
```

### Sunset Theme

```bash
$ cex -theme sunset 100 usd to eur gbp
```

```
Currency Exchange converting 100.00 USD

┌─────────────────────┬────────────────────┐
│ Converted Value     │ Rate (1 USD = X)   │
├─────────────────────┼────────────────────┤
│ 85.60 EUR          │ 1 USD = 0.8560 EUR │
│ 73.30 GBP          │ 1 USD = 0.7330 GBP │
└─────────────────────┴────────────────────┘
```

### Cards Format

```bash
$ cex -format cards 100 usd to eur
```

```
Currency Exchange converting 100.00 USD

┌────────────────────────────────┐
│ EUR                            │
│ Converted: 85.60               │
│ Rate: 1 USD = 0.8560 EUR       │
└────────────────────────────────┘
```

### List Format

```bash
$ cex -format list 100 usd to eur gbp
```

```
Currency Exchange converting 100.00 USD

  EUR    → 85.60            (1 USD = 0.8560 EUR)
  GBP    → 73.30            (1 USD = 0.7330 GBP)
```

### Minimal Format

```bash
$ cex -format minimal 100 usd to eur gbp
```

```
100.00 USD →
  EUR=85.60  |  GBP=73.30
```

## Features in Detail

### K Suffix Support

Use `k` suffix for thousands to make large numbers more readable:

```bash
cex 700k chf to czk    # Same as: cex 700000 chf to czk
cex 1.5m usd to eur   # Supports decimals: 1.5m = 1,500,000
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

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Exchange rates provided by [exchangerate-api.com](https://www.exchangerate-api.com/)
- Built with Go and beautiful terminal output using ANSI colors
