#!/bin/bash
set -e

echo "Setting up Go Chat development environment..."

# Install PulseAudio client for audio passthrough to host
sudo apt-get update && sudo apt-get install -y pulseaudio-utils alsa-utils

# Create afplay shim for macOS audio compatibility (Claude Code uses afplay)
sudo tee /usr/local/bin/afplay > /dev/null << 'SHIM'
#!/bin/bash
paplay "$1" 2>/dev/null || aplay "$1" 2>/dev/null
SHIM
sudo chmod +x /usr/local/bin/afplay

# Install Claude Code CLI
sudo npm install -g @anthropic-ai/claude-code

# Install project dependencies if package.json exists
if [ -f "package.json" ]; then
    npm install
fi

# Git configuration
git config --global init.defaultBranch main
git config --global pull.rebase false

# GitLab SSL certificate setup (for self-hosted gitlab.yuki.lan)
# Option A: If you have the CA cert file, copy it to .devcontainer/gitlab-ca.crt and uncomment:
# if [ -f ".devcontainer/gitlab-ca.crt" ]; then
#     sudo cp .devcontainer/gitlab-ca.crt /usr/local/share/ca-certificates/gitlab-yuki-lan.crt
#     sudo update-ca-certificates
#     echo "GitLab CA certificate installed"
# fi

# Option B: Disable SSL verification for GitLab only (less secure, but works)
git config --global http.https://gitlab.yuki.lan/.sslVerify false

echo ""
echo "Setup complete!"
echo "Run: claude auth login"
echo "Start: npm run dev"
