# Garden Request: GR-003

## Title
Dynamic port configuration for devcontainer templates

## Priority
Medium

## Category
Enhancement

## Source Project
Go Chat

## Date
2026-01-01

---

## Summary

Devcontainer templates currently hardcode single service ports (e.g., 3000 for Next.js, 8000 for FastAPI). Fullstack projects with separate frontend and backend services require manual editing to add additional ports. This enhancement proposes adding template placeholders for dynamic port configuration during project creation.

## Current Behavior

Each template has a fixed `forwardPorts` configuration:

```json
// node-nextjs template
"forwardPorts": [3000],
"portsAttributes": {
  "3000": {
    "label": "Next.js App",
    "onAutoForward": "notify"
  }
}
```

When creating a fullstack project (e.g., Next.js + Go API), developers must manually edit the devcontainer.json to add backend ports.

## Proposed Behavior

Add template placeholders that get populated during project creation:

```json
"forwardPorts": [{{FRONTEND_PORT}}{{#BACKEND_PORT}}, {{BACKEND_PORT}}{{/BACKEND_PORT}}],

"portsAttributes": {
  "{{FRONTEND_PORT}}": {
    "label": "{{FRONTEND_LABEL}}",
    "onAutoForward": "notify"
  }{{#BACKEND_PORT}},
  "{{BACKEND_PORT}}": {
    "label": "{{BACKEND_LABEL}}",
    "onAutoForward": "notify"
  }{{/BACKEND_PORT}}
}
```

During `/plant` or project creation:
1. Agent infers project type from user description (e.g., "Next.js + Go API")
2. Agent automatically selects conventional ports based on detected stack
3. Only prompt user if there's a port conflict or non-standard requirement

**Port inference logic:**
- Frontend detected → use framework default (Next.js=3000, Vite=5173)
- Backend detected → use framework default (Go=8080, FastAPI=8000, Express=3001)
- Conflict (e.g., two Node services) → ask user or auto-increment (3000, 3001)

## Implementation Notes

Two approaches:

**Option A: Mustache-style conditionals** (shown above)
- More flexible, supports optional sections
- Requires template engine that handles conditionals

**Option B: Separate fullstack templates**
- Create `fullstack-nextjs-go`, `fullstack-nextjs-python`, etc.
- Simpler implementation but more templates to maintain

**Option C: Post-processing script**
- Keep simple placeholders: `{{PORTS_CONFIG}}`
- Generate the JSON block in setup.sh based on project type

Recommend Option A for flexibility without template explosion.

## Files to Modify

| File | Change |
|------|--------|
| `.claude/templates/devcontainer-templates/*/devcontainer.json` | Add port placeholders |
| `.claude/skills/plant-project/README.md` | Add fullstack prompts to workflow |
| `.claude/scripts/plant.py` (or equivalent) | Handle port placeholder substitution |

## Testing

1. Run `/plant` with description like "Next.js frontend with Go API backend"
2. Verify agent infers ports automatically (3000 + 8080)
3. Check generated `devcontainer.json` has both ports configured
4. Build devcontainer and verify both ports forward correctly
5. Test conflict case: "Express frontend + Express API" - verify agent handles gracefully

## Notes

- Backwards compatible: existing single-service projects continue working
- Agent-driven: ports are inferred from tech stack, not asked of user
- User only prompted when ambiguity exists (rare)
- Standard port conventions make this reliable for 95% of cases
