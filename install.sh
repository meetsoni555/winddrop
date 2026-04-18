#!/bin/bash

set -e

echo "🌬️ Installing WindDrop..."

# Step 1: Build binary
echo "🔨 Building binary..."
go build -o winddrop

# Step 2: Make executable
chmod +x winddrop

# Step 3: Install to /usr/local/bin
echo "📦 Installing to /usr/local/bin..."
sudo cp winddrop /usr/local/bin/

# Step 4: Verify install
if command -v winddrop &> /dev/null
then
    echo "✅ WindDrop installed successfully!"
    echo ""
    echo "👉 Try:"
    echo "   winddrop send <file>"
    echo "   winddrop send <file> --expire 5m"
    echo "   winddrop send <file> --once expire 2m"
else
    echo "❌ Installation failed. Check PATH."
fi
