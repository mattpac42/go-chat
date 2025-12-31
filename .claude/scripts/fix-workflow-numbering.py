#!/usr/bin/env python3
"""
Fix workflow step numbering in all marketplace agent files.
"""

import re
from pathlib import Path

MARKETPLACE_AGENTS = Path("/Users/mattpacione/git/the_garden/the_garden/marketplace/agents")
SKIP_FILES = {"README.md", "RECRUITMENT-CREW-README.md"}


def fix_numbering(filepath: Path) -> bool:
    """Fix workflow step numbering. Returns True if modified."""
    content = filepath.read_text()
    original = content

    # Find workflow section
    workflow_match = re.search(r'(## Workflow\n\n)((?:\d+\.\s+[^\n]+\n)+)', content)
    if not workflow_match:
        return False

    workflow_header = workflow_match.group(1)
    workflow_content = workflow_match.group(2)

    # Parse steps (extract content after number)
    steps = re.findall(r'\d+\.\s+([^\n]+)\n', workflow_content)

    # Rebuild with sequential numbering
    new_workflow = ""
    for i, step in enumerate(steps, 1):
        new_workflow += f"{i}. {step}\n"

    if new_workflow != workflow_content:
        content = content.replace(workflow_match.group(0), workflow_header + new_workflow)
        filepath.write_text(content)
        return True

    return False


def main():
    fixed = 0
    for filepath in sorted(MARKETPLACE_AGENTS.glob("*.md")):
        if filepath.name in SKIP_FILES:
            continue
        if fix_numbering(filepath):
            print(f"FIXED: {filepath.name}")
            fixed += 1

    print(f"\nFixed {fixed} files")


if __name__ == "__main__":
    main()
