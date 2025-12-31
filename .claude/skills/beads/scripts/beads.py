#!/usr/bin/env python3
"""
Beads - Git-backed issue tracking for AI agent workflows.

A lightweight implementation inspired by Steve Yegge's beads system.
Uses JSONL for persistence, integrates with Garden workflows.

Usage:
    beads.py init                      Initialize .beads/ directory
    beads.py add "title" [options]     Create a new bead
    beads.py list [filters]            List beads
    beads.py show <id>                 Show bead details
    beads.py update <id> [options]     Update a bead
    beads.py close <id> [--note]       Close a bead
    beads.py link <id> <rel> <target>  Link beads
    beads.py context                   Show session context
    beads.py import <tasks-file>       Import tasks from PRD
"""

import argparse
import hashlib
import json
import os
import sys
from datetime import datetime
from pathlib import Path
from typing import Optional

# Bead statuses
STATUS_OPEN = "open"
STATUS_IN_PROGRESS = "in-progress"
STATUS_BLOCKED = "blocked"
STATUS_CLOSED = "closed"

# Bead types
TYPE_TASK = "task"
TYPE_BUG = "bug"
TYPE_FEATURE = "feature"
TYPE_DISCOVERY = "discovery"

# Relationship types
REL_BLOCKS = "blocks"
REL_BLOCKED_BY = "blocked-by"
REL_PARENT = "parent"
REL_CHILD = "child"
REL_RELATED = "related"
REL_DISCOVERED_FROM = "discovered-from"


def find_beads_root() -> Optional[Path]:
    """Find .beads folder in current or parent directories."""
    current = Path.cwd()
    while current != current.parent:
        beads_path = current / ".beads"
        if beads_path.exists():
            return beads_path
        current = current.parent
    return None


def get_beads_path() -> Path:
    """Get or create .beads path."""
    root = find_beads_root()
    if root:
        return root
    return Path.cwd() / ".beads"


def generate_id(title: str, timestamp: str) -> str:
    """Generate short hash-based ID for a bead."""
    content = f"{title}{timestamp}"
    hash_val = hashlib.sha256(content.encode()).hexdigest()[:6]
    return f"bd-{hash_val}"


def load_beads(beads_path: Path) -> list[dict]:
    """Load all beads from JSONL file."""
    issues_file = beads_path / "issues.jsonl"
    if not issues_file.exists():
        return []

    beads = []
    with open(issues_file, "r") as f:
        for line in f:
            line = line.strip()
            if line:
                beads.append(json.loads(line))
    return beads


def save_beads(beads_path: Path, beads: list[dict]) -> None:
    """Save all beads to JSONL file."""
    issues_file = beads_path / "issues.jsonl"
    with open(issues_file, "w") as f:
        for bead in beads:
            f.write(json.dumps(bead) + "\n")


def find_bead(beads: list[dict], bead_id: str) -> Optional[dict]:
    """Find a bead by ID (supports partial match)."""
    for bead in beads:
        if bead["id"] == bead_id or bead["id"].endswith(bead_id):
            return bead
    return None


def is_blocked(bead: dict, beads: list[dict]) -> bool:
    """Check if a bead is blocked by any open beads."""
    blocked_by = bead.get("relationships", {}).get("blocked-by", [])
    for blocker_id in blocked_by:
        blocker = find_bead(beads, blocker_id)
        if blocker and blocker["status"] != STATUS_CLOSED:
            return True
    return False


def cmd_init(args) -> int:
    """Initialize .beads directory."""
    beads_path = Path.cwd() / ".beads"

    if beads_path.exists():
        print(f"Beads already initialized at {beads_path}")
        return 0

    beads_path.mkdir(parents=True)
    (beads_path / "issues.jsonl").touch()

    # Create config
    config = {
        "version": "1.0.0",
        "created": datetime.now().isoformat(),
        "settings": {
            "auto_close_children": True,
            "require_note_on_close": False
        }
    }
    with open(beads_path / "config.json", "w") as f:
        json.dump(config, f, indent=2)

    print(f"Initialized beads at {beads_path}")
    print("  issues.jsonl - Issue storage")
    print("  config.json  - Configuration")
    return 0


def cmd_add(args) -> int:
    """Create a new bead."""
    beads_path = get_beads_path()
    if not beads_path.exists():
        print("Error: Beads not initialized. Run 'beads init' first.")
        return 1

    beads = load_beads(beads_path)
    timestamp = datetime.now().isoformat()

    bead = {
        "id": generate_id(args.title, timestamp),
        "title": args.title,
        "type": args.type or TYPE_TASK,
        "status": STATUS_OPEN,
        "created": timestamp,
        "updated": timestamp,
        "agent": args.agent,
        "priority": args.priority or "medium",
        "tags": args.tags.split(",") if args.tags else [],
        "description": args.description or "",
        "files": args.files.split(",") if args.files else [],
        "relationships": {},
        "notes": []
    }

    # Handle parent relationship
    if args.parent:
        parent = find_bead(beads, args.parent)
        if parent:
            bead["relationships"]["parent"] = parent["id"]
            # Add child relationship to parent
            if "children" not in parent.get("relationships", {}):
                parent["relationships"]["children"] = []
            parent["relationships"]["children"].append(bead["id"])

    beads.append(bead)
    save_beads(beads_path, beads)

    print(f"Created: {bead['id']} - {bead['title']}")
    if args.agent:
        print(f"  Agent: {args.agent}")
    return 0


def cmd_list(args) -> int:
    """List beads with optional filters."""
    beads_path = find_beads_root()
    if not beads_path:
        print("Error: No .beads directory found.")
        return 1

    beads = load_beads(beads_path)

    # Apply filters
    filtered = beads

    if args.status:
        filtered = [b for b in filtered if b["status"] == args.status]
    elif not args.all:
        # By default, show non-closed
        filtered = [b for b in filtered if b["status"] != STATUS_CLOSED]

    if args.agent:
        filtered = [b for b in filtered if b.get("agent") == args.agent]

    if args.type:
        filtered = [b for b in filtered if b.get("type") == args.type]

    if args.ready:
        # Show only unblocked, open beads
        filtered = [b for b in filtered
                   if b["status"] == STATUS_OPEN and not is_blocked(b, beads)]

    if args.blocked:
        filtered = [b for b in filtered if is_blocked(b, beads)]

    # Output
    if not filtered:
        print("No beads found matching criteria.")
        return 0

    # Group by status for better readability
    by_status = {}
    for bead in filtered:
        status = bead["status"]
        if status not in by_status:
            by_status[status] = []
        by_status[status].append(bead)

    status_order = [STATUS_IN_PROGRESS, STATUS_OPEN, STATUS_BLOCKED, STATUS_CLOSED]
    status_icons = {
        STATUS_IN_PROGRESS: "🔄",
        STATUS_OPEN: "⭕",
        STATUS_BLOCKED: "🚫",
        STATUS_CLOSED: "✅"
    }

    for status in status_order:
        if status in by_status:
            print(f"\n{status_icons.get(status, '')} {status.upper()}")
            for bead in by_status[status]:
                agent_str = f" [{bead['agent']}]" if bead.get('agent') else ""
                blocked_str = " (blocked)" if is_blocked(bead, beads) else ""
                print(f"  {bead['id']}: {bead['title']}{agent_str}{blocked_str}")

    print(f"\nTotal: {len(filtered)} beads")
    return 0


def cmd_show(args) -> int:
    """Show detailed bead information."""
    beads_path = find_beads_root()
    if not beads_path:
        print("Error: No .beads directory found.")
        return 1

    beads = load_beads(beads_path)
    bead = find_bead(beads, args.id)

    if not bead:
        print(f"Error: Bead '{args.id}' not found.")
        return 1

    print(f"\n{'='*60}")
    print(f"ID: {bead['id']}")
    print(f"Title: {bead['title']}")
    print(f"Type: {bead['type']}")
    print(f"Status: {bead['status']}")
    print(f"Priority: {bead.get('priority', 'medium')}")
    if bead.get('agent'):
        print(f"Agent: {bead['agent']}")
    print(f"Created: {bead['created']}")
    print(f"Updated: {bead['updated']}")

    if bead.get('tags'):
        print(f"Tags: {', '.join(bead['tags'])}")

    if bead.get('description'):
        print(f"\nDescription:\n{bead['description']}")

    if bead.get('files'):
        print(f"\nFiles: {', '.join(bead['files'])}")

    if bead.get('relationships'):
        print(f"\nRelationships:")
        for rel_type, rel_ids in bead['relationships'].items():
            if isinstance(rel_ids, list):
                print(f"  {rel_type}: {', '.join(rel_ids)}")
            else:
                print(f"  {rel_type}: {rel_ids}")

    if bead.get('notes'):
        print(f"\nNotes:")
        for note in bead['notes']:
            print(f"  [{note['timestamp']}] {note['text']}")

    print(f"{'='*60}\n")
    return 0


def cmd_update(args) -> int:
    """Update a bead."""
    beads_path = find_beads_root()
    if not beads_path:
        print("Error: No .beads directory found.")
        return 1

    beads = load_beads(beads_path)
    bead = find_bead(beads, args.id)

    if not bead:
        print(f"Error: Bead '{args.id}' not found.")
        return 1

    # Update fields
    if args.status:
        bead["status"] = args.status
    if args.title:
        bead["title"] = args.title
    if args.agent:
        bead["agent"] = args.agent
    if args.priority:
        bead["priority"] = args.priority
    if args.description:
        bead["description"] = args.description

    if args.note:
        if "notes" not in bead:
            bead["notes"] = []
        bead["notes"].append({
            "timestamp": datetime.now().isoformat(),
            "text": args.note
        })

    bead["updated"] = datetime.now().isoformat()
    save_beads(beads_path, beads)

    print(f"Updated: {bead['id']}")
    return 0


def cmd_close(args) -> int:
    """Close a bead."""
    beads_path = find_beads_root()
    if not beads_path:
        print("Error: No .beads directory found.")
        return 1

    beads = load_beads(beads_path)
    bead = find_bead(beads, args.id)

    if not bead:
        print(f"Error: Bead '{args.id}' not found.")
        return 1

    bead["status"] = STATUS_CLOSED
    bead["closed"] = datetime.now().isoformat()
    bead["updated"] = bead["closed"]

    if args.note:
        if "notes" not in bead:
            bead["notes"] = []
        bead["notes"].append({
            "timestamp": bead["closed"],
            "text": f"Closed: {args.note}"
        })

    save_beads(beads_path, beads)

    # Check if this unblocks anything
    unblocked = []
    for other in beads:
        if other["id"] == bead["id"]:
            continue
        blocked_by = other.get("relationships", {}).get("blocked-by", [])
        if bead["id"] in blocked_by and not is_blocked(other, beads):
            unblocked.append(other)

    print(f"Closed: {bead['id']} - {bead['title']}")
    if unblocked:
        print("Unblocked:")
        for u in unblocked:
            print(f"  {u['id']}: {u['title']}")

    return 0


def cmd_link(args) -> int:
    """Create a relationship between beads."""
    beads_path = find_beads_root()
    if not beads_path:
        print("Error: No .beads directory found.")
        return 1

    beads = load_beads(beads_path)
    source = find_bead(beads, args.id)
    target = find_bead(beads, args.target)

    if not source:
        print(f"Error: Bead '{args.id}' not found.")
        return 1
    if not target:
        print(f"Error: Bead '{args.target}' not found.")
        return 1

    rel_type = args.relationship
    inverse_type = {
        REL_BLOCKS: REL_BLOCKED_BY,
        REL_BLOCKED_BY: REL_BLOCKS,
        REL_PARENT: REL_CHILD,
        REL_CHILD: REL_PARENT,
        REL_RELATED: REL_RELATED,
        REL_DISCOVERED_FROM: None
    }.get(rel_type)

    # Add relationship to source
    if "relationships" not in source:
        source["relationships"] = {}
    if rel_type not in source["relationships"]:
        source["relationships"][rel_type] = []
    if target["id"] not in source["relationships"][rel_type]:
        source["relationships"][rel_type].append(target["id"])

    # Add inverse relationship to target
    if inverse_type:
        if "relationships" not in target:
            target["relationships"] = {}
        if inverse_type not in target["relationships"]:
            target["relationships"][inverse_type] = []
        if source["id"] not in target["relationships"][inverse_type]:
            target["relationships"][inverse_type].append(source["id"])

    source["updated"] = datetime.now().isoformat()
    target["updated"] = source["updated"]
    save_beads(beads_path, beads)

    print(f"Linked: {source['id']} --{rel_type}--> {target['id']}")
    return 0


def cmd_context(args) -> int:
    """Show session context - what to work on."""
    beads_path = find_beads_root()
    if not beads_path:
        print("Error: No .beads directory found.")
        return 1

    beads = load_beads(beads_path)

    in_progress = [b for b in beads if b["status"] == STATUS_IN_PROGRESS]
    ready = [b for b in beads
             if b["status"] == STATUS_OPEN and not is_blocked(b, beads)]
    blocked = [b for b in beads if is_blocked(b, beads)]
    recently_closed = sorted(
        [b for b in beads if b["status"] == STATUS_CLOSED],
        key=lambda x: x.get("closed", ""),
        reverse=True
    )[:5]

    print("\n" + "="*60)
    print("BEADS SESSION CONTEXT")
    print("="*60)

    if in_progress:
        print(f"\n🔄 IN PROGRESS ({len(in_progress)})")
        for b in in_progress:
            agent_str = f" [{b['agent']}]" if b.get('agent') else ""
            print(f"   {b['id']}: {b['title']}{agent_str}")

    if ready:
        print(f"\n⭕ READY TO WORK ({len(ready)})")
        for b in ready[:10]:  # Show top 10
            agent_str = f" [{b['agent']}]" if b.get('agent') else ""
            print(f"   {b['id']}: {b['title']}{agent_str}")
        if len(ready) > 10:
            print(f"   ... and {len(ready) - 10} more")

    if blocked:
        print(f"\n🚫 BLOCKED ({len(blocked)})")
        for b in blocked[:5]:
            blockers = b.get("relationships", {}).get("blocked-by", [])
            print(f"   {b['id']}: {b['title']} (by: {', '.join(blockers)})")

    if recently_closed:
        print(f"\n✅ RECENTLY CLOSED")
        for b in recently_closed:
            print(f"   {b['id']}: {b['title']}")

    # Summary
    total = len(beads)
    open_count = len([b for b in beads if b["status"] == STATUS_OPEN])
    closed_count = len([b for b in beads if b["status"] == STATUS_CLOSED])

    print(f"\n{'='*60}")
    print(f"Total: {total} | Open: {open_count} | Closed: {closed_count}")
    print("="*60 + "\n")

    return 0


def cmd_import(args) -> int:
    """Import tasks from a PRD tasks file."""
    beads_path = get_beads_path()
    if not beads_path.exists():
        print("Initializing beads...")
        beads_path.mkdir(parents=True)
        (beads_path / "issues.jsonl").touch()

    tasks_file = Path(args.file)
    if not tasks_file.exists():
        print(f"Error: File '{args.file}' not found.")
        return 1

    content = tasks_file.read_text()
    beads = load_beads(beads_path)
    timestamp = datetime.now().isoformat()

    # Parse markdown tasks (simple parser)
    # Looks for: ## Task N: Title and - [ ] subtasks
    current_parent = None
    imported = []

    lines = content.split("\n")
    for line in lines:
        line = line.strip()

        # Parent task
        if line.startswith("## Task"):
            # Extract title after colon
            if ":" in line:
                title = line.split(":", 1)[1].strip()
            else:
                title = line.replace("## Task", "").strip()

            bead = {
                "id": generate_id(title, timestamp + str(len(imported))),
                "title": title,
                "type": TYPE_TASK,
                "status": STATUS_OPEN,
                "created": timestamp,
                "updated": timestamp,
                "agent": None,
                "priority": "medium",
                "tags": ["imported"],
                "description": "",
                "files": [],
                "relationships": {},
                "notes": []
            }
            beads.append(bead)
            imported.append(bead)
            current_parent = bead

        # Agent assignment
        elif line.startswith("**Agent**:") and current_parent:
            agent = line.replace("**Agent**:", "").strip()
            current_parent["agent"] = agent

        # Subtask
        elif line.startswith("- [ ]") and current_parent:
            # Extract subtask text
            subtask_text = line.replace("- [ ]", "").strip()
            # Remove leading numbers like "1.1 "
            if subtask_text and subtask_text[0].isdigit():
                parts = subtask_text.split(" ", 1)
                if len(parts) > 1:
                    subtask_text = parts[1]

            # Extract files if present
            files = []
            if "Files:" in subtask_text:
                text_part, files_part = subtask_text.split("Files:", 1)
                subtask_text = text_part.strip().rstrip(" -")
                files = [f.strip() for f in files_part.split(",")]

            bead = {
                "id": generate_id(subtask_text, timestamp + str(len(imported))),
                "title": subtask_text,
                "type": TYPE_TASK,
                "status": STATUS_OPEN,
                "created": timestamp,
                "updated": timestamp,
                "agent": current_parent.get("agent"),
                "priority": "medium",
                "tags": ["imported", "subtask"],
                "description": "",
                "files": files,
                "relationships": {
                    "parent": current_parent["id"]
                },
                "notes": []
            }

            # Add child reference to parent
            if "children" not in current_parent["relationships"]:
                current_parent["relationships"]["children"] = []
            current_parent["relationships"]["children"].append(bead["id"])

            beads.append(bead)
            imported.append(bead)

    save_beads(beads_path, beads)

    print(f"Imported {len(imported)} beads from {args.file}")
    for b in imported:
        parent_indicator = "" if not b["relationships"].get("parent") else "  "
        print(f"{parent_indicator}{b['id']}: {b['title']}")

    return 0


def cmd_progress(args) -> int:
    """Mark a bead as in-progress."""
    beads_path = find_beads_root()
    if not beads_path:
        print("Error: No .beads directory found.")
        return 1

    beads = load_beads(beads_path)
    bead = find_bead(beads, args.id)

    if not bead:
        print(f"Error: Bead '{args.id}' not found.")
        return 1

    bead["status"] = STATUS_IN_PROGRESS
    bead["updated"] = datetime.now().isoformat()

    save_beads(beads_path, beads)
    print(f"In progress: {bead['id']} - {bead['title']}")
    return 0


def main():
    parser = argparse.ArgumentParser(
        description="Beads - Git-backed issue tracking for AI workflows"
    )
    subparsers = parser.add_subparsers(dest="command", help="Commands")

    # init
    subparsers.add_parser("init", help="Initialize .beads directory")

    # add
    add_parser = subparsers.add_parser("add", help="Create a new bead")
    add_parser.add_argument("title", help="Bead title")
    add_parser.add_argument("--type", "-t", choices=[TYPE_TASK, TYPE_BUG, TYPE_FEATURE, TYPE_DISCOVERY])
    add_parser.add_argument("--agent", "-a", help="Assigned agent")
    add_parser.add_argument("--priority", "-p", choices=["low", "medium", "high", "critical"])
    add_parser.add_argument("--tags", help="Comma-separated tags")
    add_parser.add_argument("--description", "-d", help="Description")
    add_parser.add_argument("--files", "-f", help="Comma-separated file paths")
    add_parser.add_argument("--parent", help="Parent bead ID")

    # list
    list_parser = subparsers.add_parser("list", help="List beads")
    list_parser.add_argument("--status", "-s", choices=[STATUS_OPEN, STATUS_IN_PROGRESS, STATUS_BLOCKED, STATUS_CLOSED])
    list_parser.add_argument("--agent", "-a", help="Filter by agent")
    list_parser.add_argument("--type", "-t", help="Filter by type")
    list_parser.add_argument("--ready", "-r", action="store_true", help="Show only ready (unblocked) beads")
    list_parser.add_argument("--blocked", "-b", action="store_true", help="Show only blocked beads")
    list_parser.add_argument("--all", action="store_true", help="Include closed beads")

    # show
    show_parser = subparsers.add_parser("show", help="Show bead details")
    show_parser.add_argument("id", help="Bead ID")

    # update
    update_parser = subparsers.add_parser("update", help="Update a bead")
    update_parser.add_argument("id", help="Bead ID")
    update_parser.add_argument("--status", "-s", choices=[STATUS_OPEN, STATUS_IN_PROGRESS, STATUS_BLOCKED, STATUS_CLOSED])
    update_parser.add_argument("--title", help="New title")
    update_parser.add_argument("--agent", "-a", help="Assigned agent")
    update_parser.add_argument("--priority", "-p", choices=["low", "medium", "high", "critical"])
    update_parser.add_argument("--description", "-d", help="Description")
    update_parser.add_argument("--note", "-n", help="Add a note")

    # close
    close_parser = subparsers.add_parser("close", help="Close a bead")
    close_parser.add_argument("id", help="Bead ID")
    close_parser.add_argument("--note", "-n", help="Closing note")

    # progress
    progress_parser = subparsers.add_parser("progress", help="Mark bead as in-progress")
    progress_parser.add_argument("id", help="Bead ID")

    # link
    link_parser = subparsers.add_parser("link", help="Link beads")
    link_parser.add_argument("id", help="Source bead ID")
    link_parser.add_argument("relationship", choices=[REL_BLOCKS, REL_BLOCKED_BY, REL_PARENT, REL_CHILD, REL_RELATED, REL_DISCOVERED_FROM])
    link_parser.add_argument("target", help="Target bead ID")

    # context
    subparsers.add_parser("context", help="Show session context")

    # import
    import_parser = subparsers.add_parser("import", help="Import from tasks file")
    import_parser.add_argument("file", help="Path to tasks markdown file")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        return 0

    commands = {
        "init": cmd_init,
        "add": cmd_add,
        "list": cmd_list,
        "show": cmd_show,
        "update": cmd_update,
        "close": cmd_close,
        "progress": cmd_progress,
        "link": cmd_link,
        "context": cmd_context,
        "import": cmd_import,
    }

    return commands[args.command](args)


if __name__ == "__main__":
    sys.exit(main())
