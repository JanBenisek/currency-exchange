# currency-exchange

Simple CLI tool to convert currencies.

## Install

### Homebrew (recommended)
```bash
brew tap JanBenisek/cex
brew install cex
```

### From source
```bash
cd /Users/janbenisek/github/currency-exchange
go build -o cex main.go
mv cex /usr/local/bin/
```

### From releases (future)
```bash
curl -L https://github.com/janbenisek/currency-exchange/releases/latest/download/cex-darwin-amd64 -o cex
chmod +x cex
mv cex /usr/local/bin/
```

## Usage

```bash
cex 100 czk to pln eur chf
```

## Example

```
$ cex 100 czk to pln eur chf

┌────────────┬─────────────┬────────────────────┐
│ Currency   │   Converted │ Rate               │
├────────────┼─────────────┼────────────────────┤
│ PLN        │       17.90 │ 1 CZK = 0.1790 PLN │
├────────────┼─────────────┼────────────────────┤
│ EUR        │        4.15 │ 1 CZK = 0.0415 EUR │
├────────────┼─────────────┼────────────────────┤
│ CHF        │        3.88 │ 1 CZK = 0.0388 CHF │
└────────────┴─────────────┴────────────────────┘
```

## Development

```bash
go run main.go 100 usd to eur gbp
```
