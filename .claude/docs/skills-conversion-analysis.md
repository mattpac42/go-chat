# Skills Conversion Analysis

**Created**: 2025-12-20
**Version**: 1.0.0
**Status**: Research & Analysis

---

## Executive Summary

This document analyzes the conversion of Garden system commands and agent workflows into Claude Code skills. After researching the skills system architecture, we provide specific conversion plans for 11 workflows, assess trade-offs, and recommend an implementation roadmap.

**Key Finding**: Skills are best suited for **autonomous, repeatable workflows** that Claude should invoke automatically. Complex interactive workflows (like agent delegation) should remain as commands or be reimplemented as hybrid approaches.

---

## Table of Contents

1. [Skills System Overview](#1-skills-system-overview)
2. [Comparison: Skills vs Commands vs Agents](#2-comparison-skills-vs-commands-vs-agents)
3. [Conversion Plans for 11 Workflows](#3-conversion-plans-for-11-workflows)
4. [Architecture Implications](#4-architecture-implications)
5. [Implementation Roadmap](#5-implementation-roadmap)
6. [Trade-offs and Recommendations](#6-trade-offs-and-recommendations)
7. [Sources](#sources)

---

## 1. Skills System Overview

### What Are Skills?

**Agent Skills** are modular capabilities that extend Claude's functionality through organized folders containing instructions, templates, and scripts. They are part of the open [Agent Skills specification](https://agentskills.io) maintained by Anthropic.

### Core Characteristics

| Aspect | Details |
|--------|---------|
| **Invocation** | Model-invoked (autonomous) - Claude decides when to use based on description |
| **Discovery** | Auto-discovered from `~/.claude/skills/` (personal) and `.claude/skills/` (project) |
| **Structure** | Directory with `SKILL.md` + optional supporting files (scripts, templates, docs) |
| **Portability** | Cross-platform - same format works across AI tools adopting agentskills.io standard |
| **Scope** | Best for focused, repeatable workflows; not for complex multi-agent orchestration |

### File Structure

```
.claude/skills/my-skill/
├── SKILL.md              # Required: Instructions + frontmatter
├── reference.md          # Optional: Documentation
├── examples.md           # Optional: Usage examples
├── scripts/              # Optional: Helper scripts
│   └── helper.py
└── templates/            # Optional: Templates
    └── template.txt
```

### SKILL.md Format

```yaml
---
name: skill-identifier
description: What this skill does and when to use it. Include specific
  triggers like file types, operations, keywords. Be detailed.
allowed-tools: Read, Grep, Glob  # Optional: Restrict tool usage
---

# Skill Title

## Instructions
Step-by-step guidance for Claude to follow

## Examples
Concrete usage examples

## Guidelines
Best practices and considerations
```

### Critical Frontmatter Fields

- **`name`**: Lowercase, hyphens, max 64 chars (e.g., `session-handoff`)
- **`description`**: Max 1024 chars - **CRITICAL** for skill discovery. Must include:
  - What the skill does
  - When to use it (specific triggers/keywords)
  - File types, operations, domain terms
- **`allowed-tools`**: Optional - restricts which tools Claude can use

### Key Limitations

1. **No interactive prompts** - Skills execute autonomously without user input during execution
2. **Description-based discovery** - Claude uses description to decide when to invoke; vague descriptions = poor discovery
3. **Single focus** - Skills should be narrow and focused; complex workflows should be broken into multiple skills
4. **No state persistence** - Skills don't maintain state between invocations
5. **Tool restrictions** - `allowed-tools` only works in Claude Code, not claude.ai

---

## 2. Comparison: Skills vs Commands vs Agents

### Quick Reference Matrix

| Aspect | Skills | Commands | Agents |
|--------|--------|----------|--------|
| **Invocation** | Model-invoked (automatic) | User-invoked (`/cmd`) | Task-tool-invoked (delegated) |
| **Discovery** | Description-based matching | Memorized by users | Selected by main agent |
| **Interactivity** | No - autonomous execution | Yes - can prompt user | Yes - multi-turn conversations |
| **Context** | Shares main agent context | Shares main agent context | Separate 200k context window |
| **Complexity** | Simple to moderate workflows | Simple utilities/shortcuts | Complex, specialized work |
| **State** | Stateless (no persistence) | Stateless | Session history files |
| **Best For** | Repeatable, autonomous tasks | Quick utilities | Deep expertise domains |
| **Examples** | Commit message generation | `/catch-up`, `/handoff` | Code implementation, security audits |

### Detailed Comparison

#### Skills

**Strengths**:
- Automatic invocation - no user memorization needed
- Cross-platform compatibility (agentskills.io standard)
- Clean separation of concerns
- Can restrict tools with `allowed-tools`
- Good for focused, repeatable workflows

**Weaknesses**:
- No interactive prompts during execution
- Discovery depends on description quality
- Can't maintain state between invocations
- Limited to simple workflows
- Shares main agent context (no isolation)

**Best Use Cases**:
- Automated documentation generation
- Commit message formatting
- File processing workflows
- Template-based document creation
- Read-only analysis tasks

---

#### Commands (Current `.claude/commands/`)

**Strengths**:
- Explicit user control (`/command`)
- Can prompt user for input
- Familiar pattern for users
- Direct execution path
- Good for utilities

**Weaknesses**:
- Users must remember commands
- Not discoverable by Claude
- Manual invocation required
- Shares main agent context
- No cross-platform portability

**Best Use Cases**:
- Session management (`/handoff`, `/catch-up`)
- Onboarding workflows (`/onboard`)
- System utilities (`/sync-baseline`)
- User-initiated workflows

---

#### Agents (Current `.claude/agents/`)

**Strengths**:
- Separate 200k context window
- Deep domain expertise
- Multi-turn interactive conversations
- Session history persistence
- Main agent orchestrates delegation

**Weaknesses**:
- Requires main agent delegation
- Higher overhead (agent briefing)
- Not directly user-invokable
- Garden-specific implementation
- More complex setup

**Best Use Cases**:
- Code implementation (tactical-software-engineer)
- Security audits (tactical-cybersecurity)
- Product planning (tactical-product-manager)
- System design (strategic-software-engineer)
- Complex multi-step workflows

---

### Decision Framework: Which to Use?

```
┌─────────────────────────────────────────────────────────────────┐
│                    DECISION TREE                                 │
└─────────────────────────────────────────────────────────────────┘

Is it a complex, multi-turn workflow requiring deep expertise?
├─ YES → Use AGENT (separate context, specialized knowledge)
└─ NO → Continue...

Does it require user interaction or prompts during execution?
├─ YES → Use COMMAND (user control, interactive)
└─ NO → Continue...

Should Claude invoke it automatically based on context?
├─ YES → Use SKILL (autonomous, description-based)
└─ NO → Use COMMAND (explicit user trigger)

Is it a simple, focused, repeatable task?
├─ YES → Use SKILL (autonomous, clean separation)
└─ NO → Use COMMAND or AGENT based on complexity
```

---

## 3. Conversion Plans for 11 Workflows

### 3.1 Handoff (Session Handoff Creation)

**Current Implementation**: `/handoff` command
- Creates `[NUMBER]-SESSION.md` (session summary)
- Creates/updates `HANDOFF-SESSION.md` (forward-looking handoff)
- Uses templates: `session-summary-template.md`, `handoff-session-template.md`
- Triggered manually at 75% context usage

**Conversion Assessment**: ⚠️ **Partial Conversion**

**Recommendation**: **Hybrid Approach**

**Reasoning**:
- Core handoff logic is deterministic and repeatable → Good fit for skill
- But requires context awareness and timing judgment → Better as command
- User may want to review before creating handoff → Command provides control

**Hybrid Approach**:

1. **Keep `/handoff` command** for user-initiated handoffs
2. **Create `session-handoff` skill** for automatic context-threshold handoffs

**Skill Implementation**:

```yaml
---
name: session-handoff
description: Automatically create session handoff documentation when context
  usage reaches 75% threshold. Use when context tracker triggers handoff
  or when user mentions session continuity, context limits, or starting
  new session. Creates numbered session summary and HANDOFF-SESSION.md files.
allowed-tools: Read, Write, Bash
---

# Session Handoff Skill

## Purpose
Automatically create session handoff documentation at 75% context threshold.

## Instructions

1. **Determine Session Number**
   - Check `.claude/context/session-history/` for existing files
   - Find highest number (e.g., 011-SESSION.md → next is 012)
   - Use zero-padded 3-digit format

2. **Create Session Summary**
   - Read template: `.claude/templates/session-summary-template.md`
   - Create: `.claude/context/session-history/[NUMBER]-SESSION.md`
   - Fill all required sections with actual session data

3. **Create/Update Handoff File**
   - Read template: `.claude/templates/handoff-session-template.md`
   - Overwrite: `.claude/context/session-history/HANDOFF-SESSION.md`
   - Include forward-looking next steps and context

4. **Confirm Creation**
   - Tell user: "Created `.claude/context/session-history/[NUMBER]-SESSION.md`"
   - Tell user: "Updated `.claude/context/session-history/HANDOFF-SESSION.md`"

## Triggers
- Context tracker outputs handoff message
- User mentions "create handoff", "session continuity", "new session"
- Context usage at or above 75%

## Quality Validation
- Ensure DELEGATION MANDATE section included
- Verify template sections filled (not placeholders)
- Check next steps are prioritized and actionable
- Validate git status and file paths accurate
```

**Migration Path**:
1. Create `.claude/skills/session-handoff/SKILL.md`
2. Copy templates to skill directory for reference
3. Keep `/handoff` command for manual triggering
4. Update context tracker to mention "session handoff" keyword

**Trade-offs**:
- **Pro**: Automatic invocation when threshold hit
- **Pro**: Reduces user burden
- **Con**: Less user control over timing
- **Con**: Shares main agent context (adds to context usage)

---

### 3.2 Catch-Up (Session Restoration)

**Current Implementation**: `/catch-up` command
- Reads PROJECT_CONTEXT.md, CLAUDE.md, HANDOFF-SESSION.md
- Reviews session history files
- Checks active tasks
- Reviews git status
- Synthesizes comprehensive summary

**Conversion Assessment**: ✅ **Good Fit for Skill**

**Recommendation**: **Convert to Skill**

**Reasoning**:
- Deterministic workflow (always same steps)
- No user interaction needed during execution
- Should trigger automatically when Claude needs project context
- Read-only operation (safe for autonomous execution)

**Skill Implementation**:

```yaml
---
name: catch-up
description: Get fully up to speed with project state, context, recent work,
  and next steps. Use when starting new session, returning to project after
  break, or when user asks to "catch up", "get context", "what's the status",
  or "where were we". Reads PROJECT_CONTEXT.md, CLAUDE.md, session history,
  active tasks, and git status. Provides comprehensive project summary.
allowed-tools: Read, Bash, Grep, Glob
---

# Project Catch-Up Skill

## Purpose
Restore full project context by reading core documentation and current state.

## Instructions

### 1. Load Core Documentation
- Read PROJECT_CONTEXT.md (project purpose, tech stack, constraints)
- Read CLAUDE.md (agent rules, workflows, protocols)
- Read INITIAL_CONTEXT.md (if exists)

### 2. Review Session History
- Read most recent `.claude/context/session-history/[NUMBER]-SESSION.md`
- Scan 2-3 previous session files for historical context
- **CRITICAL**: Read `.claude/context/session-history/HANDOFF-SESSION.md` for current state

### 3. Check Active Tasks
- List `.claude/tasks/2_active/` task files
- Read each task file (objectives, progress, blockers)
- Cross-reference with HANDOFF-SESSION.md next steps

### 4. Review PRD Backlog (if applicable)
- Check `.claude/tasks/1_backlog/` and `.claude/tasks/2_active/` for PRDs
- Read project-status.md if exists

### 5. Check Git Status
- Run `git status` for uncommitted changes
- Run `git log -5 --oneline` for recent commits

### 6. Present Summary
Provide structured output:

```markdown
# 🎯 Project Catch-Up Summary

## 📋 Project Overview
- **Project Name**: [name]
- **Purpose**: [1-2 sentences]
- **Current Phase**: [phase]
- **Tech Stack**: [key technologies]

## 📊 Current Status
- **Active Tasks**: [count] - [summary]
- **Recent Work**: [last session completion]
- **Git Status**: [branch, changes]

## ✅ Recent Accomplishments
[Bullet list from recent sessions]

## 🔄 What I'm Picking Up
**Immediate Next Steps** (from HANDOFF):
1. [Priority 1]
2. [Priority 2]
3. [Priority 3]

**Active Tasks**:
- [Task 1]: [status]
- [Task 2]: [status]

## 🚧 Active Blockers/Issues
[Unresolved issues]

## 🔑 Critical Context
[Important decisions, constraints]

## 📁 Key Files
[Files mentioned in handoff/tasks]

## 💡 Readiness
[Confirm understanding, ask which task to prioritize]
```

## Triggers
- User says: "catch up", "get context", "status", "where were we"
- New session starting (if HANDOFF-SESSION.md exists)
- User asks: "what should I work on", "what's next"

## Success Criteria
- All core docs read
- Session history reviewed
- Active tasks identified
- Git status checked
- Summary presented
- Next steps clear
```

**Migration Path**:
1. Create `.claude/skills/catch-up/SKILL.md`
2. Test with various trigger phrases
3. Keep `/catch-up` command initially for backward compatibility
4. Deprecate command after skill proven reliable

**Trade-offs**:
- **Pro**: Automatic invocation at session start
- **Pro**: No user action needed
- **Pro**: Read-only (safe)
- **Con**: May trigger when not needed (tunable with better description)

---

### 3.3 PRD Generation

**Current Implementation**: `.claude/tasks/1_create-prd.md` template
- Ask clarifying questions before writing
- Generate PRD using template structure
- Save to `1_backlog/[NUMBER]-[feature-name]/prd-[NUMBER]-[feature-name].md`
- Get user approval before task generation

**Conversion Assessment**: ⚠️ **Not Suitable for Pure Skill**

**Recommendation**: **Keep as Template + Agent Workflow**

**Reasoning**:
- Requires extensive user interaction (clarifying questions)
- Multi-turn conversation workflow
- Needs approval gates
- Best handled by tactical-product-manager agent
- PRD quality depends on thorough discovery

**Alternative**: **Create Helper Skill for PRD Structuring**

Instead of full PRD generation, create a skill that helps structure PRD content after discovery is complete.

**Skill Implementation**:

```yaml
---
name: prd-structuring
description: Structure and format Product Requirements Document (PRD) content
  into proper markdown format. Use when user has PRD information and needs it
  formatted, or when reviewing/refining existing PRD. Works with features,
  requirements, user stories, and acceptance criteria.
allowed-tools: Read, Write
---

# PRD Structuring Skill

## Purpose
Format PRD content into standardized structure using template.

## Instructions

1. **Read PRD Template**
   - Load `.claude/templates/prd-template.md` (if exists)
   - Or use standard PRD structure

2. **Organize Content**
   - Introduction/Overview
   - Goals (specific, measurable)
   - User Stories (detailed narratives)
   - Functional Requirements (numbered, clear)
   - Non-Goals (out of scope)
   - Design Considerations (optional)
   - Technical Considerations (optional)
   - Success Metrics
   - Open Questions

3. **Format Requirements**
   - Number all functional requirements
   - Use clear, concise language
   - Explicit and unambiguous phrasing
   - Suitable for junior developer understanding

4. **Determine File Location**
   - Check for existing PRD folders in `1_backlog/`
   - Use next sequential number (001, 002, 003...)
   - Create folder: `1_backlog/[NUMBER]-[feature-name]/`
   - Save as: `prd-[NUMBER]-[feature-name].md`

## Triggers
- User says: "format this as a PRD", "structure PRD", "create PRD document"
- User has PRD content that needs formatting
- Reviewing/refining existing PRD

## Note
This skill formats PRD content. For PRD discovery interviews and requirement
gathering, delegate to tactical-product-manager agent.
```

**Migration Path**:
1. Create `.claude/skills/prd-structuring/SKILL.md` for formatting
2. Keep `.claude/tasks/1_create-prd.md` as reference template
3. Continue using tactical-product-manager agent for full PRD workflow
4. Skill handles quick formatting; agent handles complex discovery

**Trade-offs**:
- **Pro**: Useful for quick PRD formatting
- **Pro**: Autonomous structuring
- **Con**: Doesn't replace full PRD workflow
- **Con**: Can't ask clarifying questions

---

### 3.4 Task List Generation

**Current Implementation**: `.claude/tasks/2_generate-tasks.md` template
- Analyze PRD requirements
- Generate parent tasks (pause for user approval)
- Generate subtasks after "Go" confirmation
- Assign agents to each task
- Identify relevant files

**Conversion Assessment**: ⚠️ **Not Suitable for Pure Skill**

**Recommendation**: **Keep as Agent Workflow**

**Reasoning**:
- Requires user approval gate (parent tasks → "Go" → subtasks)
- Interactive workflow with pauses
- Complex codebase analysis needed
- Agent assignment logic requires context
- Best handled by tactical-product-manager agent

**Alternative**: **No Conversion Recommended**

Task generation is inherently interactive and requires:
1. PRD analysis
2. Codebase review
3. Parent task generation
4. User confirmation ("Go")
5. Subtask decomposition
6. Agent assignment

This workflow is already optimal as agent delegation.

**Migration Path**: None - keep current implementation

**Trade-offs**: N/A

---

### 3.5 Task List Processing

**Current Implementation**: `.claude/tasks/3_process-task-list.md` template
- Move PRD folders between stages (backlog → active → completed → OBE)
- Mark tasks completed
- Run tests before commits
- Create commits for completed parent tasks
- Update project-status.md

**Conversion Assessment**: ⚠️ **Not Suitable for Skill**

**Recommendation**: **Keep as Workflow Protocol in CLAUDE.md**

**Reasoning**:
- Requires continuous user interaction and approval
- One subtask at a time with user permission
- Human judgment needed for task completion
- Test failures require investigation
- Git operations need review

**Alternative**: **Create Helper Skills for Specific Operations**

Break task processing into focused skills:

1. **`task-completion` skill** - Mark tasks complete, move folders
2. **`test-and-commit` skill** - Run tests, create commit if pass

**Skill 1: Task Completion**

```yaml
---
name: task-completion
description: Mark tasks as completed and move PRD folders between stages.
  Use when subtask finished, parent task completed, or feature done. Moves
  folders between 1_backlog, 2_active, 3_completed, or 0_obe. Updates
  project-status.md. Keywords: "mark complete", "task done", "move to completed".
allowed-tools: Read, Write, Bash
---

# Task Completion Skill

## Purpose
Mark tasks complete and manage PRD folder lifecycle.

## Instructions

1. **Identify Completion Type**
   - Subtask complete → Mark `[x]` in task file
   - Parent task complete → Mark parent `[x]`, check if all subtasks done
   - Feature complete → Move folder to 3_completed/

2. **Update Task File**
   - Read current task file
   - Change `[ ]` to `[x]` for completed items
   - Save updated file

3. **Move Folder If Needed**
   - Feature complete → `mv 2_active/001-feature/ 3_completed/`
   - OBE → `mv [current]/ 0_obe/` and add `obe-reason.md`
   - Starting work → `mv 1_backlog/001-feature/ 2_active/`

4. **Update Project Status**
   - Update `.claude/tasks/project-status.md`
   - Reflect new counts and status

## Triggers
- User says: "mark complete", "task done", "feature finished"
- All subtasks checked off
- Moving between backlog/active/completed
```

**Skill 2: Test and Commit**

```yaml
---
name: test-and-commit
description: Run test suite and create git commit if all tests pass. Use
  when parent task completed and ready to commit. Runs pytest/npm test/rails
  test based on project. Only commits if tests green. Keywords: "run tests",
  "commit changes", "tests and commit".
allowed-tools: Bash
---

# Test and Commit Skill

## Purpose
Ensure code quality by running tests before committing.

## Instructions

1. **Detect Test Framework**
   - Check for pytest.ini, package.json, or Gemfile
   - Determine test command

2. **Run Test Suite**
   - Execute: `pytest` or `npm test` or `bin/rails test`
   - Capture output

3. **Evaluate Results**
   - All tests pass → Proceed to commit
   - Any failures → Stop, report failures, do NOT commit

4. **Create Commit (Only If Tests Pass)**
   - Stage changes: `git add .`
   - Create commit with conventional format:
     ```
     git commit -m "feat: [task description]" \
                -m "- [change 1]" \
                -m "- [change 2]" \
                -m "Related to Task X.Y from prd-feature-name.md"
     ```

5. **Report Results**
   - Tests passed + committed → Confirm commit hash
   - Tests failed → Show failures, DO NOT commit

## Triggers
- User says: "run tests and commit", "commit if tests pass"
- Parent task complete, ready to commit
```

**Migration Path**:
1. Create both skills
2. Update CLAUDE.md to reference skills for task operations
3. Keep manual approval for starting subtasks
4. Skills handle deterministic operations only

**Trade-offs**:
- **Pro**: Automate repetitive operations
- **Pro**: Ensure test-commit discipline
- **Con**: Still need user approval for subtask progression
- **Con**: Complex workflows remain manual

---

### 3.6 Product Vision Creation

**Current Implementation**: strategic-product-visionary agent
- Conduct 5-phase discovery interview
- Ask clarifying questions (sequential pattern)
- Create `product-vision.md` and `strategic-themes.md`
- Populate templates with user responses

**Conversion Assessment**: ❌ **Not Suitable for Skill**

**Recommendation**: **Keep as Agent**

**Reasoning**:
- Highly interactive multi-turn conversation
- Requires strategic thinking and synthesis
- Sequential question pattern (CLAUDE.md mandate)
- Deep discovery interview (5 phases)
- Separate context window valuable
- Agent provides specialized product strategy expertise

**Alternative**: None - agent is optimal approach

**Migration Path**: None - keep current implementation

**Trade-offs**: N/A

---

### 3.7 Workspace Setup with Color Schemes

**Current Implementation**: `/onboard` command
- Detect project context
- Initialize VS Code workspace
- Analyze/apply color schemes
- Populate completion one-liners
- Optional product vision discovery
- Project context discovery (garden-guide)
- Agent hiring from gnomes library

**Conversion Assessment**: ⚠️ **Partial Conversion**

**Recommendation**: **Hybrid Approach**

**Reasoning**:
- Color scheme logic is deterministic → Good for skill
- Agent hiring requires analysis and decisions → Keep as command/agent
- Product vision requires interaction → Keep as agent
- User may want control over onboarding flow → Command better

**Hybrid Approach**:

1. **Keep `/onboard` command** for full interactive onboarding
2. **Create `workspace-colors` skill** for automatic color application

**Skill Implementation**:

```yaml
---
name: workspace-colors
description: Apply workspace color schemes to VS Code based on project name
  or keywords. Use when setting up new project, changing workspace theme, or
  user mentions "workspace colors", "theme", "color scheme". Detects project
  type and applies matching colors. Keywords: "setup workspace", "apply colors",
  "workspace theme".
allowed-tools: Read, Write, Bash
---

# Workspace Colors Skill

## Purpose
Automatically apply appropriate VS Code color scheme based on project context.

## Instructions

1. **Detect Project Context**
   - Get current directory name
   - Extract project name from git remote (if available)
   - Read README or package.json for project type

2. **Initialize VS Code Workspace**
   - Check if `.vscode/settings.json` exists
   - If not, create directory and copy template
   - If exists, read current colors

3. **Infer Theme from Project Name**
   - Match keywords against theme mapping:
     - "garden", "green" → Green + Brown
     - "ocean", "blue", "sea" → Blue + Teal
     - "fire", "red" → Red + Orange
     - "royal", "purple" → Purple + Blue
     - "tech", "cyber", "code" → Cyan + Navy

4. **Apply Color Scheme**
   - Update `.vscode/settings.json` with theme colors
   - Replace all 7 color instances
   - Update theme name comment

5. **Confirm Application**
   - Tell user: "Applied [Theme Name] workspace colors"
   - List primary and secondary colors
   - Suggest: "Reload VS Code to see changes"

## Theme Mappings
[Include full theme mapping from onboard.md]

## Triggers
- New project setup
- User says: "setup workspace", "apply colors", "change theme"
- .vscode/settings.json doesn't exist
```

**Migration Path**:
1. Create `.claude/skills/workspace-colors/SKILL.md`
2. Extract color logic from `/onboard` command
3. Keep full `/onboard` for comprehensive setup
4. Skill handles just color application

**Trade-offs**:
- **Pro**: Automatic color application for new projects
- **Pro**: Faster setup for simple cases
- **Con**: Loses interactive options (A/B/C selection)
- **Con**: Full onboarding still needs command

---

### 3.8 Agent Hiring from Gnome Folder / GitLab Marketplace

**Current Implementation**: Part of `/onboard` command (Step 11)
- Discover agents from `gnomes/agents/` library
- Analyze project context for domain keywords
- Build recommendation matrix
- Present categorized recommendations
- Copy selected agents to `.claude/agents/`

**Conversion Assessment**: ⚠️ **Complex - Not Recommended for Skill**

**Recommendation**: **Keep as Command/Agent Workflow**

**Reasoning**:
- Requires project context analysis
- User decision-making needed (A/B/C/D options)
- Dynamic library discovery
- Interactive selection process
- Best as part of onboarding command

**Alternative**: **Create Skill for Agent Discovery Only**

```yaml
---
name: agent-discovery
description: Discover and list available agents from gnomes library based on
  project needs. Use when user asks "what agents are available", "find agents
  for [domain]", "recommend agents", or during project setup. Analyzes project
  context and suggests matching specialists.
allowed-tools: Read, Bash, Grep, Glob
---

# Agent Discovery Skill

## Purpose
Help users find appropriate specialized agents for their project.

## Instructions

1. **Locate Agent Library**
   - Check: `../the_gnomes/agents/`
   - Or: `~/git/the_garden/the_gnomes/agents/`
   - If not found, report library unavailable

2. **Scan Available Agents**
   - List all .md files in library
   - Parse frontmatter (name, description, domain)
   - Categorize by prefix (software-, platform-, data-, etc.)

3. **Analyze Project Context**
   - Read PROJECT_CONTEXT.md for tech stack
   - Read PRODUCT_VISION.md for domain focus (if exists)
   - Extract keywords

4. **Build Recommendations**
   - Match keywords to agent categories
   - Create recommendation matrix
   - Prioritize by relevance

5. **Present Findings**
   ```markdown
   👥 Available Agents for Your Project

   **Recommended Based on Your Stack**:
   - [agent-name]: [description]

   **Available by Category**:
   - Software Engineering: [count] agents
   - Platform & Infrastructure: [count] agents
   - Data & Analytics: [count] agents
   [etc.]

   Use /onboard to hire agents, or manually copy from library.
   ```

## Triggers
- User asks: "what agents", "find agents", "recommend agents"
- Project setup phase
- New domain work starting
```

**Migration Path**:
1. Create `agent-discovery` skill for listing/recommending
2. Keep `/onboard` for full hiring workflow
3. Skill provides read-only discovery
4. Command handles interactive selection and copying

**Trade-offs**:
- **Pro**: Easy agent discovery
- **Pro**: Automatic recommendations
- **Con**: Doesn't copy agents (read-only)
- **Con**: Full hiring still needs command

---

### 3.9 Summary Report Writing for Subagent Work

**Current Implementation**: Agent session history protocol
- Every agent creates `[YYYYMMDD-HHMMSS]-[AGENT_TYPE]-[DESCRIPTION]-[SEQUENCE].md`
- Uses `.claude/templates/agent-session-template.md`
- Documents: task, work performed, decisions, deliverables, issues, metrics
- Main agent references in session handoff

**Conversion Assessment**: ✅ **Good Fit for Skill**

**Recommendation**: **Convert to Skill**

**Reasoning**:
- Deterministic workflow (same structure always)
- Template-based (agent-session-template.md)
- No user interaction needed
- Should trigger automatically when agent work completes
- Provides audit trail

**Skill Implementation**:

```yaml
---
name: agent-session-summary
description: Create session history documentation for specialized agent work.
  Use when agent completes task and returns to main agent. Generates timestamped
  session history file in .claude/context/agent-history/ using standard template.
  Keywords: "agent completed", "agent work done", "create agent history",
  "document agent session".
allowed-tools: Read, Write, Bash
---

# Agent Session Summary Skill

## Purpose
Document all specialized agent work with standardized session history files.

## Instructions

1. **Determine Session File Name**
   - Get current timestamp: `YYYYMMDD-HHMMSS`
   - Extract agent type (e.g., tactical-software-engineer)
   - Create description (1-4 words, kebab-case, max 25 chars)
   - Find next sequence number for this agent type today
   - Format: `[YYYYMMDD-HHMMSS]-[AGENT_TYPE]-[DESCRIPTION]-[SEQUENCE].md`

2. **Load Template**
   - Read `.claude/templates/agent-session-template.md`

3. **Populate Template Sections**
   - **Session Metadata**: Date, agent, duration, sequence
   - **Task Assignment**: What was requested, scope, success criteria
   - **Work Performed**: Detailed actions taken, decisions made
   - **Deliverables**: Files created/modified, outputs produced
   - **Recommendations**: Next steps, considerations for main agent
   - **Issues & Resolutions**: Problems encountered and solutions
   - **Performance Metrics**: Context usage, time, quality assessment

4. **Create History File**
   - Ensure directory exists: `.claude/context/agent-history/`
   - Write session file to directory
   - Confirm creation

5. **Return Summary to Main Agent**
   - Provide brief summary of work
   - Reference session file location
   - Highlight key deliverables and next steps

## Description Guidelines
- 1-4 words, hyphenated (kebab-case)
- Maximum ~25 characters
- Action-focused: verb-noun (e.g., `add-auth-flow`, `fix-bug`)
- Or feature-focused: noun-noun (e.g., `context-docs`, `api-security`)

## Example File Names
- `20251220-143022-tactical-software-engineer-add-auth-flow-001.md`
- `20251220-150815-tactical-cybersecurity-security-audit-001.md`
- `20251220-163045-tactical-platform-engineer-docker-setup-001.md`

## Triggers
- Agent work completed (before returning to main agent)
- User requests: "document agent work", "create session history"
- End of agent invocation
```

**Migration Path**:
1. Create `.claude/skills/agent-session-summary/SKILL.md`
2. Update agent template instructions to mention skill
3. Test with agent completions
4. Eventually embed in agent exit protocol

**Trade-offs**:
- **Pro**: Automatic session documentation
- **Pro**: Consistent audit trail
- **Pro**: Reduces agent burden
- **Con**: Shares main context (but small overhead)
- **Con**: May need refinement of descriptions

---

### 3.10 Garden Guide Flow (As Skill, Not Agent)

**Current Implementation**: garden-guide.md agent
- Project setup coaching
- Workflow guidance
- Context optimization
- System troubleshooting
- Interactive discovery interviews
- PROJECT_CONTEXT.md population

**Conversion Assessment**: ⚠️ **Partial Conversion Only**

**Recommendation**: **Hybrid - Keep Agent, Add Helper Skills**

**Reasoning**:
- Core value is interactive coaching → Requires agent
- Discovery interviews need multi-turn conversation
- Strategic guidance requires separate context
- Some helper functions could be skills

**Hybrid Approach**:

1. **Keep garden-guide agent** for complex setup and coaching
2. **Create helper skills** for specific quick operations:

**Skill 1: Quick Setup Validation**

```yaml
---
name: setup-validation
description: Validate Claude Agent System setup and identify configuration
  issues. Use when checking project setup, troubleshooting system, or user
  asks "is my setup correct", "validate setup", "check configuration". Reviews
  folder structure, required files, and common issues.
allowed-tools: Read, Bash, Glob, Grep
---

# Setup Validation Skill

## Purpose
Quick check of Claude Agent System setup and configuration.

## Instructions

1. **Check Required Files**
   - PROJECT_CONTEXT.md (exists and populated?)
   - CLAUDE.md (exists and loaded?)
   - .claude/agents/ (agents present?)
   - .claude/templates/ (templates available?)
   - .claude/tasks/ (directory structure correct?)

2. **Validate Directory Structure**
   - .claude/agents/
   - .claude/commands/
   - .claude/context/agent-history/
   - .claude/context/session-history/
   - .claude/docs/
   - .claude/hooks/
   - .claude/tasks/0_obe/, 1_backlog/, 2_active/, 3_completed/
   - .claude/templates/

3. **Check for Common Issues**
   - PROJECT_CONTEXT.md is template (not populated)
   - No agents in .claude/agents/
   - Missing required templates
   - Incorrect folder naming
   - Git not configured

4. **Report Findings**
   ```markdown
   ✅ Setup Validation Results

   **Status**: [Healthy / Issues Found]

   **Required Files**:
   - PROJECT_CONTEXT.md: [✓ Found / ✗ Missing / ⚠️ Template]
   - CLAUDE.md: [✓ Found / ✗ Missing]
   - Agents: [count] agents available

   **Issues Detected**:
   - [Issue 1]
   - [Issue 2]

   **Recommendations**:
   - [Fix 1]
   - [Fix 2]

   Run /onboard to complete setup, or invoke garden-guide agent for detailed help.
   ```

## Triggers
- User asks: "validate setup", "check configuration", "is my setup correct"
- New project setup
- Troubleshooting
```

**Skill 2: PROJECT_CONTEXT Template**

```yaml
---
name: project-context-template
description: Generate PROJECT_CONTEXT.md template with project-specific
  placeholders. Use when creating new PROJECT_CONTEXT.md or user asks to
  "create project context", "setup project context", "initialize context".
  Creates template ready for population.
allowed-tools: Read, Write
---

# PROJECT_CONTEXT Template Skill

## Purpose
Create PROJECT_CONTEXT.md template for user to fill out or guide to populate.

## Instructions

1. **Load Template**
   - Read `.claude/templates/project-context-template.md`

2. **Detect Project Info**
   - Get project name from directory or git remote
   - Detect tech stack from package.json, Gemfile, requirements.txt, etc.
   - Identify project type (web app, CLI, API, data pipeline)

3. **Customize Template**
   - Replace [Project Name] with actual name
   - Add detected tech stack to placeholders
   - Pre-fill what can be inferred

4. **Create PROJECT_CONTEXT.md**
   - Write file to project root
   - Confirm creation

5. **Guide User**
   ```
   ✅ Created PROJECT_CONTEXT.md template

   I've pre-filled some sections based on your project structure.

   Next steps:
   1. Review and complete the template sections
   2. Or invoke garden-guide agent for interactive discovery interview
   3. Or use /onboard for full setup workflow

   The more detail you provide, the better agents will understand your project.
   ```

## Triggers
- User says: "create project context", "setup context", "initialize PROJECT_CONTEXT"
- PROJECT_CONTEXT.md missing
- New project setup
```

**Migration Path**:
1. Create validation and template skills
2. Keep garden-guide agent for complex coaching
3. Skills handle quick checks and template creation
4. Agent handles discovery interviews and strategic guidance

**Trade-offs**:
- **Pro**: Quick validation and setup helpers
- **Pro**: Reduce agent invocations for simple tasks
- **Con**: Complex guidance still needs agent
- **Con**: Skills can't ask clarifying questions

---

### 3.11 Skills from docs/ Folder

**Current Documentation**: `.claude/docs/` contains:
- agent-invocation-examples.md
- code-quality-examples.md
- task-management-examples.md
- testing-and-implementation.md
- vision-workflow-guide.md
- strategic-agents-quick-start.md
- context-display-guide.md
- gitlab-cicd-guide.md
- tdd-workflow.md
- automated-context-tracking.md
- context-management-strategies.md

**Conversion Assessment**: ✅ **Several Good Candidates**

**Recommendation**: **Create Documentation-Based Skills**

**Reasoning**:
- Documentation contains best practices and patterns
- Can be packaged as autonomous guidance skills
- Claude can invoke when relevant tasks detected
- Reduces need to read full docs

**Recommended Skills**:

#### Skill 1: TDD Workflow

```yaml
---
name: tdd-workflow
description: Guide test-driven development workflow. Use when implementing
  code features, writing tests, or user mentions TDD, testing, unit tests,
  integration tests. Enforces write-test-first discipline. Keywords: "write
  tests", "TDD", "test-driven", "implement feature with tests".
---

# TDD Workflow Skill

## Purpose
Enforce test-driven development best practices.

## Mandatory Workflow

1. **Write Failing Test First**
   - Create test file: `[module].test.[ext]`
   - Write test that fails (feature not implemented yet)
   - Run test to confirm it fails
   - Commit: `test: add failing test for [feature]`

2. **Implement Minimal Code**
   - Write just enough code to make test pass
   - No extra features
   - Focus on test passing

3. **Verify Test Passes**
   - Run test suite
   - Confirm new test passes
   - Confirm no regressions

4. **Refactor If Needed**
   - Improve code quality
   - Keep tests green
   - Run tests after each refactor

5. **Integration/E2E Tests**
   - Add integration test if needed
   - Verify end-to-end workflow
   - Run full suite

6. **Commit Only If All Pass**
   - Stage changes
   - Commit with conventional format
   - Reference task and tests

## Quality Standards
- >80% unit test coverage for modified code
- Integration test for user-facing features
- All tests passing before commit
- Test files committed with code

## Enforcement
- Main agent verifies test files exist
- Tasks incomplete until tests pass
- Commits rejected without tests

[Include full TDD protocol from tdd-workflow.md]
```

#### Skill 2: GitLab CI/CD Guide

```yaml
---
name: gitlab-cicd
description: Write GitLab CI/CD pipeline configurations following YAML best
  practices. Use when creating .gitlab-ci.yml, configuring pipelines, or user
  mentions GitLab CI, pipeline, CI/CD. Prevents common YAML parsing errors.
---

# GitLab CI/CD Configuration Skill

## Purpose
Ensure GitLab CI/CD files follow correct YAML formatting.

## Critical Rules

1. **Use `--` Instead of `:` in Echo Statements**
   ```yaml
   # Wrong
   echo "Status: Running tests"

   # Correct
   echo "Status-- Running tests"
   ```

2. **Use YAML Literal Block Scalar for Multi-line Scripts**
   ```yaml
   script: |
     echo "Multi-line script"
     echo "Uses pipe character"
   ```

3. **Reference Arrays Directly for YAML Anchors**
   [Include examples from gitlab-cicd-guide.md]

[Full guide content]
```

#### Skill 3: Context Display

```yaml
---
name: context-display
description: Display context usage with emoji visualization bar. Use at end
  of every response to show token usage. Mandatory protocol. Keywords: "context
  usage", "token usage", "show context".
---

# Context Display Skill

## Purpose
Mandatory display of context usage at end of every response.

## Display Format
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 40% (80k/200k)
```

## Instructions

1. **Get Context Usage**
   - Use `/context` command output
   - Or estimate conservatively if between checks

2. **Calculate Blocks**
   - Total: 20 blocks
   - Filled blocks = round(percentage / 5)

3. **Apply Colors**
   - Blocks 1-12 (0-60%): 🟩 Green
   - Blocks 13-15 (60-75%): 🟨 Yellow
   - Blocks 16-17 (75-85%): 🟧 Orange
   - Blocks 18-20 (85-100%): 🟥 Red

4. **Add Status Messages**
   - 60-75%: "⚠️ Approaching handoff"
   - 75-85%: "🔄 Session handoff recommended"
   - 85%+: "🚨 New session recommended"

[Full display guide from context-display-guide.md]
```

#### Skill 4: Code Quality Comments

```yaml
---
name: code-quality-comments
description: Add clear, purposeful code comments using standard prefixes.
  Use when writing or reviewing code. Enforces REASON, WHY, NOTE, HACK, TODO
  comment patterns.
---

# Code Quality Comments Skill

## Purpose
Ensure code comments are clear, consistent, and purposeful.

## Comment Prefixes

- `REASON:` - Explain approach or method choice
- `WHY:` - Explain business logic or requirement
- `NOTE:` - Important detail or caveat
- `HACK:` - Temporary solution or workaround
- `TODO:` - Future improvement needed

[Include examples from code-quality-examples.md]
```

**Migration Path**:
1. Create skills for each doc-based best practice
2. Test invocation with relevant keywords
3. Keep docs as reference for comprehensive reading
4. Skills provide just-in-time guidance

**Trade-offs**:
- **Pro**: Best practices automatically applied
- **Pro**: Reduces need to read full docs
- **Pro**: Consistent enforcement
- **Con**: Skills less comprehensive than full docs
- **Con**: May trigger when not needed

---

## 4. Architecture Implications

### 4.1 Directory Structure

**Current Structure**:
```
.claude/
├── agents/              # Specialized agents (separate context)
├── commands/            # User-invoked commands
├── context/
│   ├── agent-history/   # Agent session files
│   └── session-history/ # Session handoff files
├── docs/                # Documentation
├── hooks/               # Git hooks, context tracker
├── tasks/               # PRD workflow folders
└── templates/           # Templates for PRDs, tasks, etc.
```

**With Skills Added**:
```
.claude/
├── agents/              # Keep for complex workflows
├── commands/            # Keep for interactive commands
├── context/
│   ├── agent-history/
│   └── session-history/
├── docs/                # Keep as reference docs
├── hooks/
├── skills/              # NEW: Autonomous workflows
│   ├── session-handoff/
│   │   └── SKILL.md
│   ├── catch-up/
│   │   └── SKILL.md
│   ├── agent-session-summary/
│   │   └── SKILL.md
│   ├── workspace-colors/
│   │   └── SKILL.md
│   ├── tdd-workflow/
│   │   └── SKILL.md
│   └── [other-skills]/
│       └── SKILL.md
├── tasks/
└── templates/
```

### 4.2 Workflow Integration

**Current Flow**:
```
User Request
    ↓
Main Agent (orchestrator)
    ↓
├─→ Delegate to Agent (complex work, separate context)
├─→ Execute Command (user-initiated utility)
└─→ Direct Implementation (minimal, discouraged)
```

**With Skills Integrated**:
```
User Request
    ↓
Main Agent (orchestrator)
    ↓
├─→ Skill Auto-Invoked (autonomous, description-matched)
├─→ Delegate to Agent (complex work, separate context)
├─→ Execute Command (user-initiated utility)
└─→ Direct Implementation (minimal, discouraged)
```

**Decision Logic**:
```
Task Received
    ↓
Is description-matched skill available?
    ├─ YES → Invoke skill autonomously
    └─ NO → Continue...
         ↓
User explicitly requested command (/cmd)?
    ├─ YES → Execute command
    └─ NO → Continue...
         ↓
Requires deep expertise + separate context?
    ├─ YES → Delegate to agent
    └─ NO → Minimal direct handling
```

### 4.3 Context Consumption Analysis

**Context Impact** (per 200k window):

| Workflow Type | Typical Cost | Notes |
|---------------|--------------|-------|
| **Skill** | 2-5k tokens | Shares main context; adds instructions + execution |
| **Command** | 3-8k tokens | Shares main context; may include templates |
| **Agent** | 15-30k tokens | Main agent briefing + output integration; agent has separate 200k |
| **Direct** | Variable | Depends on task complexity |

**Skill Overhead**:
- Skill discovery: ~1-2k tokens (all skill descriptions loaded)
- Skill execution: 2-5k tokens (instructions + work)
- **Total**: 3-7k tokens per skill invocation

**Implication**: Skills reduce context vs. agents (no briefing needed) but still consume main context. Use for focused, quick tasks only.

### 4.4 Skill Discovery Optimization

**Challenge**: Too many skills → high discovery overhead

**Mitigation Strategies**:

1. **Keep Skills Focused**
   - Each skill handles one specific workflow
   - Don't create "mega-skills" with broad scope

2. **Write Specific Descriptions**
   - Include exact keywords and triggers
   - Avoid vague terms
   - Make discovery signal-to-noise high

3. **Limit Total Skills**
   - Recommended: 10-15 skills maximum per project
   - More than 20 skills → discovery confusion
   - Use categories/prefixes for organization

4. **Periodic Review**
   - Remove unused skills
   - Merge overlapping skills
   - Refine descriptions based on usage

---

## 5. Implementation Roadmap

### Phase 1: Core Skills (Weeks 1-2)

**Goal**: Convert highest-value, lowest-risk workflows

**Skills to Create**:
1. ✅ `catch-up` - Session restoration (read-only, safe)
2. ✅ `agent-session-summary` - Agent work documentation
3. ✅ `context-display` - Context usage visualization
4. ✅ `setup-validation` - Quick setup checks

**Success Criteria**:
- All 4 skills created and tested
- Invocation working with trigger phrases
- No regressions in existing workflows
- User feedback collected

**Testing Checklist**:
- [ ] Skills auto-invoke correctly
- [ ] Descriptions precise enough for discovery
- [ ] No conflicts with existing commands
- [ ] Context overhead acceptable (<5k per skill)

---

### Phase 2: Helper Skills (Weeks 3-4)

**Goal**: Add specialized helper skills for common operations

**Skills to Create**:
1. ✅ `workspace-colors` - Automatic color scheme application
2. ✅ `tdd-workflow` - TDD guidance and enforcement
3. ✅ `gitlab-cicd` - CI/CD configuration best practices
4. ✅ `task-completion` - Mark tasks complete, move folders

**Success Criteria**:
- 4 additional skills operational
- Total 8 skills without discovery issues
- Measurable time savings on routine tasks
- Positive user feedback

---

### Phase 3: Advanced Skills (Weeks 5-6)

**Goal**: Add more complex automation

**Skills to Create**:
1. ✅ `session-handoff` - Automatic handoff at 75%
2. ✅ `test-and-commit` - Run tests, commit if pass
3. ✅ `prd-structuring` - Format PRD content
4. ✅ `code-quality-comments` - Comment prefix guidance

**Success Criteria**:
- 12 total skills operational
- Complex workflows (handoff, test-commit) working reliably
- Context tracking integrated with skills
- Documentation updated

---

### Phase 4: Optimization & Refinement (Weeks 7-8)

**Goal**: Optimize skill discovery and performance

**Tasks**:
1. Refine skill descriptions based on usage patterns
2. Merge or deprecate underused skills
3. Optimize context overhead
4. Update documentation and templates
5. Gather user feedback and iterate

**Success Criteria**:
- Skill invocation accuracy >90%
- Context overhead <5k average per skill
- User satisfaction with skills vs commands
- Clear migration guide for users

---

### Phase 5: Documentation & Training (Week 9)

**Goal**: Comprehensive documentation and user guidance

**Deliverables**:
1. Update CLAUDE.md with skill protocols
2. Create `.claude/docs/skills-guide.md`
3. Update onboarding to introduce skills
4. Create skill development guide for custom skills
5. Record usage metrics and optimization tips

**Success Criteria**:
- All skills documented
- Users understand when to use skills vs commands vs agents
- Clear examples for each skill
- Migration complete from old commands (where applicable)

---

## 6. Trade-offs and Recommendations

### 6.1 When to Use Skills

✅ **Use Skills For**:

1. **Autonomous, Repeatable Tasks**
   - Session restoration (catch-up)
   - Agent work documentation
   - Context usage display
   - Setup validation

2. **Template-Based Operations**
   - PRD formatting
   - Document structuring
   - File generation from templates

3. **Best Practice Enforcement**
   - TDD workflow guidance
   - Code comment standards
   - CI/CD configuration rules

4. **Quick Utilities**
   - Workspace color application
   - Task completion marking
   - Test-and-commit automation

---

### 6.2 When NOT to Use Skills

❌ **Don't Use Skills For**:

1. **Interactive Workflows**
   - PRD discovery interviews → Use agent
   - Product vision creation → Use agent
   - Complex user decisions → Use command

2. **Multi-Turn Conversations**
   - Clarifying questions → Use agent
   - Strategic planning → Use agent
   - Discovery interviews → Use agent

3. **Complex Orchestration**
   - Agent delegation logic → Keep in main agent
   - Multi-phase workflows → Use command + agents
   - Cross-domain tasks → Use agent swarms

4. **State-Dependent Operations**
   - Workflows requiring state tracking → Use commands
   - User approval gates → Use commands
   - Progressive disclosure → Use agents

---

### 6.3 Key Recommendations

#### Recommendation 1: Hybrid Approach is Optimal

**Finding**: Pure conversion (all skills OR all commands OR all agents) is suboptimal.

**Recommendation**: Use all three mechanisms strategically:
- **Skills**: Autonomous, focused, repeatable workflows
- **Commands**: User-initiated, interactive utilities
- **Agents**: Complex, deep-expertise, multi-turn work

**Rationale**: Each mechanism has distinct strengths. Combining them provides best user experience.

---

#### Recommendation 2: Start Conservative

**Finding**: Skills can interfere with workflows if poorly designed.

**Recommendation**:
- Start with 4-5 core skills (Phase 1)
- Validate each skill thoroughly before adding more
- Monitor context overhead and invocation accuracy
- Expand only after proven success

**Rationale**: Skills have discovery overhead. Too many poorly-defined skills → confusion and misfires.

---

#### Recommendation 3: Maintain Command Fallbacks

**Finding**: Skills may not always invoke when needed (description matching failure).

**Recommendation**:
- Keep existing commands as explicit fallbacks
- Document both skill (automatic) and command (manual) paths
- Allow users to bypass skill and use command if needed

**Rationale**: Provides reliability and user control.

---

#### Recommendation 4: Focus Skill Descriptions

**Finding**: Vague descriptions cause discovery failures.

**Recommendation**:
- Include 3-5 specific trigger keywords per skill
- Mention file types, operations, domain terms explicitly
- Test descriptions with various phrasings
- Iterate based on invocation success rates

**Rationale**: Description quality directly determines skill effectiveness.

---

#### Recommendation 5: Don't Convert Complex Workflows

**Finding**: Some workflows (PRD generation, vision creation, agent hiring) are inherently interactive.

**Recommendation**:
- Keep these as agent workflows
- Don't force-fit into skills
- Create helper skills for sub-tasks only (e.g., PRD formatting, not full PRD creation)

**Rationale**: Skills can't replace deep, interactive discovery conversations.

---

#### Recommendation 6: Monitor Context Overhead

**Finding**: Skills consume main agent context (unlike agents with separate windows).

**Recommendation**:
- Track context usage per skill invocation
- Aim for <5k tokens average per skill
- Deprecate high-overhead skills
- Use agents for context-heavy work

**Rationale**: Skills should reduce context burden, not increase it.

---

#### Recommendation 7: Version and Iterate

**Finding**: Skill descriptions and instructions need refinement over time.

**Recommendation**:
- Version skill files (v1.0.0, v1.1.0, etc.)
- Track changes in skill documentation
- Collect usage metrics (invocation rate, accuracy)
- Iterate based on data

**Rationale**: Skills improve with usage feedback and iteration.

---

### 6.4 Implementation Priority Matrix

| Workflow | Convert to Skill? | Priority | Complexity | Value |
|----------|-------------------|----------|------------|-------|
| **Catch-up** | ✅ Yes | High | Low | High |
| **Agent session summary** | ✅ Yes | High | Low | High |
| **Context display** | ✅ Yes | High | Low | High |
| **Setup validation** | ✅ Yes | Medium | Low | Medium |
| **Workspace colors** | ⚠️ Partial | Medium | Medium | Medium |
| **TDD workflow** | ✅ Yes | Medium | Low | High |
| **GitLab CI/CD** | ✅ Yes | Low | Low | Medium |
| **Task completion** | ✅ Yes | Medium | Medium | Medium |
| **Session handoff** | ⚠️ Partial | Medium | Medium | High |
| **Test-and-commit** | ✅ Yes | Low | Medium | Medium |
| **PRD structuring** | ⚠️ Helper only | Low | Medium | Low |
| **Code quality comments** | ✅ Yes | Low | Low | Medium |
| **PRD generation** | ❌ No (keep agent) | N/A | High | N/A |
| **Task list generation** | ❌ No (keep agent) | N/A | High | N/A |
| **Product vision** | ❌ No (keep agent) | N/A | High | N/A |
| **Agent hiring** | ⚠️ Discovery only | Low | High | Low |
| **Garden guide** | ❌ No (keep agent) | N/A | High | N/A |

**Legend**:
- ✅ Yes = Full conversion recommended
- ⚠️ Partial = Hybrid approach (skill + command/agent)
- ❌ No = Keep as agent or command

---

## 7. Conclusion

### Key Findings

1. **Skills are powerful but limited** - Best for autonomous, focused workflows; not suitable for complex interactive work

2. **Hybrid approach is optimal** - Combine skills (autonomous), commands (user-initiated), and agents (deep expertise)

3. **Start small and iterate** - Begin with 4-5 core skills, validate, then expand based on success

4. **Description quality is critical** - Precise descriptions with specific triggers determine skill effectiveness

5. **Don't force-fit complex workflows** - PRD generation, vision creation, and agent hiring should remain as agent workflows

### Recommended Immediate Actions

1. **Phase 1 Implementation** (Weeks 1-2):
   - Create `catch-up`, `agent-session-summary`, `context-display`, `setup-validation` skills
   - Test thoroughly with various trigger phrases
   - Validate context overhead <5k per skill

2. **Update Documentation**:
   - Add skills section to CLAUDE.md
   - Create `.claude/docs/skills-guide.md`
   - Update onboarding to mention skills

3. **Monitor and Iterate**:
   - Track skill invocation accuracy
   - Collect user feedback
   - Refine descriptions based on usage

### Long-Term Vision

**Goal**: Garden system leverages skills for routine automation while maintaining command and agent flexibility for complex work.

**Success Metrics**:
- 10-15 well-designed skills operational
- >90% invocation accuracy for skills
- <5k average context overhead per skill
- Positive user feedback on skills vs manual commands
- Documented migration guide for users

---

## Sources

- [How to create custom Skills | Claude Help Center](https://support.claude.com/en/articles/12512198-how-to-create-custom-skills)
- [Agent Skills - Claude Code Docs](https://code.claude.com/docs/en/skills)
- [GitHub - anthropics/skills: Public repository for Agent Skills](https://github.com/anthropics/skills)
- [Introducing Agent Skills | Claude](https://claude.com/blog/skills)
- [How to create Skills for Claude: steps and examples | Claude](https://claude.com/blog/how-to-create-skills-key-steps-limitations-and-examples)
- [Using Skills in Claude | Claude Help Center](https://support.claude.com/en/articles/12512180-using-skills-in-claude)
- [How to Create and Use Skills in Claude and Claude Code](https://apidog.com/blog/claude-skills/)
- [What are Skills? | Claude Help Center](https://support.claude.com/en/articles/12512176-what-are-skills)
- [Building Skills for Claude Code | Claude](https://claude.com/blog/building-skills-for-claude-code)
- [Claude Agent Skills: A First Principles Deep Dive](https://leehanchung.github.io/blogs/2025/10/26/claude-skills-deep-dive/)
