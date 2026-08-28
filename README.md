# `cex` - Currency Exchange CLI

A colorful command-line currency converter. Type an amount, get rates — no browser needed. For everyone who wants to keep their money data locally.

```bash
$ cex 100 usd to eur gbp jpy
╭──────────┬───────────┬──────────╮
│ CURRENCY │ CONVERTED │     RATE │
├──────────┼───────────┼──────────┤
│ EUR      │ 85.60     │   0.8560 │
│ GBP      │ 73.30     │   0.7330 │
│ JPY      │ 15,895.00 │ 158.9500 │
╰──────────┴───────────┴──────────╯
```

## Install

```bash
brew tap janbenisek/currency-exchange https://github.com/janbenisek/currency-exchange
brew trust janbenisek/currency-exchange
brew install cex
```

Or grab a binary from [Releases](https://github.com/janbenisek/currency-exchange/releases), or build from source:

```bash
git clone https://github.com/janbenisek/currency-exchange.git
cd currency-exchange
go build -o cex .
```

## Examples

```bash
cex 100 usd to eur          # 100 USD -> EUR
cex 700k chf to czk         # 700,000 CHF -> CZK (k = thousands)
cex 50 usd to gbp jpy cad   # convert to several currencies at once
cex -f minimal 100 usd to jpy
cex -t sunset 100 eur to usd
```

## Options

| Flag | What it does | Choices |
|------|--------------|---------|
| `-t` | Color theme | `ocean` (default), `sunset`, `forest`, `neon` |
| `-f` | Output format | `table` (default), `minimal`, `csv`, `json`, `number` |
| `-a` | Rate provider | `exchangerates` (default), `fixer`, `currencylayer`, `openexchangerates` |
| `-v` | Show version | |

The default provider works with no API key. To use another one, generate a free key and export it:

```bash
export FIXER_API_KEY=your_key
export CURRENCYLAYER_API_KEY=your_key
export OPENEXCHANGERATES_API_KEY=your_key
```

## Development

```bash
go test ./...
go build -o cex .
```

## License

MIT — see [LICENSE](LICENSE).
