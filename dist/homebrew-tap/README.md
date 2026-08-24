# Homebrew Tap

This is the Homebrew tap for [currency-exchange](https://github.com/janbenisek/currency-exchange), a simple CLI tool to convert currencies.

## Installation

```bash
# Add the tap
brew tap janbenisek/tap

# Install currency-exchange
brew install currency-exchange
```

## Updating

```bash
# Upgrade to the latest version
brew upgrade currency-exchange
```

## Uninstalling

```bash
# Remove the formula
brew uninstall currency-exchange

# Remove the tap (optional)
brew untap janbenisek/tap
```

## Usage

After installation, you can use the `currency-exchange` command:

```bash
# Convert 100 USD to EUR
currency-exchange 100 usd to eur

# Convert to multiple currencies
currency-exchange 100 usd to eur gbp jpy

# Use a different theme
currency-exchange -t sunset 100 usd to eur

# Use minimal format
currency-exchange -f minimal 100 usd to eur
```

For more usage examples and options, visit the [currency-exchange repository](https://github.com/janbenisek/currency-exchange).

## Formula

The formula file is automatically updated by GitHub Actions when new releases are published.

## License

MIT - see [currency-exchange license](https://github.com/janbenisek/currency-exchange/blob/main/LICENSE) for details.
