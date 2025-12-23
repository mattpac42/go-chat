#!/bin/bash
set -e

echo "Setting up {{PROJECT_NAME}} development environment..."

# System dependencies
sudo apt-get update
sudo apt-get install -y build-essential curl ca-certificates

# Install Node.js (required for Claude Code CLI)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install Claude Code CLI
sudo npm install -g @anthropic-ai/claude-code

# Python dependencies
pip install --upgrade pip
if [ -f "requirements.txt" ]; then
    pip install -r requirements.txt
else
    pip install flask flask-sqlalchemy python-dotenv
fi

# Development tools
pip install black flake8 pytest

# Git configuration
git config --global init.defaultBranch main
git config --global pull.rebase false

echo ""
echo "Setup complete!"
echo "Run: claude auth login"
echo "Start: flask run --host=0.0.0.0"
