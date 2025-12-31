# The Garden Marketplace

Central repository for agents, plugins, skills, and knowledge resources.

## Structure

```
marketplace/
├── agents/           # 130+ specialized agents
├── plugins/          # Extension plugins (future)
├── skills/           # Community skills (future)
└── knowledge-base/   # Reference materials
```

## Directories

### agents/
130+ specialized agents organized by domain. See [agents/README.md](agents/README.md) for full inventory.

**Key domains:**
- Software Development (tactical/strategic pairs)
- Platform & Infrastructure
- Product Management
- UX/UI Design
- Security & Compliance
- Data & Analytics
- Medical Specialists
- Game Development
- Legal & Finance

**Usage pattern:** Specialized agents use tactical/strategic pairs:
- `software-tactical` + `software-strategic` (code + architecture)
- `platform-tactical` + `platform-strategic` (infra + cloud strategy)
- `product-manager-tactical` + `product-manager-strategic` (sprints + roadmap)

### plugins/
Extension plugins for Garden functionality. *Coming soon.*

### skills/
Community-contributed skills. *Coming soon.*

### knowledge-base/
Reference materials and documentation resources.

## Using Marketplace Resources

### From The Garden (baseline)
Resources are available directly:
```
marketplace/agents/software-tactical.md
marketplace/knowledge-base/...
```

### From Planted Projects
Use `/sync-baseline` to pull updates from The Garden, including marketplace resources.

## Contributing

### Adding Agents
1. Create agent file in `agents/` following naming conventions
2. Use domain prefix (e.g., `software-`, `platform-`, `medical-`)
3. Include tactical/strategic designation if applicable
4. Update `agents/README.md` with new agent

### Adding Plugins
*Guidelines coming soon*

### Adding Skills
*Guidelines coming soon*

---

**Version:** Garden 2.0
**Location:** `marketplace/`
