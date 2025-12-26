#!/bin/bash

# Script to generate Homebrew formula for alidash
# Usage: ./scripts/generate-formula.sh <version> <github-username>

set -e

VERSION=${1:-"1.0.0"}
GITHUB_USER=${2:-"your-username"}

if [ "$GITHUB_USER" = "your-username" ]; then
    echo "Usage: $0 <version> <github-username>"
    echo "Example: $0 1.0.0 your-username"
    exit 1
fi

TARBALL_URL="https://github.com/${GITHUB_USER}/alidash/archive/v${VERSION}.tar.gz"
TEMP_FILE="/tmp/alidash-${VERSION}.tar.gz"

echo "Downloading tarball to calculate SHA256..."
curl -L "$TARBALL_URL" -o "$TEMP_FILE"

if [ ! -f "$TEMP_FILE" ]; then
    echo "Error: Failed to download tarball from $TARBALL_URL"
    echo "Make sure the release v${VERSION} exists on GitHub"
    exit 1
fi

SHA256=$(shasum -a 256 "$TEMP_FILE" | cut -d' ' -f1)
echo "SHA256: $SHA256"

# Generate the formula
cat > alidash.rb << EOF
class Alidash < Formula
  desc "Terminal User Interface (TUI) application for managing Alibaba Cloud resources"
  homepage "https://github.com/${GITHUB_USER}/alidash"
  url "https://github.com/${GITHUB_USER}/alidash/archive/v${VERSION}.tar.gz"
  sha256 "${SHA256}"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd"
  end

  test do
    assert_match "alidash", shell_output("#{bin}/alidash --help 2>&1", 1)
  end
end
EOF

echo "Generated alidash.rb formula for version ${VERSION}"
echo "GitHub user: ${GITHUB_USER}"
echo "SHA256: ${SHA256}"

# Clean up
rm -f "$TEMP_FILE"

echo ""
echo "Next steps:"
echo "1. Review the generated alidash.rb file"
echo "2. Test the formula: brew install --build-from-source ./alidash.rb"
echo "3. Create a homebrew-alidash repository on GitHub"
echo "4. Copy alidash.rb to the root of that repository"
echo "5. Users can then install with: brew tap ${GITHUB_USER}/alidash && brew install alidash" 