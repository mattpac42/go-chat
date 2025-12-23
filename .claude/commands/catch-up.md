# Catch-Up Command

**Task**: Get the Claude main agent fully up to speed with the current project state, context, recent work, and next steps.

## Usage

```bash
/catch-up
```

**No arguments required** - This command performs a comprehensive project state review.

## Purpose

This command is designed for when Claude (or a new Claude session) needs to quickly understand:
- What this project is about (PROJECT_CONTEXT.md)
- How Claude should operate in this project (CLAUDE.md)
- What work has been completed recently (session history)
- What needs to be done next (HANDOFF-SESSION.md and active tasks)

## Steps

### 1. Load Core Project Documentation

Read and analyze the following files in order:

**a) PROJECT_CONTEXT.md** (Project root)
- Understand the project's purpose, goals, and technical stack
- Note the project phase, architecture decisions, and key stakeholders
- Identify business logic, constraints, and compliance requirements
- Extract the project name and primary objectives

**b) CLAUDE.md** (Project root)
- Review agent orchestration rules and delegation protocols
- Understand the mandatory workflows (PRD, task management, session handoff)
- Note the context management protocols and quality requirements
- Understand the task file lifecycle and agent history requirements

**c) INITIAL_CONTEXT.md** (Project root, if exists)
- Review any initial project setup context or decisions
- Note any foundational constraints or requirements

### 2. Review Session History

**a) Most Recent Session File** (`.claude/work/history/SESSION-[highest-number].md`)
- Read the most recent completed session summary
- Note what was accomplished, decisions made, and issues encountered
- Understand the files that were modified and commands that were run
- Extract any important context or lessons learned

**b) Previous Session Files** (if multiple exist)
- Scan the 2-3 most recent session files for relevant historical context
- Look for recurring issues, patterns, or important decisions
- Note the progression of work across sessions

**c) HANDOFF-SESSION.md** (`.claude/work/history/HANDOFF-SESSION.md`)
- **CRITICAL**: This file contains the current work state and next steps
- Read the current status and what was being worked on
- Note all prioritized next steps
- Identify active blockers and issues
- Review critical context and recent decisions
- Note the exact files and locations being worked on
- Review the quick start guide for immediate next actions

### 3. Review Active Tasks

**a) Active Tasks Directory** (`.claude/work/2_active/`)
- List all active task files
- Read each active task file to understand:
  - Task objectives and success criteria
  - Current progress (completed vs pending subtasks)
  - Any blockers or notes
  - Commands that need to be run
  - Related PRDs or dependencies

**b) Task Priority**
- Cross-reference active tasks with HANDOFF-SESSION.md next steps
- Identify which task should be worked on first
- Understand dependencies between tasks

### 4. Check PRD Backlog (if applicable)

**a) PRD Folders** (`.claude/work/1_backlog/`, `.claude/work/2_active/`)
- Check if there are any PRD folders in backlog or active stages
- Read PRD files to understand planned or in-progress features
- Check for corresponding task files

**b) Project Status** (`.claude/work/project-status.md`, if exists)
- Review overall project status and feature tracking
- Note any high-level priorities or milestones

### 5. Review Recent Agent Work (if applicable)

**a) Agent History** (`.claude/work/history/`)
- Check for recent agent session files (last 5-10 files)
- Understand what specialized agents have been working on
- Note any agent recommendations or handoff notes
- Identify patterns or issues across agent sessions

### 6. Check Git Status

**a) Repository State**
```bash
git status
```
- Note any uncommitted changes
- Check current branch
- Identify any untracked files

**b) Recent Commits**
```bash
git log -5 --oneline
```
- Review the last 5 commits to understand recent work
- Note commit messages and patterns

### 7. Synthesize and Present Summary

After reading all the above, provide a **comprehensive but concise** summary organized as follows:

```markdown
# 🎯 Project Catch-Up Summary

## 📋 Project Overview
- **Project Name**: [name]
- **Purpose**: [1-2 sentence summary]
- **Current Phase**: [phase/stage]
- **Tech Stack**: [key technologies]

## 📊 Current Status
- **Active Tasks**: [count] - [brief summary]
- **Recent Work**: [what was completed in last session]
- **Git Status**: [branch, uncommitted changes]

## ✅ Recent Accomplishments
[Bullet list of key completed items from recent sessions]

## 🔄 What I'm Picking Up
**Immediate Next Steps** (from HANDOFF-SESSION.md):
1. [First priority with context]
2. [Second priority with context]
3. [Third priority with context]

**Active Tasks Requiring Attention**:
- [Task 1]: [status and next action]
- [Task 2]: [status and next action]

## 🚧 Active Blockers/Issues
[Any unresolved issues or blockers]

## 🔑 Critical Context to Remember
[Any important decisions, constraints, or context that affects current work]

## 📁 Key Files Currently Being Worked On
[List of files mentioned in HANDOFF or active tasks]

## 💡 My Understanding & Readiness
[Brief statement confirming understanding and readiness to proceed]

---

**Ready to continue!** Should I proceed with [the most logical next step based on handoff]?
```

### 8. Await User Confirmation

After presenting the summary:
- Ask if the summary is accurate
- Confirm which task/next step to prioritize
- Request any corrections or additional context needed
- Wait for explicit user instruction before starting work

## Success Criteria

- ✅ All core documentation files read and understood
- ✅ Session history reviewed (most recent + HANDOFF)
- ✅ All active tasks identified and understood
- ✅ Git status checked and noted
- ✅ Comprehensive summary provided to user
- ✅ Next steps clearly identified and prioritized
- ✅ User confirms understanding before proceeding

## Notes

- **This command is read-only** - it does NOT make any changes or start any work
- **Always wait for user confirmation** after presenting the summary
- **If files are missing**, note them in the summary and ask user for guidance
- **Be thorough but concise** - users want clarity without walls of text
- **Focus on actionability** - the summary should make it clear what to do next

## When to Use This Command

Use `/catch-up` when:
- Starting a new Claude session on an existing project
- Returning to a project after a break
- Taking over work from another session (after context handoff)
- Needing to refresh understanding of project state
- Before making any significant changes or decisions

## Integration with Existing Workflows

This command respects and leverages:
- **PRD Workflow**: Identifies active PRDs and their stage
- **Task Management**: Uses active tasks as source of truth
- **Session Handoff Protocol**: Relies on HANDOFF-SESSION.md as primary next-step guide
- **Agent History**: Incorporates specialized agent work and recommendations
- **Context Management**: Designed to operate within context budget efficiently
