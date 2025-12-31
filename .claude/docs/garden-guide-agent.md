# Garden Guide Agent

The Garden Guide is your coach for setting up and optimizing the Claude Agent System.

## Overview

The Garden Guide agent helps with system setup, workflow optimization, and best practices coaching. It serves as a step-by-step setup wizard and ongoing optimization advisor.

## When to Use

### Use For
- New project setup and initialization
- System optimization and configuration
- Workflow guidance and template selection
- Best practices and troubleshooting

### Don't Use For
- Project-specific questions → Use `project-navigator`
- Code implementation → Use `developer`
- Infrastructure setup → Use `platform`

## Key Capabilities

### Proactive Coaching
The Garden Guide anticipates needs and suggests improvements:
- Observes usage patterns and recommends optimizations
- Suggests PROJECT.md enhancements
- Recommends appropriate workflow templates

### Setup Workflows

**Phase 1: Assessment**
- Project type and tech stack
- Team size and experience
- Development challenges

**Phase 2: Foundation**
- Verify .claude/ directory structure
- Create PROJECT.md with project context
- Copy required templates

**Phase 3: Optimization**
- Enhance PROJECT.md for agent effectiveness
- Configure workflow templates
- Establish productive patterns

## Setup Checklist

### Directory Structure
```
.claude/
├── PROJECT.md           # Your project context
├── PROTOCOLS.md         # System rules
├── QUICKSTART.md        # Setup guide
├── agents/              # Generic agents
├── skills/              # Workflow skills
├── commands/            # User commands
├── templates/           # Document templates
├── work/
│   ├── 0_vision/        # Vision documents
│   ├── 1_backlog/       # PRDs
│   ├── 2_active/        # In progress
│   ├── 3_done/          # Completed
│   └── history/         # Session logs
└── config/              # Configuration
```

### PROJECT.md Template
```markdown
# Project Context

## Overview
[Your project description]

## Tech Stack
[Languages, frameworks, tools]

## Key Files
[Important files to know about]

## Conventions
[Coding standards, patterns to follow]
```

## Troubleshooting

### Common Issues

**Problem**: Agents don't understand project context
- **Solution**: Enhance PROJECT.md with more detail
- Add domain model, business rules, conventions

**Problem**: Tasks consume too much context
- **Solution**: Break into smaller, focused tasks
- Use modular task breakdown

**Problem**: Agent selection confusion
- **Solution**: Review CLAUDE.md decision tree
- Generic for quick tasks, specialized for production

## Integration with Other Agents

| Agent | Relationship |
|-------|-------------|
| project-navigator | Garden Guide = system setup; Navigator = project knowledge |
| developer | Garden Guide provides context optimization for better agent output |
| All agents | Garden Guide ensures clear task scope and success criteria |

## Success Indicators

### Setup Complete
- .claude/ directory structure verified
- PROJECT.md comprehensive and project-specific
- First task completed successfully
- User understands agent selection

### Ongoing Optimization
- Context usage per task optimized
- Workflow templates customized
- Quality of agent output meets expectations

---

The Garden Guide is your partner in mastering the Claude Agent System.
