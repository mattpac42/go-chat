# Project Navigator Agent

The Project Navigator is your project's institutional memory and knowledge companion.

## Overview

The Project Navigator builds comprehensive understanding of your project through code analysis, documentation review, and conversation learning. It becomes your source for project-specific questions, historical context, and institutional knowledge.

## When to Use

### Use For
- Project-specific questions: "Where is the authentication logic?"
- Historical context: "Why did we choose this architecture?"
- Decision rationale: "What were the trade-offs for this approach?"
- Pattern queries: "What's our established error handling pattern?"
- Business logic: "How does the payment flow work?"

### Don't Use For
- General programming → Use `developer`
- System setup → Use `garden-guide`
- Infrastructure → Use `platform`
- Security analysis → Use specialized security agents

## Key Capabilities

### Continuous Learning
The Navigator learns from:

**Git History**
- Commit messages for decision context
- Code diffs for pattern identification
- PR descriptions for rationale

**Documentation**
- README updates and project docs
- Code comments and inline documentation
- Architecture decision records

**Conversations**
- Technical discussions about decisions
- Architecture explanations
- Business logic clarifications

### Knowledge Synthesis

**Decision Tracking**
- Why technology X was chosen over Y
- Constraints at time of decision
- How decisions worked in practice

**Pattern Recognition**
- Error handling approaches
- API response formats
- Database migration patterns
- Testing conventions

**Business Logic Mapping**
- User flows and journeys
- Payment processing rules
- Validation requirements

## Question Response Patterns

### Architecture Questions
```
Q: "Where is user authentication implemented?"

Response includes:
- PRIMARY LOCATION: Specific files and line ranges
- RELATED COMPONENTS: Dependencies and connected systems
- DECISION CONTEXT: Why this architecture was chosen
- USAGE PATTERNS: How other parts interact with it
```

### Decision History Questions
```
Q: "Why did we choose PostgreSQL over MongoDB?"

Response includes:
- ORIGINAL DECISION: What was decided and when
- CONTEXT: Requirements and constraints at the time
- ALTERNATIVES: Other options considered
- TRADE-OFFS: Advantages/disadvantages weighed
- OUTCOME: How it worked in practice
```

### Pattern Questions
```
Q: "What's our error handling pattern?"

Response includes:
- ESTABLISHED PATTERN: Standard approach
- CODE EXAMPLES: Where pattern is implemented
- VARIATIONS: Context-specific adaptations
- ANTI-PATTERNS: What to avoid
```

## Integration with Development

### Context for Other Agents
When agents are invoked, Navigator provides:
- Relevant code patterns for the task
- Historical decisions affecting implementation
- Business logic context and constraints
- Performance considerations from past work

### Learning from Agent Work
When agents complete tasks, Navigator captures:
- Implementation patterns used
- Decisions made during development
- Performance considerations discovered
- Business logic insights gained

## Cross-Session Continuity

When returning to work, Navigator provides:
- Context from previous sessions
- Decisions made since last conversation
- New patterns identified in recent commits
- Updated business logic understanding

## Integration with Other Agents

| Agent | Relationship |
|-------|-------------|
| garden-guide | Guide = system setup; Navigator = project knowledge |
| developer | Navigator provides context for better implementations |
| architect | Navigator provides historical decisions for design choices |

---

The Project Navigator preserves the "why" behind every decision, helping you build on knowledge gained throughout your project's evolution.
