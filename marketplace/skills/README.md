# Marketplace Skills

Community and extended skills for specialized workflows.

## Core vs Marketplace Skills

| Location | Purpose | Copied to Projects |
|----------|---------|-------------------|
| `.claude/skills/` | Core workflows (PRD, tasks, handoff) | Always |
| `marketplace/skills/` | Extended/community skills | Yes, with marketplace |

## What Are Marketplace Skills?

Extended skills that provide:
- Domain-specific workflows
- Community-contributed automation
- Specialized integrations
- Advanced tooling

## Skill Structure

```
skills/
├── README.md           # This file
└── <skill-name>/
    ├── SKILL.md        # Skill definition and instructions
    ├── README.md       # Usage documentation (optional)
    └── scripts/        # Supporting scripts (optional)
```

## Skill Definition (SKILL.md)

```markdown
# Skill Name

Brief description of what this skill does.

## Triggers

- "trigger phrase one"
- "trigger phrase two"

## Workflow

1. Step one
2. Step two
3. Step three

## Outputs

- What this skill produces
```

## Using Marketplace Skills

Skills in this directory are available to all planted projects. They can be:
- Invoked directly by trigger phrases
- Referenced by agents
- Chained with core skills

## Contributing Skills

1. Create a folder with your skill name (kebab-case)
2. Add `SKILL.md` with skill definition
3. Add `README.md` if additional documentation needed
4. Include any supporting scripts in `scripts/`
5. Submit to The Garden

## Available Skills

*No marketplace skills installed yet. This directory is ready for community contributions.*

---

**Status:** Ready for contributions
**Version:** Garden 2.0
