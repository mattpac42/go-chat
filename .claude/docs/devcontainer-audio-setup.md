# Devcontainer Audio Setup

Enable Claude Code audio notifications in devcontainers. This is optional - containers work without audio.

## How It Works

1. **Container side** (automatic): Devcontainer templates include PulseAudio client tools and an `afplay` shim that translates macOS audio commands to Linux
2. **Host side** (manual): You run a PulseAudio server that accepts connections from containers

## Host Setup

### macOS with Podman (Default)

Templates default to Podman's hostname. Just start PulseAudio:

```bash
# Install PulseAudio
brew install pulseaudio

# Start with TCP and anonymous auth
pulseaudio --load="module-native-protocol-tcp auth-anonymous=1" --exit-idle-time=-1 --daemon

# Verify it's running
pactl info
```

### macOS with Docker Desktop

Docker uses a different hostname for host access. Update your devcontainer.json:

```json
"containerEnv": {
  "PULSE_SERVER": "host.docker.internal"
}
```

Then start PulseAudio (same commands as Podman above).

### Linux Host

```bash
# Usually PulseAudio is already running
# Enable TCP module
pactl load-module module-native-protocol-tcp auth-anonymous=1

# Or add to /etc/pulse/default.pa for persistence:
# load-module module-native-protocol-tcp auth-anonymous=1
```

## Auto-Start PulseAudio (macOS)

Create a LaunchAgent for automatic startup:

```bash
mkdir -p ~/Library/LaunchAgents

cat > ~/Library/LaunchAgents/org.pulseaudio.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>org.pulseaudio</string>
    <key>ProgramArguments</key>
    <array>
        <string>/opt/homebrew/bin/pulseaudio</string>
        <string>--load=module-native-protocol-tcp auth-anonymous=1</string>
        <string>--exit-idle-time=-1</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

# Load it
launchctl load ~/Library/LaunchAgents/org.pulseaudio.plist
```

## Troubleshooting

### No sound in container

```bash
# Check PULSE_SERVER is set
echo $PULSE_SERVER

# Test connection
paplay /usr/share/sounds/alsa/Front_Center.wav

# Check if host is reachable
nc -zv $PULSE_SERVER 4713
```

### PulseAudio not accepting connections

```bash
# On host, check if TCP module is loaded
pactl list modules | grep tcp

# If not, load it
pactl load-module module-native-protocol-tcp auth-anonymous=1
```

### Wrong hostname (Podman vs Docker)

| Runtime | Hostname |
|---------|----------|
| Podman (default) | `host.containers.internal` |
| Docker | `host.docker.internal` |

Docker users: update `PULSE_SERVER` in devcontainer.json to `host.docker.internal`.

## Security Note

`auth-anonymous=1` allows any connection to your PulseAudio server. This is fine for local development but should not be used in production or shared environments.

For more secure setups, use PulseAudio cookie authentication or SSH tunneling.
