# Claude Agent System Documentation

This directory contains reference documentation for the Claude Agent orchestration system.

## Quick Reference

### Core Rules (CLAUDE.md)
- **PRD Workflow**: Mandatory structured workflow for feature development
- **Plan Adherence Protocol**: Zero deviation policy for technology stacks
- **Context Display**: Mandatory visual display at end of every response
- **Session Handoff**: Automatic at 60% context usage
- **Project Discovery**: Required before system setup
- **Agent Selection**: Domain-based routing to specialized agents
- **Task Management**: Mandatory task file workflow

### Extended Documentation

#### Examples and Templates
- **agent-invocation-examples.md**: Detailed agent briefing templates and examples
- **code-quality-examples.md**: Code commenting examples and safety protocols
- **task-management-examples.md**: Task discovery and completion examples

#### Implementation Details
- **testing-and-implementation.md**: Repository setup, integration commands, testing protocols

## File Token Usage

Approximate token usage for reference:
- CLAUDE.md (core rules): ~9.2k tokens (optimized)
- Extended docs: ~2-3k tokens each (loaded on demand)

## Usage Pattern

1. **Main Claude agent**: Loads CLAUDE.md automatically (core rules only)
2. **Reference as needed**: Points to extended docs for detailed examples
3. **Context efficiency**: Reduces main file bloat while preserving comprehensive guidance
