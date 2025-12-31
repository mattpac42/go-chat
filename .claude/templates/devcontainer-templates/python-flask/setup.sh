#!/bin/bash
set -e

echo "Setting up {{PROJECT_NAME}} development environment..."

# System dependencies
sudo apt-get update
sudo apt-get install -y build-essential curl ca-certificates openssl

# Add GitLab (gitlab.yuki.lan) certificate to trusted store
echo "Fetching gitlab.yuki.lan certificate..."
echo | openssl s_client -servername gitlab.yuki.lan -connect gitlab.yuki.lan:443 2>/dev/null | \
    openssl x509 -outform PEM | sudo tee /usr/local/share/ca-certificates/gitlab.yuki.lan.crt > /dev/null
sudo update-ca-certificates

# Audio support for Claude Code notifications (optional - fails gracefully)
sudo apt-get install -y -qq pulseaudio-utils alsa-utils 2>/dev/null || true

# Create afplay shim (Claude uses macOS afplay command)
sudo tee /usr/local/bin/afplay > /dev/null << 'SHIM'
#!/bin/bash
paplay "$1" 2>/dev/null || aplay "$1" 2>/dev/null || true
SHIM
sudo chmod +x /usr/local/bin/afplay

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
