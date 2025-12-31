#!/bin/bash
set -e

echo "Setting up {{PROJECT_NAME}} development environment..."

# Install system dependencies and CA certificates
sudo apt-get update -qq
sudo apt-get install -y -qq ca-certificates openssl

# Add GitLab (gitlab.yuki.lan) certificate to trusted store
echo "Fetching gitlab.yuki.lan certificate..."
echo | openssl s_client -servername gitlab.yuki.lan -connect gitlab.yuki.lan:443 2>/dev/null | \
    openssl x509 -outform PEM | sudo tee /usr/local/share/ca-certificates/gitlab.yuki.lan.crt > /dev/null
sudo update-ca-certificates

# Install Claude Code CLI
sudo npm install -g @anthropic-ai/claude-code

# Audio support for Claude Code notifications (optional - fails gracefully)
sudo apt-get install -y -qq pulseaudio-utils alsa-utils 2>/dev/null || true

# Create afplay shim (Claude uses macOS afplay command)
sudo tee /usr/local/bin/afplay > /dev/null << 'SHIM'
#!/bin/bash
paplay "$1" 2>/dev/null || aplay "$1" 2>/dev/null || true
SHIM
sudo chmod +x /usr/local/bin/afplay

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
