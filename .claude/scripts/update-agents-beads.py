#!/usr/bin/env python3
"""
Batch update marketplace agents to be beads-aware.

Updates:
1. Removes `skills: agent-session-summary` from frontmatter
2. Adds beads steps to Workflow sections
3. Updates Handoff sections with beads commands
"""

import re
from pathlib import Path

MARKETPLACE_AGENTS = Path("/Users/mattpacione/git/the_garden/the_garden/marketplace/agents")

# Skip non-agent files
SKIP_FILES = {"README.md", "RECRUITMENT-CREW-README.md"}


def update_agent_file(filepath: Path) -> bool:
    """Update a single agent file. Returns True if modified."""
    content = filepath.read_text()
    original = content

    # 1. Remove skills: agent-session-summary from frontmatter
    content = re.sub(
        r'^(---\n.*?)skills:\s*agent-session-summary\n(.*?---)',
        r'\1\2',
        content,
        flags=re.MULTILINE | re.DOTALL
    )

    # 2. Update Workflow section - add beads steps and renumber
    if "## Workflow" in content:
        # Find workflow section
        workflow_match = re.search(r'(## Workflow\n\n)((?:\d+\.\s+[^\n]+\n)+)', content)
        if workflow_match:
            workflow_header = workflow_match.group(1)
            workflow_content = workflow_match.group(2)

            # Parse existing steps
            steps = re.findall(r'\d+\.\s+([^\n]+)\n', workflow_content)

            # Check if already has beads
            has_check_beads = any("Check Beads" in s for s in steps)
            has_update_beads = any("Update Beads" in s for s in steps)

            if not has_check_beads or not has_update_beads:
                new_steps = []

                # Add Check Beads at start if not present
                if not has_check_beads:
                    new_steps.append("**Check Beads**: Run `beads context` to see assigned work, mark bead as in-progress")

                # Add original steps (except Update Beads if we're adding it)
                for step in steps:
                    if "Check Beads" not in step and "Update Beads" not in step:
                        new_steps.append(step)

                # Add Update Beads before last step if not present
                if not has_update_beads:
                    # Insert before the last step
                    last_step = new_steps.pop()
                    new_steps.append("**Update Beads**: Close completed beads, add new beads for discovered work")
                    new_steps.append(last_step)

                # Rebuild with proper numbering
                new_workflow = ""
                for i, step in enumerate(new_steps, 1):
                    new_workflow += f"{i}. {step}\n"

                content = content.replace(workflow_match.group(0), workflow_header + new_workflow)

    # 3. Update Handoff section
    if "## Handoff" in content and "*Session history auto-created via" in content:
        # Replace the old handoff pattern
        content = re.sub(
            r'(\n## Handoff\n\nBefore returning control:\n)((?:\d+\.\s+[^\n]+\n)+)\n\*Session history auto-created via `agent-session-summary` skill\.\*',
            lambda m: f'{m.group(1)}1. Close completed beads with notes: `beads close <id> --note "summary"`\n2. Add beads for discovered work: `beads add "task" --parent <id>`\n{renumber_steps(m.group(2), start=3)}\n*Beads track execution state - no separate session files needed.*',
            content
        )

    if content != original:
        filepath.write_text(content)
        return True
    return False


def renumber_steps(steps_text: str, start: int) -> str:
    """Renumber steps starting from given number."""
    lines = steps_text.strip().split('\n')
    result = []
    num = start
    for line in lines:
        # Skip empty lines
        if not line.strip():
            continue
        # Replace leading number
        new_line = re.sub(r'^\d+\.', f'{num}.', line)
        result.append(new_line)
        num += 1
    return '\n'.join(result)


def main():
    updated = 0
    skipped = 0
    errors = []

    agent_files = list(MARKETPLACE_AGENTS.glob("*.md"))
    print(f"Found {len(agent_files)} markdown files in marketplace/agents/")

    for filepath in sorted(agent_files):
        if filepath.name in SKIP_FILES:
            print(f"  SKIP: {filepath.name}")
            skipped += 1
            continue

        try:
            if update_agent_file(filepath):
                print(f"  UPDATED: {filepath.name}")
                updated += 1
            else:
                print(f"  NO CHANGE: {filepath.name}")
        except Exception as e:
            print(f"  ERROR: {filepath.name} - {e}")
            errors.append((filepath.name, str(e)))

    print(f"\nSummary:")
    print(f"  Updated: {updated}")
    print(f"  Skipped: {skipped}")
    print(f"  Errors: {len(errors)}")

    if errors:
        print("\nErrors:")
        for name, err in errors:
            print(f"  {name}: {err}")


if __name__ == "__main__":
    main()
