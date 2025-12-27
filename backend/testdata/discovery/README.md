# Discovery Flow Test Fixtures

This directory contains JSON fixtures for testing the guided discovery flow without making actual Claude API calls.

## Fixtures Overview

Based on the bakery/cake order manager example from the UX design document.

### Stage 1: Welcome
- `welcome_response.json` - Root's initial welcome message

### Stage 2: Problem Discovery
- `problem_response.json` - Response after user describes their business
- `problem_followup_response.json` - Follow-up about pain points and goals

### Stage 3: User Personas
- `personas_response.json` - Initial question about users
- `personas_followup_response.json` - Follow-up about permissions

### Stage 4: MVP Scope
- `mvp_response.json` - Ask for essential features
- `mvp_followup_response.json` - Reflect back priorities
- `mvp_confirm_response.json` - Confirm MVP and future roadmap

### Stage 5: Summary
- `summary_response.json` - Show the summary card
- `summary_confirm_response.json` - Ask for confirmation

### Completion
- `complete_response.json` - Discovery complete, hand off to developer

## Fixture Format

```json
{
  "stage": "welcome",
  "response": "The actual message text...",
  "metadata": {
    "stage_complete": false,
    "next_stage": "problem",
    "extracted": {
      "business_context": "...",
      "problem_statement": "...",
      // ... more extracted data
    }
  },
  "summary_card": {
    // Only present in summary stage
    "project_name": "...",
    "solves_statement": "...",
    "users": [...],
    "mvp_features": [...],
    "future_features": [...]
  }
}
```

## Usage

### With MockClaudeService

```go
mock, err := service.NewMockClaudeService("backend/testdata/discovery")
if err != nil {
    log.Fatal(err)
}

// Send messages and get fixture responses
stream, err := mock.SendMessage(ctx, "discovery system prompt", messages)
```

### With MockDiscoveryRepository

```go
repo := repository.NewMockDiscoveryRepository()

// Create discovery state for a project
state, err := repo.Create(ctx, projectID)

// Update as discovery progresses
repo.SetBusinessContext(ctx, state.ID, "Custom cake bakery")
repo.UpdateStage(ctx, state.ID, repository.DiscoveryStageProblem)
```

## Example Conversation Flow

1. User starts new project
2. System returns `welcome_response.json`
3. User: "I run a small bakery that does custom cake orders"
4. System returns `problem_response.json`
5. User: "Tracking orders is chaos. We use paper and WhatsApp"
6. System returns `problem_followup_response.json`
7. ... continues through all stages
8. User confirms summary
9. System returns `complete_response.json`
