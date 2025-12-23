# First Task

> Complete your first task with Garden.

## Example Task

Let's add a simple feature to see Garden in action.

### Step 1: Make a Request

Ask Claude:

```
"Create a utility function that formats dates in ISO format"
```

### Step 2: Watch Delegation

Garden will:

1. **Understand** - Parse your request
2. **Delegate** - Invoke the developer agent
3. **Implement** - Agent writes the code
4. **Return** - Show you the results

You'll see something like:

```
I'll have the developer agent handle this.

[Developer agent creates the function]

Here's the utility function:
[code output]

Context: 🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 30%
```

### Step 3: Review Results

The developer agent will:
- Write the function with tests (TDD)
- Document the code
- Create a session history file

### Step 4: Continue Working

Ask for more:
- "Add error handling for invalid dates"
- "Create a companion function for parsing"
- "Write tests for edge cases"

Each request goes through the same delegation pattern.

## More Complex Example

For multi-domain tasks, multiple agents work together:

```
"Design and implement user authentication"
```

Garden will:
1. **architect** → Design the auth system
2. **product** → Define requirements
3. **developer** → Implement the code
4. **platform** → Set up infrastructure

Agents may run in parallel for independent work.

## Track Progress

Create a task file for larger work:

```
/Users/project/.garden/work/active/add-auth.md
```

Garden tracks subtasks and progress automatically.

## Session Transition

When context reaches 75%, you'll see:

```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟨🟨🟨⬛⬛⬛⬛⬛ 75%
🔄 Session handoff recommended
```

Run `/handoff` to create transition files, then start a new session.

## Tips

1. **Be specific** - Clear requests get better results
2. **Trust delegation** - Let agents do their work
3. **Watch context** - Know when to hand off
4. **Review history** - Check `.garden/work/history/` for agent work

## Next Steps

- Read [Working with Agents](../guides/working-with-agents.md) for advanced patterns
- Explore [Context Management](../guides/context-management.md) for session handling
- Check [Agents Reference](../reference/agents.md) for full capabilities
