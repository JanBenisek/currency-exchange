# This file is automatically updated by GitHub Actions on release
# Manual edits will be overwritten

class CurrencyExchange < Formula
  desc "Simple CLI tool to convert currencies"
  homepage "https://github.com/janbenisek/currency-exchange"
  version "1.0.0"

  on_macos do
    on_arm do
      url "https://github.com/janbenisek/currency-exchange/releases/download/v1.0.0/cex-darwin-arm64"
      sha256 "PLACEHOLDER_ARM64"
    end
    on_intel do
      url "https://github.com/janbenisek/currency-exchange/releases/download/v1.0.0/cex-darwin-amd64"
      sha256 "PLACEHOLDER_AMD64"
    end
  end

  on_linux do
    url "https://github.com/janbenisek/currency-exchange/releases/download/v1.0.0/cex-linux-amd64"
    sha256 "PLACEHOLDER_LINUX"
  end

  def install
    if Hardware::CPU.arm?
      bin.install "cex-darwin-arm64" => "currency-exchange"
    elsif Hardware::CPU.intel?
      bin.install "cex-darwin-amd64" => "currency-exchange"
    elsif OS.linux?
      bin.install "cex-linux-amd64" => "currency-exchange"
    end
  end

  test do
    system bin/"currency-exchange", "1", "usd", "to", "eur"
  end
end
