#!/bin/bash
set -e

echo "Setting up Go Chat development environment..."

# Install Claude Code CLI
sudo npm install -g @anthropic-ai/claude-code

# Install project dependencies if package.json exists
if [ -f "package.json" ]; then
    npm install
fi

# Git configuration
git config --global init.defaultBranch main
git config --global pull.rebase false

echo ""
echo "Setup complete!"
echo "Run: claude auth login"
echo "Start: npm run dev"
