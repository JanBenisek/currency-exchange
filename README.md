# Currency Exchange CLI

Why?
A sim, colorful command-line tool for currency conversion with support for multiple display formats and themes.

## TODO
- add nice short video, better shorter readme

## Features

- **Multiple API providers**: Exchange Rates API (default), Fixer, Currency Layer, and Open Exchange Rates
- **4 beautiful color themes**: Ocean, Sunset, Forest, and Neon
- **Multiple display formats**: Table, Minimal, CSV, JSON, and Number
- **Smart number formatting**: Thousand separators for easy reading
- **K suffix support**: Use `700k` instead of `700000`
- **Multiple currency conversions**: Convert to multiple currencies in a single command

## Installation

### Homebrew (Recommended)

```bash
# Add the tap
brew tap janbenisek/currency-exchange

# Install
brew install cex
```

To upgrade:

```bash
brew upgrade cex
```

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

# Use Fixer API provider
export FIXER_API_KEY=your_key
cex -a fixer 100 usd to eur gbp

# Use CurrencyLayer provider
export CURRENCYLAYER_API_KEY=your_key
cex -a currencylayer 100 usd to eur

# Combine theme, format, and API provider
cex -t neon -f csv -a fixer 200 gbp to usd eur chf

# Flags work anywhere in the command
cex 100 -t forest -a exchangerates usd to eur

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

### `-a, --api` - API Provider

Choose from 4 currency exchange rate APIs:

| Provider            | Description                                    | API Key Required                         |
|---------------------|------------------------------------------------|------------------------------------------|
| `exchangerates`     | Exchange Rates API (v4 free, or v6 with key)  | Optional (`EXCHANGERATE_API_KEY` for v6) |
| `fixer`             | Fixer.io API                                    | Yes (`FIXER_API_KEY`)                     |
| `currencylayer`     | CurrencyLayer API                               | Yes (`CURRENCYLAYER_API_KEY`)            |
| `openexchangerates` | Open Exchange Rates API                         | Yes (`OPENEXCHANGERATES_API_KEY`)         |

**Default:** `exchangerates`

#### Setting API Keys

For providers that require API keys, set them as environment variables:

```bash
# Fixer.io
export FIXER_API_KEY=your_key_here
cex -a fixer 100 usd to eur

# CurrencyLayer
export CURRENCYLAYER_API_KEY=your_key_here
cex -a currencylayer 100 usd to eur

# Open Exchange Rates
export OPENEXCHANGERATES_API_KEY=your_key_here
cex -a openexchangerates 100 usd to eur
```

**Using a .env file:**

For convenience, you can use a `.env` file to store your API keys:

```bash
# Copy the example file
cp .env.example .env

# Edit .env with your actual API keys
nano .env  # or your preferred editor

# Source the file in your shell
source .env
```

⚠️ **Security Note:** The `.env` file is git-ignored for security. Never commit your API keys to version control.

#### API Provider Details

**Exchange Rates API** (`exchangerates`)
- **v4**: Free tier with no API key required (up to 1,500 requests/month)
- **v6**: Optional API key for higher limits and more features
- Reliable rates for 170+ currencies
- Updated every 60 seconds
- Website: [exchangerate-api.com](https://www.exchangerate-api.com/)
- **How to get API key (v6 - optional):**
  1. Visit [exchangerate-api.com](https://www.exchangerate-api.com/)
  2. Sign up for a free account
  3. Your API key will be available in the dashboard
  4. Free v6 tier: 10,000 requests/month (vs 1,500 for v4)
  5. Set `EXCHANGERATE_API_KEY` environment variable to use v6

**Fixer** (`fixer`)
- Requires free API key from [fixer.io](https://fixer.io/)
- Supports historical exchange rates
- 170+ currencies available
- **How to get API key:**
  1. Visit [apilayer.com/fixer](https://apilayer.com/fixer)
  2. Sign up for a free account (no credit card required)
  3. Verify your email address
  4. Your API key will be available in the dashboard
  5. Free tier: 1,000 requests/month

**Currency Layer** (`currencylayer`)
- Requires free API key from [currencylayer.com](https://currencylayer.com/)
- Live and historical forex rates
- JSON API for exchange data
- **How to get API key:**
  1. Visit [currencylayer.com/product](https://currencylayer.com/product)
  2. Sign up for a free account
  3. Verify your email address
  4. Your API access key will be in your account dashboard
  5. Free tier: 1,000 requests/month

**Open Exchange Rates** (`openexchangerates`)
- Requires free API key from [openexchangerates.org](https://openexchangerates.org/)
- Reliable, consistent exchange rate data
- Used by startups and organizations worldwide
- **How to get API key:**
  1. Visit [openexchangerates.org/signup](https://openexchangerates.org/signup)
  2. Sign up for a free account (no credit card required)
  3. Your App ID (API key) will be displayed immediately and sent to your email
  4. Free tier: 1,000 requests/month

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

- Exchange rates provided by multiple APIs:
  - [Exchange Rate API](https://www.exchangerate-api.com/) (default)
  - [Fixer.io](https://fixer.io/)
  - [CurrencyLayer](https://currencylayer.com/)
  - [Open Exchange Rates](https://openexchangerates.org/)
- Built with Go and beautiful terminal output using ANSI colors