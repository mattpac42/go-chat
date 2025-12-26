package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

func TestMockDiscoveryRepository_Create(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, err := repo.Create(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to create discovery state: %v", err)
	}

	if state.ProjectID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, state.ProjectID)
	}
	if state.Stage != model.StageWelcome {
		t.Errorf("expected initial stage to be welcome, got %s", state.Stage)
	}
	if state.ID == uuid.Nil {
		t.Error("expected state ID to be generated")
	}

	// Creating again should return an error
	_, err = repo.Create(ctx, projectID)
	if err == nil {
		t.Error("expected error when creating discovery for same project")
	}
}

func TestMockDiscoveryRepository_GetByProjectID(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	// Get non-existent state
	_, err := repo.GetByProjectID(ctx, projectID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// Create and retrieve
	created, _ := repo.Create(ctx, projectID)
	retrieved, err := repo.GetByProjectID(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to get discovery state: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, retrieved.ID)
	}
}

func TestMockDiscoveryRepository_GetByID(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	created, _ := repo.Create(ctx, projectID)

	retrieved, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to get discovery state by ID: %v", err)
	}

	if retrieved.ProjectID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, retrieved.ProjectID)
	}
}

func TestMockDiscoveryRepository_UpdateStage(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)

	err := repo.UpdateStage(ctx, state.ID, model.StageProblem)
	if err != nil {
		t.Fatalf("failed to update stage: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if updated.Stage != model.StageProblem {
		t.Errorf("expected stage problem, got %s", updated.Stage)
	}
}

func TestMockDiscoveryRepository_SetBusinessContext(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)
	context := "Custom cake bakery"

	err := repo.SetBusinessContext(ctx, state.ID, context)
	if err != nil {
		t.Fatalf("failed to set business context: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if updated.BusinessContext != context {
		t.Errorf("expected business context '%s', got '%s'", context, updated.BusinessContext)
	}
}

func TestMockDiscoveryRepository_SetUsers(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)
	users := &DiscoveryUsers{
		Description:     "Owner and 2 employees",
		Count:           3,
		HasPermissions:  true,
		PermissionNotes: "Employees can manage orders only",
	}

	err := repo.SetUsers(ctx, state.ID, users)
	if err != nil {
		t.Fatalf("failed to set users: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if updated.Users == nil {
		t.Fatal("expected users to be set")
	}
	if updated.Users.Count != 3 {
		t.Errorf("expected user count 3, got %d", updated.Users.Count)
	}
	if !updated.Users.HasPermissions {
		t.Error("expected has_permissions to be true")
	}
}

func TestMockDiscoveryRepository_SetMVPFeatures(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)
	features := []DiscoveryMVPFeature{
		{Name: "Order list", Priority: 1},
		{Name: "Order form", Priority: 2},
		{Name: "Due dates", Priority: 3},
	}

	err := repo.SetMVPFeatures(ctx, state.ID, features)
	if err != nil {
		t.Fatalf("failed to set MVP features: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if len(updated.MVPFeatures) != 3 {
		t.Errorf("expected 3 MVP features, got %d", len(updated.MVPFeatures))
	}
	if updated.MVPFeatures[0].Name != "Order list" {
		t.Errorf("expected first feature 'Order list', got '%s'", updated.MVPFeatures[0].Name)
	}
}

func TestMockDiscoveryRepository_SetSummary(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)
	summary := &DiscoverySummary{
		ProjectName:     "Cake Order Manager",
		SolvesStatement: "Replaces paper and WhatsApp chaos",
	}

	err := repo.SetSummary(ctx, state.ID, summary)
	if err != nil {
		t.Fatalf("failed to set summary: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if updated.Summary == nil {
		t.Fatal("expected summary to be set")
	}
	if updated.Summary.ProjectName != "Cake Order Manager" {
		t.Errorf("expected project name 'Cake Order Manager', got '%s'", updated.Summary.ProjectName)
	}
}

func TestMockDiscoveryRepository_AddEdit(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)
	edit := DiscoveryEdit{
		Stage:         model.StageMVP,
		OriginalValue: "3 features",
		NewValue:      "4 features",
	}

	err := repo.AddEdit(ctx, state.ID, edit)
	if err != nil {
		t.Fatalf("failed to add edit: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if len(updated.EditHistory) != 1 {
		t.Errorf("expected 1 edit in history, got %d", len(updated.EditHistory))
	}
	if updated.EditHistory[0].NewValue != "4 features" {
		t.Errorf("expected new value '4 features', got '%s'", updated.EditHistory[0].NewValue)
	}
}

func TestMockDiscoveryRepository_MarkComplete(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)

	// Set summary first
	summary := &DiscoverySummary{
		ProjectName: "Test Project",
	}
	repo.SetSummary(ctx, state.ID, summary)

	err := repo.MarkComplete(ctx, state.ID)
	if err != nil {
		t.Fatalf("failed to mark complete: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if updated.Stage != model.StageComplete {
		t.Errorf("expected stage complete, got %s", updated.Stage)
	}
	if updated.Summary.ConfirmedAt.IsZero() {
		t.Error("expected ConfirmedAt to be set")
	}
}

func TestMockDiscoveryRepository_Delete(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)

	err := repo.Delete(ctx, state.ID)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	_, err = repo.GetByID(ctx, state.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}

	_, err = repo.GetByProjectID(ctx, projectID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for project after delete, got %v", err)
	}
}

func TestMockDiscoveryRepository_Reset(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()

	// Create multiple states
	repo.Create(ctx, uuid.New())
	repo.Create(ctx, uuid.New())

	if repo.Count() != 2 {
		t.Errorf("expected 2 states, got %d", repo.Count())
	}

	repo.Reset()

	if repo.Count() != 0 {
		t.Errorf("expected 0 states after reset, got %d", repo.Count())
	}
}

func TestMockDiscoveryRepository_DeepCopy(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)

	// Set some values
	repo.SetGoals(ctx, state.ID, []string{"Goal 1", "Goal 2"})

	// Get the state
	retrieved, _ := repo.GetByID(ctx, state.ID)

	// Modify the retrieved state's goals
	retrieved.Goals = append(retrieved.Goals, "Goal 3")

	// Get again and verify original is unchanged
	fresh, _ := repo.GetByID(ctx, state.ID)
	if len(fresh.Goals) != 2 {
		t.Errorf("expected 2 goals (deep copy protection), got %d", len(fresh.Goals))
	}
}

func TestMockDiscoveryRepository_GetAll(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()

	// Create multiple states
	repo.Create(ctx, uuid.New())
	repo.Create(ctx, uuid.New())
	repo.Create(ctx, uuid.New())

	all := repo.GetAll()
	if len(all) != 3 {
		t.Errorf("expected 3 states, got %d", len(all))
	}
}

func TestMockDiscoveryRepository_Update(t *testing.T) {
	repo := NewMockDiscoveryRepository()
	ctx := context.Background()
	projectID := uuid.New()

	state, _ := repo.Create(ctx, projectID)

	// Modify state
	state.BusinessContext = "Updated context"
	state.ProblemStatement = "Updated problem"
	state.Goals = []string{"New goal"}

	err := repo.Update(ctx, state)
	if err != nil {
		t.Fatalf("failed to update: %v", err)
	}

	updated, _ := repo.GetByID(ctx, state.ID)
	if updated.BusinessContext != "Updated context" {
		t.Errorf("expected business context 'Updated context', got '%s'", updated.BusinessContext)
	}
	if updated.ProblemStatement != "Updated problem" {
		t.Errorf("expected problem statement 'Updated problem', got '%s'", updated.ProblemStatement)
	}
	if len(updated.Goals) != 1 || updated.Goals[0] != "New goal" {
		t.Error("expected goals to be updated")
	}

	// UpdatedAt should be set
	if updated.UpdatedAt.Before(time.Now().Add(-time.Minute)) {
		t.Error("expected UpdatedAt to be recent")
	}
}
