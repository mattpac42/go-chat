package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the application.
type Config struct {
	// Server settings
	Port string `envconfig:"PORT" default:"8080"`

	// Database settings
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`

	// Claude API settings
	ClaudeAPIKey    string `envconfig:"CLAUDE_API_KEY" required:"true"`
	ClaudeModel     string `envconfig:"CLAUDE_MODEL" default:"claude-sonnet-4-20250514"`
	ClaudeMaxTokens int    `envconfig:"CLAUDE_MAX_TOKENS" default:"4096"`

	// Context settings
	ContextMessageLimit int `envconfig:"CONTEXT_MESSAGE_LIMIT" default:"20"`

	// Logging settings
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// CORS settings
	CORSOrigins string `envconfig:"CORS_ORIGINS" default:"*"`

	// WebSocket settings
	WSPingInterval time.Duration `envconfig:"WS_PING_INTERVAL" default:"30s"`
	WSPongTimeout  time.Duration `envconfig:"WS_PONG_TIMEOUT" default:"60s"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MustLoad reads configuration and panics on error.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
