#!/usr/bin/env python3
"""
Check for Garden updates from a planted project.

Usage:
    check_updates.py                    # Check and display
    check_updates.py --json             # Output as JSON
    check_updates.py --quiet            # Exit code only (0=no updates, 1=updates available)
    check_updates.py --changelog        # Show full changes since planted version
    check_updates.py --dismiss VERSION  # Dismiss notification for version
"""

import argparse
import json
import os
import re
import subprocess
import sys
from datetime import datetime
from pathlib import Path
from typing import Optional, Tuple, List, Dict, Any


def find_lineage() -> Optional[Path]:
    """Find lineage.json in current or parent directories."""
    current = Path.cwd()
    while current != current.parent:
        lineage_path = current / ".claude" / "lineage.json"
        if lineage_path.exists():
            return lineage_path
        current = current.parent
    return None


def parse_version(version_str: str) -> Tuple[int, int, int]:
    """Parse semantic version string to tuple."""
    match = re.match(r'(\d+)\.(\d+)\.(\d+)', version_str.strip())
    if match:
        return (int(match.group(1)), int(match.group(2)), int(match.group(3)))
    return (0, 0, 0)


def compare_versions(v1: str, v2: str) -> int:
    """Compare two version strings. Returns -1, 0, or 1."""
    t1, t2 = parse_version(v1), parse_version(v2)
    if t1 < t2:
        return -1
    if t1 > t2:
        return 1
    return 0


def read_garden_version_local(garden_path: Path) -> Optional[str]:
    """Read VERSION file from local Garden path."""
    version_file = garden_path / "VERSION"
    if version_file.exists():
        return version_file.read_text().strip()
    return None


def read_garden_version_remote() -> Optional[str]:
    """Read VERSION from git remote named 'garden'."""
    try:
        # Check if remote exists
        result = subprocess.run(
            ["git", "remote", "get-url", "garden"],
            capture_output=True,
            text=True,
            timeout=5
        )
        if result.returncode != 0:
            return None

        # Fetch from remote
        subprocess.run(
            ["git", "fetch", "garden", "--quiet"],
            capture_output=True,
            timeout=30
        )

        # Read VERSION from remote
        result = subprocess.run(
            ["git", "show", "garden/main:VERSION"],
            capture_output=True,
            text=True,
            timeout=5
        )
        if result.returncode == 0:
            return result.stdout.strip()
    except (subprocess.TimeoutExpired, FileNotFoundError):
        pass
    return None


def read_changelog_local(garden_path: Path) -> Optional[str]:
    """Read CHANGELOG.md from local Garden path."""
    changelog_file = garden_path / "CHANGELOG.md"
    if changelog_file.exists():
        return changelog_file.read_text()
    return None


def read_changelog_remote() -> Optional[str]:
    """Read CHANGELOG.md from git remote."""
    try:
        result = subprocess.run(
            ["git", "show", "garden/main:CHANGELOG.md"],
            capture_output=True,
            text=True,
            timeout=10
        )
        if result.returncode == 0:
            return result.stdout
    except (subprocess.TimeoutExpired, FileNotFoundError):
        pass
    return None


def parse_changelog(content: str, since_version: str) -> List[Dict[str, Any]]:
    """Parse CHANGELOG.md and return entries since given version."""
    entries = []
    current_version = None
    current_entry: Dict[str, Any] = {"version": None, "date": None, "sections": {}}
    current_section = None

    for line in content.split('\n'):
        # Version header: ## [3.2.0] - 2025-01-02
        version_match = re.match(r'## \[(\d+\.\d+\.\d+)\](?: - (\d{4}-\d{2}-\d{2}))?', line)
        if version_match:
            # Save previous entry if it's newer than since_version
            if current_version and compare_versions(current_version, since_version) > 0:
                entries.append(current_entry)

            current_version = version_match.group(1)
            current_entry = {
                "version": current_version,
                "date": version_match.group(2),
                "sections": {}
            }
            current_section = None
            continue

        # Section header: ### Added, ### Breaking, etc.
        section_match = re.match(r'### (\w+)', line)
        if section_match and current_version:
            current_section = section_match.group(1).lower()
            current_entry["sections"][current_section] = []
            continue

        # List item: - Something changed
        if line.startswith('- ') and current_version and current_section:
            current_entry["sections"][current_section].append(line[2:])

    # Don't forget the last entry
    if current_version and compare_versions(current_version, since_version) > 0:
        entries.append(current_entry)

    return entries


def check_updates(lineage_path: Path, quiet: bool = False) -> Dict[str, Any]:
    """Check for available updates."""
    lineage = json.loads(lineage_path.read_text())
    garden_path_str = lineage.get("garden", {}).get("source_path", "")
    garden_path = Path(garden_path_str) if garden_path_str else None

    result = {
        "garden_accessible": False,
        "source": None,
        "current_version": lineage.get("garden", {}).get("version", "unknown"),
        "latest_version": None,
        "update_available": False,
        "breaking_changes": False,
        "changes": [],
        "error": None
    }

    # Try local path first
    if garden_path and garden_path.exists():
        result["garden_accessible"] = True
        result["source"] = "local"
        result["latest_version"] = read_garden_version_local(garden_path)
        changelog_content = read_changelog_local(garden_path)
    else:
        # Fall back to git remote
        result["latest_version"] = read_garden_version_remote()
        if result["latest_version"]:
            result["garden_accessible"] = True
            result["source"] = "remote"
            changelog_content = read_changelog_remote()
        else:
            changelog_content = None

    if not result["garden_accessible"]:
        result["error"] = "Garden not accessible via local path or git remote"
        if not quiet:
            print(f"Warning: {result['error']}")
            if garden_path:
                print(f"  Local path: {garden_path} (not accessible)")
            print("  Git remote: 'garden' not configured")
            print("\nTo configure git remote:")
            print("  git remote add garden <garden-repo-url>")
        return result

    if not result["latest_version"]:
        result["error"] = "Could not determine Garden version"
        if not quiet:
            print(f"Warning: {result['error']}")
        return result

    current = result["current_version"]
    latest = result["latest_version"]

    if current == "unknown" or compare_versions(latest, current) > 0:
        result["update_available"] = True

        # Parse changelog for changes
        if changelog_content and current != "unknown":
            result["changes"] = parse_changelog(changelog_content, current)

            # Check for breaking changes
            for entry in result["changes"]:
                if "breaking" in entry.get("sections", {}):
                    result["breaking_changes"] = True
                    break

    return result


def display_update_notification(result: Dict[str, Any], verbose: bool = False):
    """Display update notification in Garden style."""
    if not result["update_available"]:
        current = result["current_version"]
        print(f"Garden is up to date (v{current})")
        if result.get("source"):
            print(f"  Source: {result['source']}")
        return

    current = result["current_version"]
    latest = result["latest_version"]

    # Visual notification bar
    print()
    print("=" * 60)

    if result["breaking_changes"]:
        status = "BREAKING CHANGES"
    else:
        status = "UPDATE AVAILABLE"

    print(f"Garden {status}: {current} -> {latest}")
    print("=" * 60)

    if result["changes"]:
        for entry in result["changes"]:
            if verbose:
                print(f"\n## v{entry['version']}" + (f" ({entry['date']})" if entry.get('date') else ""))

            # Show breaking changes first (always)
            if "breaking" in entry["sections"]:
                print("\nBREAKING:")
                for item in entry["sections"]["breaking"]:
                    print(f"   - {item}")

            if verbose:
                # Show other sections
                section_icons = {
                    "added": "ADDED",
                    "changed": "CHANGED",
                    "fixed": "FIXED",
                    "removed": "REMOVED",
                    "deprecated": "DEPRECATED"
                }
                for section, items in entry["sections"].items():
                    if section == "breaking":
                        continue
                    label = section_icons.get(section, section.upper())
                    print(f"\n{label}:")
                    for item in items[:5]:
                        print(f"   - {item}")
                    if len(items) > 5:
                        print(f"   ... and {len(items) - 5} more")
            else:
                # Summary mode - just show counts
                summary_parts = []
                for section, items in entry["sections"].items():
                    if section != "breaking" and items:
                        summary_parts.append(f"{len(items)} {section}")
                if summary_parts:
                    print(f"\n  {', '.join(summary_parts)}")

    print()
    print("Run `/sync-baseline` to update or `/updates --dismiss` to hide")
    print("=" * 60)
    print()


def dismiss_version(lineage_path: Path, version: str):
    """Add version to dismissed list."""
    lineage = json.loads(lineage_path.read_text())

    if "notifications" not in lineage:
        lineage["notifications"] = {}
    if "dismissed_versions" not in lineage["notifications"]:
        lineage["notifications"]["dismissed_versions"] = []

    if version not in lineage["notifications"]["dismissed_versions"]:
        lineage["notifications"]["dismissed_versions"].append(version)
        lineage_path.write_text(json.dumps(lineage, indent=2))
        print(f"Dismissed update notification for v{version}")
    else:
        print(f"v{version} was already dismissed")


def update_last_check(lineage_path: Path):
    """Update the last_check timestamp."""
    lineage = json.loads(lineage_path.read_text())

    if "notifications" not in lineage:
        lineage["notifications"] = {}

    lineage["notifications"]["last_check"] = datetime.now().isoformat()
    lineage_path.write_text(json.dumps(lineage, indent=2))


def main():
    parser = argparse.ArgumentParser(description="Check for Garden updates")
    parser.add_argument("--json", action="store_true", help="Output as JSON")
    parser.add_argument("--quiet", "-q", action="store_true", help="Quiet mode, exit code only")
    parser.add_argument("--changelog", action="store_true", help="Show full changelog")
    parser.add_argument("--dismiss", metavar="VERSION", help="Dismiss notification for version")
    args = parser.parse_args()

    lineage_path = find_lineage()
    if not lineage_path:
        if not args.quiet:
            print("Error: Not in a Garden project (no lineage.json found)")
        return 1

    # Handle dismiss
    if args.dismiss:
        dismiss_version(lineage_path, args.dismiss)
        return 0

    result = check_updates(lineage_path, args.quiet)

    # Update last check timestamp (unless quiet mode)
    if not args.quiet:
        update_last_check(lineage_path)

    if args.json:
        print(json.dumps(result, indent=2, default=str))
        return 0 if not result["update_available"] else 1

    if args.quiet:
        return 0 if not result["update_available"] else 1

    display_update_notification(result, verbose=args.changelog)
    return 0 if not result["update_available"] else 1


if __name__ == "__main__":
    sys.exit(main())
