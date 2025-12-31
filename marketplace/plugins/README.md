# Marketplace Plugins

Extension plugins for Garden functionality.

## What Are Plugins?

Plugins extend Garden's core capabilities with:
- Custom tool integrations
- Third-party service connectors
- Specialized workflows
- Domain-specific automation

## Plugin Structure

```
plugins/
├── README.md           # This file
└── <plugin-name>/
    ├── plugin.json     # Plugin metadata and configuration
    ├── README.md       # Plugin documentation
    └── ...             # Plugin-specific files
```

## Plugin Metadata (plugin.json)

```json
{
  "name": "example-plugin",
  "version": "1.0.0",
  "description": "What this plugin does",
  "author": "Your Name",
  "requires": {
    "garden": ">=2.0.0"
  },
  "provides": {
    "tools": [],
    "skills": [],
    "agents": []
  }
}
```

## Using Plugins

Plugins are included when projects are planted from The Garden. To use a plugin:

1. Ensure the plugin exists in `marketplace/plugins/`
2. Reference it in your project's `.claude/config/plugins.json`
3. Follow plugin-specific setup instructions

## Contributing Plugins

1. Create a folder with your plugin name
2. Add `plugin.json` with metadata
3. Add `README.md` with documentation
4. Include any required files (scripts, templates, etc.)
5. Submit to The Garden

## Available Plugins

*No plugins installed yet. This directory is ready for future extensions.*

---

**Status:** Ready for contributions
**Version:** Garden 2.0
