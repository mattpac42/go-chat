#!/usr/bin/env python3
"""
Project Initializer - Creates a new project or adds Garden to existing repository

Usage:
    init_project.py <project-name> --path <path> --type <project-type> [--agents agent1,agent2,...] [--description "Project description"] [--stack python-flask] [--existing]

Examples:
    # New project
    init_project.py my-web-app --path ~/projects --type webapp --agents developer,architect,ux-tactical --stack python-flask
    init_project.py api-service --path /work/projects --type api --agents developer,architect,platform-tactical --stack python-fastapi
    init_project.py my-cli --path .. --type cli --agents developer,architect

    # Existing repository
    init_project.py my-repo --path ~/existing/repo --type webapp --agents developer,architect --existing
"""

import sys
import json
import re
import shutil
import subprocess
from pathlib import Path
from datetime import datetime


# Agent descriptions for CLAUDE.md generation
AGENT_DESCRIPTIONS = {
    # Core agents
    "developer": "Code implementation, testing, debugging",
    "architect": "System design, patterns, technical strategy",
    "product": "Requirements, PRDs, feature planning",
    "platform": "Infrastructure, DevOps, CI/CD",
    "researcher": "Analysis, exploration, information gathering",
    "garden-guide": "Project setup, workflow guidance",
    "project-navigator": "Project-specific context and memory",
    "prompt-optimizer": "Prompt engineering and optimization",
    "product-visionary": "Strategic product vision",
    # Specialized agents (common ones)
    "software-tactical": "Hands-on code implementation",
    "software-strategic": "Architecture and technical direction",
    "platform-tactical": "Infrastructure hands-on work",
    "platform-strategic": "Cloud strategy and planning",
    "ux-tactical": "UI implementation and components",
    "ux-strategic": "Design systems and UX strategy",
    "cicd-tactical": "Pipeline implementation",
    "cicd-strategic": "CI/CD architecture",
    "product-manager-tactical": "Sprint execution and tickets",
    "product-manager-strategic": "Roadmap and product strategy",
    "sre-tactical": "Reliability implementation",
    "sre-strategic": "SRE strategy and planning",
    "business-strategic": "Business strategy and planning",
}

# Tech stack to devcontainer template mapping
TECH_STACKS = {
    "python-flask": {
        "template_dir": "python-flask",
        "display_name": "Python Flask",
        "services": [{"name": "flask", "port": 5000, "label": "Flask App"}],
    },
    "python-fastapi": {
        "template_dir": "python-fastapi",
        "display_name": "Python FastAPI",
        "services": [{"name": "fastapi", "port": 8000, "label": "FastAPI App"}],
    },
    "node-react": {
        "template_dir": "node-react",
        "display_name": "React",
        "services": [{"name": "react", "port": 3000, "label": "React App"}],
    },
    "node-nextjs": {
        "template_dir": "node-nextjs",
        "display_name": "Next.js",
        "services": [{"name": "nextjs", "port": 3000, "label": "Next.js App"}],
    },
}

# Port registry for dynamic service detection
PORT_REGISTRY = {
    # Frontend frameworks
    "nextjs": {"port": 3000, "label": "Next.js App", "category": "frontend"},
    "react": {"port": 3000, "label": "React App", "category": "frontend"},
    "vite": {"port": 5173, "label": "Vite Dev Server", "category": "frontend"},
    "vue": {"port": 5173, "label": "Vue App", "category": "frontend"},
    "angular": {"port": 4200, "label": "Angular App", "category": "frontend"},

    # Backend frameworks
    "fastapi": {"port": 8000, "label": "FastAPI Server", "category": "backend"},
    "flask": {"port": 5000, "label": "Flask App", "category": "backend"},
    "django": {"port": 8000, "label": "Django Server", "category": "backend"},
    "express": {"port": 3001, "label": "Express Server", "category": "backend"},
    "go": {"port": 8080, "label": "Go Server", "category": "backend"},
    "rust": {"port": 8080, "label": "Rust Server", "category": "backend"},
    "spring": {"port": 8080, "label": "Spring Boot", "category": "backend"},
}


def detect_services(description: str) -> list[dict]:
    """Parse project description to detect frontend/backend services."""
    description_lower = description.lower()
    detected = []

    patterns = [
        (r'\bnext\.?js\b', 'nextjs'),
        (r'\breact\b', 'react'),
        (r'\bvite\b', 'vite'),
        (r'\bvue\b', 'vue'),
        (r'\bangular\b', 'angular'),
        (r'\bfastapi\b', 'fastapi'),
        (r'\bflask\b', 'flask'),
        (r'\bdjango\b', 'django'),
        (r'\bexpress\b', 'express'),
        (r'\bgo\b(?:\s+api|\s+backend|\s+server)?', 'go'),
        (r'\brust\b', 'rust'),
        (r'\bspring\b', 'spring'),
    ]

    for pattern, service_name in patterns:
        if re.search(pattern, description_lower):
            if service_name in PORT_REGISTRY:
                detected.append({"name": service_name, **PORT_REGISTRY[service_name]})

    return detected


def resolve_port_conflicts(services: list[dict]) -> list[dict]:
    """Resolve port conflicts by incrementing duplicate ports."""
    used_ports = set()
    resolved = []

    for service in services:
        port = service["port"]
        while port in used_ports:
            port += 1
        used_ports.add(port)
        resolved.append({**service, "port": port})

    return resolved


# Which project types support which devcontainer stacks
PROJECT_TYPE_STACKS = {
    "webapp": ["node-react", "node-nextjs", "python-flask"],
    "api": ["python-flask", "python-fastapi", "node-nextjs"],
    "cli": [],
    "mobile": [],
    "library": ["python-flask", "node-react"],
    "data": ["python-flask", "python-fastapi"],
    "devops": [],
    "business": [],
}


def get_garden_root():
    """Get the root directory of The Garden (where this script lives)."""
    script_path = Path(__file__).resolve()
    # Script is at .claude/skills/plant-project/scripts/init_project.py
    # Garden root is 5 levels up
    return script_path.parents[4]


def get_git_commit_hash(garden_root):
    """Get the current git commit hash of The Garden."""
    try:
        result = subprocess.run(
            ["git", "rev-parse", "HEAD"],
            cwd=garden_root,
            capture_output=True,
            text=True
        )
        if result.returncode == 0:
            return result.stdout.strip()[:8]  # Short hash
    except Exception:
        pass
    return "unknown"


def get_garden_version(garden_root):
    """Get the current version of The Garden from VERSION file."""
    version_file = garden_root / "VERSION"
    if version_file.exists():
        return version_file.read_text().strip()
    return "unknown"


def get_garden_remote_url(garden_root):
    """Get the git remote URL of The Garden for remote sync support."""
    try:
        result = subprocess.run(
            ["git", "remote", "get-url", "origin"],
            cwd=garden_root,
            capture_output=True,
            text=True
        )
        if result.returncode == 0:
            return result.stdout.strip()
    except Exception:
        pass
    return None


def init_beads(project_dir, project_name):
    """Initialize .beads directory for execution state tracking.

    Creates the beads structure so projects can use beads immediately
    without requiring manual 'beads init'.
    """
    beads_path = project_dir / ".beads"

    if beads_path.exists():
        return beads_path  # Already initialized

    beads_path.mkdir(parents=True)
    (beads_path / "issues.jsonl").touch()

    # Create config matching beads.py structure
    config = {
        "version": "1.0.0",
        "created": datetime.now().isoformat(),
        "project": project_name,
        "settings": {
            "auto_close_children": True,
            "require_note_on_close": False
        }
    }
    with open(beads_path / "config.json", "w") as f:
        json.dump(config, f, indent=2)

    return beads_path


def title_case_name(name):
    """Convert hyphenated name to Title Case for display."""
    return ' '.join(word.capitalize() for word in name.split('-'))


def generate_agent_table(agents):
    """Generate markdown table of agents for CLAUDE.md."""
    lines = []
    for agent in agents:
        description = AGENT_DESCRIPTIONS.get(agent, "Specialized agent")
        lines.append(f"| {agent} | {description} |")
    return '\n'.join(lines)


def create_lineage_json(project_dir, garden_root, project_name, project_type, description, agents, tech_stack=None):
    """Create lineage.json to track Garden connection."""
    commit_hash = get_git_commit_hash(garden_root)
    garden_version = get_garden_version(garden_root)
    remote_url = get_garden_remote_url(garden_root)

    lineage = {
        "schema_version": "2.1",
        "garden": {
            "source_path": str(garden_root),
            "commit_hash": commit_hash,
            "version": garden_version,
            "remote_url": remote_url
        },
        "project": {
            "name": project_name,
            "type": project_type,
            "tech_stack": tech_stack,
            "description": description,
            "rooted_at": datetime.now().isoformat()
        },
        "included": {
            "agents": agents,
            "templates": [
                "agent-template.md",
                "prd.md",
                "product-vision-template.md",
                "task.md",
                "tdd-task-template.md",
                "handoff.md",
                "session-summary-template.md"
            ],
            "skills": [
                "catch-up",
                "handoff",
                "context-display",
                "agent-session-summary",
                "setup-validation",
                "version-notify",
                "garden-request"
            ],
            "commands": [
                "commit.md",
                "catch-up.md",
                "handoff.md",
                "mr.md",
                "onboard.md",
                "sync-baseline.md",
                "updates.md",
                "garden-request.md"
            ],
            "devcontainer": tech_stack is not None
        },
        "sync": {
            "enabled": True,
            "last_sync": None,
            "last_version": None,
            "auto_notify": True,
            "excluded_paths": []
        },
        "notifications": {
            "dismissed_versions": [],
            "last_check": None
        }
    }

    lineage_path = project_dir / ".claude" / "lineage.json"
    lineage_path.write_text(json.dumps(lineage, indent=2))
    return lineage_path


def create_claude_md(project_dir, project_name, description, agents, rooted_date):
    """Create customized CLAUDE.md for the new project."""
    agent_table = generate_agent_table(agents)

    content = f"""# {title_case_name(project_name)}

{description}

**Main agent orchestrates. Agents implement. No exceptions.**

## Agents

| Agent | Domain |
|-------|--------|
{agent_table}

## Context Thresholds

| Level | Action |
|-------|--------|
| 60% | Warning - approaching limit |
| 75% | Handoff triggered automatically |
| 85% | New session recommended |

Display context bar at 50%+ or after agent completion:
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 50% (100k/200k)
```

## Quick Rules

1. **Delegate all implementation** to specialized agents
2. **Main agent only**: reads context, asks questions, invokes agents, tracks progress
3. **Use skills** for workflows (PRD creation, task management, etc.)
4. **Parallel execution**: Run independent agent tasks simultaneously
5. **Session handoff**: Auto-created at 75% context

## Key Files

- `.claude/PROTOCOLS.md` - Detailed rules and workflows
- `.claude/QUICKSTART.md` - 2-minute setup guide
- `.claude/skills/` - Workflow skills
- `.claude/agents/` - Project agents
- `.claude/lineage.json` - Garden connection for sync

## Delegation Decision

Ask four questions:
1. Is this specialized work? → Delegate
2. Will this use >10k tokens? → Delegate
3. Is this my third attempt? → Delegate
4. Are there 2+ independent tasks? → Delegate in **PARALLEL**

If NO to all four → Handle directly

## Parallel Execution

**Default to parallel when possible.** Each agent has its own 200k context window.

**Parallelize when:**
- Multiple independent research tasks
- Different domains (e.g., frontend + backend + infra)
- Unrelated file changes
- Reviews or analysis of separate components

**Invoke parallel agents in a SINGLE message with multiple Task tool calls.**

---

## Garden Lineage

This project was rooted from The Garden on {rooted_date}.

To sync updates from The Garden:
```
/sync-baseline
```

For detailed protocols: Read `.claude/PROTOCOLS.md`
"""

    claude_md_path = project_dir / "CLAUDE.md"
    claude_md_path.write_text(content)
    return claude_md_path


def copy_core_structure(garden_root, project_dir, agents):
    """Copy core .claude structure to new project."""
    garden_claude = garden_root / ".claude"
    project_claude = project_dir / ".claude"

    # Create .claude directory
    project_claude.mkdir(parents=True, exist_ok=True)

    # Copy core files
    core_files = [
        "PROTOCOLS.md",
        "QUICKSTART.md",
        "PROJECT.md",
        "settings.json",
    ]

    for filename in core_files:
        src = garden_claude / filename
        if src.exists():
            dst = project_claude / filename
            shutil.copy2(src, dst)
            print(f"  ✅ Copied {filename}")

    # Copy directories (excluding agents - we copy those selectively)
    # Note: work/ is handled separately to avoid copying history contents
    dirs_to_copy = [
        "commands",
        "templates",
        "config",
        "docs",
        "scripts",
    ]

    for dirname in dirs_to_copy:
        src_dir = garden_claude / dirname
        if src_dir.exists():
            dst_dir = project_claude / dirname
            shutil.copytree(src_dir, dst_dir, dirs_exist_ok=True)
            print(f"  ✅ Copied {dirname}/")

    # Create work/ directory structure without copying contents
    work_dst = project_claude / "work"
    work_subdirs = ["0_vision", "1_backlog", "2_active", "3_done", "history", "garden-requests"]
    for subdir in work_subdirs:
        (work_dst / subdir).mkdir(parents=True, exist_ok=True)
    print(f"  ✅ Created work/ (empty structure)")

    # Copy skills (core ones only)
    core_skills = [
        "catch-up",
        "handoff",
        "context-display",
        "agent-session-summary",
        "setup-validation",
        "hooks-setup",
        "workspace-setup",
        "version-notify",
        "garden-request",
    ]

    skills_src = garden_claude / "skills"
    skills_dst = project_claude / "skills"
    skills_dst.mkdir(exist_ok=True)

    for skill_name in core_skills:
        src_skill = skills_src / skill_name
        if src_skill.exists():
            dst_skill = skills_dst / skill_name
            shutil.copytree(src_skill, dst_skill, dirs_exist_ok=True)
            print(f"  ✅ Copied skill: {skill_name}")

    # Copy selected agents
    agents_dst = project_claude / "agents"
    agents_dst.mkdir(exist_ok=True)

    # Core agents location
    core_agents_src = garden_claude / "agents"
    # Specialized agents location
    marketplace_src = garden_root / "marketplace" / "agents"

    for agent_name in agents:
        # Try core agents first
        agent_file = f"{agent_name}.md"
        src_agent = core_agents_src / agent_file

        if not src_agent.exists():
            # Try marketplace/agents
            src_agent = marketplace_src / agent_file

        if src_agent.exists():
            dst_agent = agents_dst / agent_file
            shutil.copy2(src_agent, dst_agent)
            print(f"  ✅ Copied agent: {agent_name}")
        else:
            print(f"  ⚠️  Agent not found: {agent_name}")

    return project_claude


def merge_to_existing(garden_root, project_dir, agents):
    """Merge Garden structure to an existing repository without overwriting.

    Args:
        garden_root: Path to The Garden root directory
        project_dir: Path to the existing repository
        agents: List of agents to add

    Returns:
        Path to .claude directory
    """
    garden_claude = garden_root / ".claude"
    project_claude = project_dir / ".claude"

    # Track what was added vs skipped
    added = []
    skipped = []

    # Create .claude directory if it doesn't exist
    project_claude.mkdir(parents=True, exist_ok=True)

    # Copy core files (skip if exists)
    core_files = [
        "PROTOCOLS.md",
        "QUICKSTART.md",
        "PROJECT.md",
        "settings.json",
    ]

    for filename in core_files:
        src = garden_claude / filename
        dst = project_claude / filename
        if src.exists():
            if dst.exists():
                skipped.append(filename)
            else:
                shutil.copy2(src, dst)
                added.append(filename)

    if added:
        print(f"  ✅ Added core files: {', '.join(added)}")
    if skipped:
        print(f"  ⏭️  Skipped existing: {', '.join(skipped)}")

    # Merge directories (add new files, skip existing)
    dirs_to_merge = [
        "commands",
        "templates",
        "config",
        "docs",
        "scripts",
    ]

    for dirname in dirs_to_merge:
        src_dir = garden_claude / dirname
        dst_dir = project_claude / dirname
        if src_dir.exists():
            dst_dir.mkdir(exist_ok=True)
            dir_added = 0
            dir_skipped = 0
            for src_file in src_dir.rglob("*"):
                if src_file.is_file():
                    rel_path = src_file.relative_to(src_dir)
                    dst_file = dst_dir / rel_path
                    if dst_file.exists():
                        dir_skipped += 1
                    else:
                        dst_file.parent.mkdir(parents=True, exist_ok=True)
                        shutil.copy2(src_file, dst_file)
                        dir_added += 1
            if dir_added > 0 or dir_skipped > 0:
                print(f"  ✅ {dirname}/: +{dir_added} new, {dir_skipped} existing")

    # Create work/ directory structure if doesn't exist
    work_dst = project_claude / "work"
    work_subdirs = ["0_vision", "1_backlog", "2_active", "3_done", "history", "garden-requests"]
    work_created = False
    for subdir in work_subdirs:
        subdir_path = work_dst / subdir
        if not subdir_path.exists():
            subdir_path.mkdir(parents=True, exist_ok=True)
            work_created = True
    if work_created:
        print(f"  ✅ Created work/ structure")
    else:
        print(f"  ⏭️  work/ already exists")

    # Merge skills (add new, skip existing)
    core_skills = [
        "catch-up",
        "handoff",
        "context-display",
        "agent-session-summary",
        "setup-validation",
        "hooks-setup",
        "workspace-setup",
        "version-notify",
        "garden-request",
    ]

    skills_src = garden_claude / "skills"
    skills_dst = project_claude / "skills"
    skills_dst.mkdir(exist_ok=True)

    skills_added = []
    skills_skipped = []
    for skill_name in core_skills:
        src_skill = skills_src / skill_name
        dst_skill = skills_dst / skill_name
        if src_skill.exists():
            if dst_skill.exists():
                skills_skipped.append(skill_name)
            else:
                shutil.copytree(src_skill, dst_skill)
                skills_added.append(skill_name)

    if skills_added:
        print(f"  ✅ Added skills: {', '.join(skills_added)}")
    if skills_skipped:
        print(f"  ⏭️  Existing skills: {', '.join(skills_skipped)}")

    # Merge agents (add new, skip existing)
    agents_dst = project_claude / "agents"
    agents_dst.mkdir(exist_ok=True)

    core_agents_src = garden_claude / "agents"
    marketplace_src = garden_root / "marketplace" / "agents"

    agents_added = []
    agents_skipped = []
    agents_notfound = []

    for agent_name in agents:
        agent_file = f"{agent_name}.md"
        dst_agent = agents_dst / agent_file

        # Skip if agent already exists
        if dst_agent.exists():
            agents_skipped.append(agent_name)
            continue

        # Try core agents first, then marketplace
        src_agent = core_agents_src / agent_file
        if not src_agent.exists():
            src_agent = marketplace_src / agent_file

        if src_agent.exists():
            shutil.copy2(src_agent, dst_agent)
            agents_added.append(agent_name)
        else:
            agents_notfound.append(agent_name)

    if agents_added:
        print(f"  ✅ Added agents: {', '.join(agents_added)}")
    if agents_skipped:
        print(f"  ⏭️  Existing agents: {', '.join(agents_skipped)}")
    if agents_notfound:
        print(f"  ⚠️  Not found: {', '.join(agents_notfound)}")

    return project_claude


def copy_devcontainer(garden_root, project_dir, project_name, tech_stack, additional_services=None):
    """Copy and customize devcontainer template with dynamic port configuration.

    Args:
        garden_root: Path to The Garden root directory
        project_dir: Path to the new project directory
        project_name: Name of the project (for template substitution)
        tech_stack: Tech stack identifier (e.g., "python-flask")
        additional_services: Optional list of additional services detected from description

    Returns:
        Path to created .devcontainer directory, or None if no stack specified
    """
    if not tech_stack or tech_stack not in TECH_STACKS:
        return None

    stack_info = TECH_STACKS[tech_stack]
    template_dir = garden_root / ".claude" / "templates" / "devcontainer-templates" / stack_info["template_dir"]

    if not template_dir.exists():
        print(f"  Warning: Devcontainer template not found: {tech_stack}")
        return None

    # Combine base services with additional services
    all_services = list(stack_info.get("services", []))
    if additional_services:
        all_services.extend(additional_services)

    # Resolve any port conflicts
    all_services = resolve_port_conflicts(all_services)

    # Generate port configurations
    forward_ports = [s["port"] for s in all_services]
    ports_attributes = {
        str(s["port"]): {"label": s["label"], "onAutoForward": "notify"}
        for s in all_services
    }

    # Create .devcontainer directory in project
    devcontainer_dst = project_dir / ".devcontainer"
    devcontainer_dst.mkdir(parents=True, exist_ok=True)

    # Copy and customize each file
    for src_file in template_dir.iterdir():
        if src_file.is_file():
            dst_file = devcontainer_dst / src_file.name

            # Read content and substitute placeholders
            content = src_file.read_text()
            content = content.replace("{{PROJECT_NAME}}", title_case_name(project_name))
            content = content.replace("{{project_name}}", project_name)
            content = content.replace("{{FORWARD_PORTS}}", json.dumps(forward_ports))
            content = content.replace("{{PORTS_ATTRIBUTES}}", json.dumps(ports_attributes, indent=4))

            dst_file.write_text(content)

            # Make shell scripts executable
            if src_file.suffix == ".sh":
                dst_file.chmod(0o755)

            print(f"  Created .devcontainer/{src_file.name}")

    return devcontainer_dst


def add_to_existing(project_name, existing_path, project_type, agents, description):
    """
    Add Garden capabilities to an existing repository.

    Args:
        project_name: Name of the project (for CLAUDE.md)
        existing_path: Path to the existing repository
        project_type: Type of project (webapp, api, cli, etc.)
        agents: List of agents to include
        description: Project description

    Returns:
        Path to the repository, or None if error
    """
    garden_root = get_garden_root()
    project_dir = Path(existing_path).resolve()

    # Verify directory exists
    if not project_dir.exists():
        print(f"❌ Error: Directory does not exist: {project_dir}")
        print("   Did you mean to create a new project? Remove --existing flag.")
        return None

    if not project_dir.is_dir():
        print(f"❌ Error: Path is not a directory: {project_dir}")
        return None

    print(f"🌿 Adding Garden to existing repository: {project_dir}")

    # Handle existing CLAUDE.md
    claude_md_path = project_dir / "CLAUDE.md"
    if claude_md_path.exists():
        backup_path = project_dir / "CLAUDE.md.backup"
        shutil.copy2(claude_md_path, backup_path)
        print(f"  📋 Backed up existing CLAUDE.md to CLAUDE.md.backup")

    # Merge Garden structure
    print("\n📁 Merging Garden structure...")
    try:
        merge_to_existing(garden_root, project_dir, agents)
    except Exception as e:
        print(f"❌ Error merging structure: {e}")
        return None

    # Create lineage.json (always create/overwrite for sync capability)
    print("\n🔗 Creating Garden lineage...")
    try:
        create_lineage_json(project_dir, garden_root, project_name, project_type, description, agents, None)
        print("  ✅ Created lineage.json")
    except Exception as e:
        print(f"❌ Error creating lineage.json: {e}")
        return None

    # Initialize beads if not present
    print("\n🔮 Checking beads...")
    beads_path = project_dir / ".beads"
    if beads_path.exists():
        print("  ⏭️  .beads/ already exists")
    else:
        try:
            init_beads(project_dir, project_name)
            print("  ✅ Created .beads/ (execution state tracking)")
        except Exception as e:
            print(f"⚠️  Could not initialize beads: {e}")

    # Create CLAUDE.md (overwrites backup we just made)
    print("\n📝 Creating CLAUDE.md...")
    rooted_date = datetime.now().strftime("%Y-%m-%d")
    try:
        create_claude_md(project_dir, project_name, description, agents, rooted_date)
        print("  ✅ Created CLAUDE.md")
    except Exception as e:
        print(f"❌ Error creating CLAUDE.md: {e}")
        return None

    # Print summary
    print(f"\n🌿 Garden added to '{project_name}' successfully!")
    print(f"   Location: {project_dir}")
    print(f"   Type: {project_type}")
    print(f"   Agents: {', '.join(agents)}")
    print("\nNext steps:")
    print("1. Review the new CLAUDE.md (original backed up to CLAUDE.md.backup)")
    print("2. Run /onboard to complete setup")
    print("3. Use /sync-baseline to pull Garden updates")

    return project_dir


def init_project(project_name, path, project_type, agents, description, tech_stack=None):
    """
    Initialize a new project rooted from The Garden.

    Args:
        project_name: Name of the project (kebab-case)
        path: Path where the project should be created
        project_type: Type of project (webapp, api, cli, etc.)
        agents: List of agents to include
        description: Project description
        tech_stack: Optional tech stack for devcontainer (e.g., "python-flask")

    Returns:
        Path to created project directory, or None if error
    """
    garden_root = get_garden_root()
    project_dir = Path(path).resolve() / project_name

    # Check if directory already exists
    if project_dir.exists():
        print(f"❌ Error: Project directory already exists: {project_dir}")
        return None

    # Create project directory
    try:
        project_dir.mkdir(parents=True, exist_ok=False)
        print(f"✅ Created project directory: {project_dir}")
    except Exception as e:
        print(f"❌ Error creating directory: {e}")
        return None

    # Copy core structure
    print("\n📁 Copying Garden structure...")
    try:
        copy_core_structure(garden_root, project_dir, agents)
    except Exception as e:
        print(f"❌ Error copying structure: {e}")
        return None

    # Copy devcontainer if tech stack specified
    if tech_stack:
        print("\n🐳 Setting up devcontainer...")
        try:
            # Detect additional services from description
            additional_services = None
            if description:
                detected = detect_services(description)
                # Filter out services that match the primary stack
                primary_service = TECH_STACKS.get(tech_stack, {}).get("services", [{}])[0].get("name")
                additional_services = [s for s in detected if s["name"] != primary_service]
                if additional_services:
                    print(f"  Detected additional services: {', '.join(s['name'] for s in additional_services)}")

            copy_devcontainer(garden_root, project_dir, project_name, tech_stack, additional_services)
        except Exception as e:
            print(f"⚠️  Error setting up devcontainer: {e}")
            # Non-fatal - continue with project creation

    # Create lineage.json
    print("\n🔗 Creating Garden lineage...")
    try:
        create_lineage_json(project_dir, garden_root, project_name, project_type, description, agents, tech_stack)
        print("  ✅ Created lineage.json")
    except Exception as e:
        print(f"❌ Error creating lineage.json: {e}")
        return None

    # Initialize beads for execution state tracking
    print("\n🔮 Initializing beads...")
    try:
        init_beads(project_dir, project_name)
        print("  ✅ Created .beads/ (execution state tracking)")
    except Exception as e:
        print(f"⚠️  Could not initialize beads: {e}")
        # Non-fatal - continue with project creation

    # Create CLAUDE.md
    print("\n📝 Creating CLAUDE.md...")
    rooted_date = datetime.now().strftime("%Y-%m-%d")
    try:
        create_claude_md(project_dir, project_name, description, agents, rooted_date)
        print("  ✅ Created CLAUDE.md")
    except Exception as e:
        print(f"❌ Error creating CLAUDE.md: {e}")
        return None

    # Print summary
    print(f"\n🌱 Project '{project_name}' planted successfully!")
    print(f"   Location: {project_dir}")
    print(f"   Type: {project_type}")
    print(f"   Agents: {', '.join(agents)}")
    if tech_stack:
        stack_info = TECH_STACKS.get(tech_stack, {})
        print(f"   Devcontainer: {stack_info.get('display_name', tech_stack)}")
    print("\nNext steps:")
    print(f"1. cd {project_dir}")
    if tech_stack:
        print("2. Open in VS Code and reopen in container")
    else:
        print("2. Run Claude Code")
    print("3. Run /onboard to complete setup")
    print("4. Use /sync-baseline to pull Garden updates")

    return project_dir


def main():
    if len(sys.argv) < 2:
        print("Usage: init_project.py <project-name> --path <path> --type <type> [--agents a,b,c] [--description 'desc'] [--stack python-flask] [--existing]")
        print("\nArguments:")
        print("  project-name    Kebab-case project identifier")
        print("  --path          Directory where project will be created (or existing repo path)")
        print("  --type          Project type: webapp, api, cli, mobile, library, data, devops, business")
        print("  --agents        Comma-separated list of agents to include")
        print("  --description   Project description (quoted string)")
        print("  --stack         Tech stack for devcontainer: python-flask, python-fastapi, node-react, node-nextjs")
        print("  --existing      Add Garden to an existing repository (won't create new directory)")
        print("\nExamples:")
        print("  # New project")
        print("  init_project.py my-app --path ~/projects --type webapp --agents developer,architect --stack python-flask")
        print("  init_project.py api --path .. --type api --description 'REST API service' --stack python-fastapi")
        print("")
        print("  # Existing repository")
        print("  init_project.py my-repo --path ~/existing/repo --type webapp --agents developer,architect --existing")
        sys.exit(1)

    # Parse arguments
    project_name = sys.argv[1]

    args = {}
    flags = set()  # For boolean flags like --existing
    i = 2
    while i < len(sys.argv):
        if sys.argv[i].startswith("--"):
            key = sys.argv[i][2:]
            # Check if next arg is a value or another flag
            if i + 1 < len(sys.argv) and not sys.argv[i + 1].startswith("--"):
                args[key] = sys.argv[i + 1]
                i += 2
            else:
                # Boolean flag (no value)
                flags.add(key)
                i += 1
        else:
            i += 1

    is_existing = "existing" in flags
    path = args.get("path", "..")
    project_type = args.get("type", "general")
    agents_str = args.get("agents", "developer,architect")
    description = args.get("description", f"A {project_type} project rooted from The Garden")
    tech_stack = args.get("stack", None)

    # Validate tech stack if provided (not used for existing repos)
    if tech_stack and tech_stack not in TECH_STACKS:
        print(f"⚠️  Unknown tech stack '{tech_stack}'. Available: {', '.join(TECH_STACKS.keys())}")
        tech_stack = None

    agents = [a.strip() for a in agents_str.split(",")]

    if is_existing:
        # Existing repository mode
        print(f"🌿 Adding Garden to existing repo: {project_name}")
        print(f"   Path: {path}")
        print(f"   Type: {project_type}")
        print(f"   Agents: {', '.join(agents)}")
        print()

        result = add_to_existing(project_name, path, project_type, agents, description)
    else:
        # New project mode
        print(f"🌱 Planting project: {project_name}")
        print(f"   Path: {path}")
        print(f"   Type: {project_type}")
        print(f"   Agents: {', '.join(agents)}")
        if tech_stack:
            print(f"   Stack: {tech_stack}")
        print()

        result = init_project(project_name, path, project_type, agents, description, tech_stack)

    if result:
        sys.exit(0)
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
