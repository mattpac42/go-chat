package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

func TestMockDiscoveryRepository_Create(t *testing.T) {
	t.Run("creates discovery successfully", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()

		discovery, err := repo.Create(ctx, projectID)

		require.NoError(t, err)
		assert.NotEmpty(t, discovery.ID)
		assert.Equal(t, projectID, discovery.ProjectID)
		assert.Equal(t, model.StageWelcome, discovery.Stage)
		assert.NotZero(t, discovery.StageStartedAt)
		assert.NotZero(t, discovery.CreatedAt)
		assert.NotZero(t, discovery.UpdatedAt)
	})

	t.Run("fails when discovery already exists for project", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()

		_, err := repo.Create(ctx, projectID)
		require.NoError(t, err)

		_, err = repo.Create(ctx, projectID)
		assert.Error(t, err)
	})
}

func TestMockDiscoveryRepository_GetByProjectID(t *testing.T) {
	t.Run("returns discovery when exists", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		created, _ := repo.Create(ctx, projectID)

		discovery, err := repo.GetByProjectID(ctx, projectID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, discovery.ID)
		assert.Equal(t, projectID, discovery.ProjectID)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		_, err := repo.GetByProjectID(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_GetByID(t *testing.T) {
	t.Run("returns discovery when exists", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		created, _ := repo.Create(ctx, projectID)

		discovery, err := repo.GetByID(ctx, created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, discovery.ID)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		_, err := repo.GetByID(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_Update(t *testing.T) {
	t.Run("updates discovery fields", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		created, _ := repo.Create(ctx, projectID)

		businessContext := "I run a bakery"
		problemStatement := "Order tracking is chaos"
		projectName := "Cake Order Manager"

		created.BusinessContext = &businessContext
		created.ProblemStatement = &problemStatement
		created.ProjectName = &projectName
		_ = created.SetGoals([]string{"track orders", "reduce errors"})

		updated, err := repo.Update(ctx, created)

		require.NoError(t, err)
		assert.Equal(t, &businessContext, updated.BusinessContext)
		assert.Equal(t, &problemStatement, updated.ProblemStatement)
		assert.Equal(t, &projectName, updated.ProjectName)

		goals, _ := updated.Goals()
		assert.Contains(t, goals, "track orders")
		assert.Contains(t, goals, "reduce errors")
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		discovery := &model.ProjectDiscovery{ID: uuid.New()}
		_, err := repo.Update(ctx, discovery)

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_UpdateStage(t *testing.T) {
	t.Run("advances stage successfully", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		created, _ := repo.Create(ctx, projectID)

		updated, err := repo.UpdateStage(ctx, created.ID, model.StageProblem)

		require.NoError(t, err)
		assert.Equal(t, model.StageProblem, updated.Stage)
		assert.True(t, updated.StageStartedAt.After(created.StageStartedAt) || updated.StageStartedAt.Equal(created.StageStartedAt))
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		_, err := repo.UpdateStage(ctx, uuid.New(), model.StageProblem)

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_MarkComplete(t *testing.T) {
	t.Run("marks discovery as complete", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		created, _ := repo.Create(ctx, projectID)

		completed, err := repo.MarkComplete(ctx, created.ID)

		require.NoError(t, err)
		assert.Equal(t, model.StageComplete, completed.Stage)
		assert.NotNil(t, completed.ConfirmedAt)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		_, err := repo.MarkComplete(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_Delete(t *testing.T) {
	t.Run("deletes discovery and associated data", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		created, _ := repo.Create(ctx, projectID)

		// Add some related data
		_, _ = repo.AddUser(ctx, &model.DiscoveryUser{
			DiscoveryID: created.ID,
			Description: "Test user",
			UserCount:   1,
		})
		_, _ = repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: created.ID,
			Name:        "Test feature",
			Priority:    1,
		})

		err := repo.Delete(ctx, created.ID)

		require.NoError(t, err)
		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		err := repo.Delete(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_Users(t *testing.T) {
	t.Run("adds and retrieves users", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		notes := "Full access"
		user, err := repo.AddUser(ctx, &model.DiscoveryUser{
			DiscoveryID:     discovery.ID,
			Description:     "Owner/Baker",
			UserCount:       1,
			HasPermissions:  true,
			PermissionNotes: &notes,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, "Owner/Baker", user.Description)
		assert.Equal(t, 1, user.UserCount)
		assert.True(t, user.HasPermissions)
		assert.Equal(t, &notes, user.PermissionNotes)

		users, err := repo.GetUsers(ctx, discovery.ID)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, user.ID, users[0].ID)
	})

	t.Run("returns empty list when no users", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		users, err := repo.GetUsers(ctx, discovery.ID)

		require.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("updates user", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		user, _ := repo.AddUser(ctx, &model.DiscoveryUser{
			DiscoveryID: discovery.ID,
			Description: "Original",
			UserCount:   1,
		})

		user.Description = "Updated"
		user.UserCount = 5
		updated, err := repo.UpdateUser(ctx, user)

		require.NoError(t, err)
		assert.Equal(t, "Updated", updated.Description)
		assert.Equal(t, 5, updated.UserCount)
	})

	t.Run("deletes user", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		user, _ := repo.AddUser(ctx, &model.DiscoveryUser{
			DiscoveryID: discovery.ID,
			Description: "To be deleted",
			UserCount:   1,
		})

		err := repo.DeleteUser(ctx, user.ID)

		require.NoError(t, err)
		users, _ := repo.GetUsers(ctx, discovery.ID)
		assert.Empty(t, users)
	})

	t.Run("delete returns error when user not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		err := repo.DeleteUser(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_Features(t *testing.T) {
	t.Run("adds and retrieves features", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		feature, err := repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Order List",
			Priority:    1,
			Version:     "v1",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, feature.ID)
		assert.Equal(t, "Order List", feature.Name)
		assert.Equal(t, 1, feature.Priority)
		assert.Equal(t, "v1", feature.Version)

		features, err := repo.GetFeatures(ctx, discovery.ID)
		require.NoError(t, err)
		assert.Len(t, features, 1)
	})

	t.Run("defaults version to v1 when empty", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		feature, err := repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Order List",
			Priority:    1,
			Version:     "", // Empty version
		})

		require.NoError(t, err)
		assert.Equal(t, "v1", feature.Version)
	})

	t.Run("separates MVP and future features", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		// Add MVP features
		_, _ = repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Order List",
			Priority:    1,
			Version:     "v1",
		})
		_, _ = repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Order Form",
			Priority:    2,
			Version:     "v1",
		})

		// Add future feature
		_, _ = repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Calendar View",
			Priority:    1,
			Version:     "v2",
		})

		mvpFeatures, err := repo.GetMVPFeatures(ctx, discovery.ID)
		require.NoError(t, err)
		assert.Len(t, mvpFeatures, 2)

		futureFeatures, err := repo.GetFutureFeatures(ctx, discovery.ID)
		require.NoError(t, err)
		assert.Len(t, futureFeatures, 1)
		assert.Equal(t, "Calendar View", futureFeatures[0].Name)
	})

	t.Run("updates feature", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		feature, _ := repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Original",
			Priority:    1,
			Version:     "v1",
		})

		feature.Name = "Updated"
		feature.Version = "v2"
		updated, err := repo.UpdateFeature(ctx, feature)

		require.NoError(t, err)
		assert.Equal(t, "Updated", updated.Name)
		assert.Equal(t, "v2", updated.Version)
	})

	t.Run("deletes feature", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		feature, _ := repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "To be deleted",
			Priority:    1,
		})

		err := repo.DeleteFeature(ctx, feature.ID)

		require.NoError(t, err)
		features, _ := repo.GetFeatures(ctx, discovery.ID)
		assert.Empty(t, features)
	})

	t.Run("delete returns error when feature not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		err := repo.DeleteFeature(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockDiscoveryRepository_EditHistory(t *testing.T) {
	t.Run("adds and retrieves edit history", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		history, err := repo.AddEditHistory(ctx, &model.DiscoveryEditHistory{
			DiscoveryID:   discovery.ID,
			Stage:         "personas",
			FieldEdited:   "users",
			OriginalValue: "2 employees",
			NewValue:      "3 employees",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, history.ID)
		assert.Equal(t, "personas", history.Stage)
		assert.Equal(t, "users", history.FieldEdited)
		assert.NotZero(t, history.EditedAt)

		histories, err := repo.GetEditHistory(ctx, discovery.ID)
		require.NoError(t, err)
		assert.Len(t, histories, 1)
	})

	t.Run("returns empty list when no history", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		histories, err := repo.GetEditHistory(ctx, discovery.ID)

		require.NoError(t, err)
		assert.Empty(t, histories)
	})
}

func TestMockDiscoveryRepository_GetSummary(t *testing.T) {
	t.Run("returns complete summary", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		// Set discovery data
		projectName := "Cake Order Manager"
		solvesStatement := "Replaces paper and WhatsApp chaos"
		discovery.ProjectName = &projectName
		discovery.SolvesStatement = &solvesStatement
		_, _ = repo.Update(ctx, discovery)

		// Add users
		_, _ = repo.AddUser(ctx, &model.DiscoveryUser{
			DiscoveryID:    discovery.ID,
			Description:    "Owner/Baker",
			UserCount:      1,
			HasPermissions: true,
		})
		_, _ = repo.AddUser(ctx, &model.DiscoveryUser{
			DiscoveryID:    discovery.ID,
			Description:    "Order Takers",
			UserCount:      2,
			HasPermissions: false,
		})

		// Add MVP features
		_, _ = repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Order List",
			Priority:    1,
			Version:     "v1",
		})
		_, _ = repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Order Form",
			Priority:    2,
			Version:     "v1",
		})

		// Add future features
		_, _ = repo.AddFeature(ctx, &model.DiscoveryFeature{
			DiscoveryID: discovery.ID,
			Name:        "Calendar View",
			Priority:    1,
			Version:     "v2",
		})

		summary, err := repo.GetSummary(ctx, discovery.ID)

		require.NoError(t, err)
		assert.Equal(t, "Cake Order Manager", summary.ProjectName)
		assert.Equal(t, "Replaces paper and WhatsApp chaos", summary.SolvesStatement)
		assert.Len(t, summary.Users, 2)
		assert.Len(t, summary.MVPFeatures, 2)
		assert.Len(t, summary.FutureFeatures, 1)
	})

	t.Run("returns error when discovery not found", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()

		_, err := repo.GetSummary(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns empty arrays when no users or features", func(t *testing.T) {
		repo := NewMockDiscoveryRepository()
		ctx := context.Background()
		projectID := uuid.New()
		discovery, _ := repo.Create(ctx, projectID)

		summary, err := repo.GetSummary(ctx, discovery.ID)

		require.NoError(t, err)
		assert.Empty(t, summary.ProjectName)
		assert.Empty(t, summary.SolvesStatement)
		assert.Empty(t, summary.Users)
		assert.Empty(t, summary.MVPFeatures)
		assert.Empty(t, summary.FutureFeatures)
	})
}

func TestDiscoveryStage(t *testing.T) {
	t.Run("NextStage returns correct next stage", func(t *testing.T) {
		tests := []struct {
			current  model.DiscoveryStage
			expected model.DiscoveryStage
		}{
			{model.StageWelcome, model.StageProblem},
			{model.StageProblem, model.StagePersonas},
			{model.StagePersonas, model.StageMVP},
			{model.StageMVP, model.StageSummary},
			{model.StageSummary, model.StageComplete},
			{model.StageComplete, ""},
		}

		for _, tt := range tests {
			t.Run(string(tt.current), func(t *testing.T) {
				next := tt.current.NextStage()
				assert.Equal(t, tt.expected, next)
			})
		}
	})

	t.Run("StageNumber returns correct number", func(t *testing.T) {
		tests := []struct {
			stage    model.DiscoveryStage
			expected int
		}{
			{model.StageWelcome, 1},
			{model.StageProblem, 2},
			{model.StagePersonas, 3},
			{model.StageMVP, 4},
			{model.StageSummary, 5},
			{model.StageComplete, 6},
		}

		for _, tt := range tests {
			t.Run(string(tt.stage), func(t *testing.T) {
				number := tt.stage.StageNumber()
				assert.Equal(t, tt.expected, number)
			})
		}
	})

	t.Run("IsComplete returns correct value", func(t *testing.T) {
		assert.False(t, model.StageWelcome.IsComplete())
		assert.False(t, model.StageSummary.IsComplete())
		assert.True(t, model.StageComplete.IsComplete())
	})

	t.Run("IsValidStage validates stages", func(t *testing.T) {
		assert.True(t, model.IsValidStage("welcome"))
		assert.True(t, model.IsValidStage("problem"))
		assert.True(t, model.IsValidStage("personas"))
		assert.True(t, model.IsValidStage("mvp"))
		assert.True(t, model.IsValidStage("summary"))
		assert.True(t, model.IsValidStage("complete"))
		assert.False(t, model.IsValidStage("invalid"))
		assert.False(t, model.IsValidStage(""))
	})
}

func TestProjectDiscovery_Goals(t *testing.T) {
	t.Run("SetGoals and Goals work correctly", func(t *testing.T) {
		discovery := &model.ProjectDiscovery{}

		err := discovery.SetGoals([]string{"goal1", "goal2"})
		require.NoError(t, err)

		goals, err := discovery.Goals()
		require.NoError(t, err)
		assert.Equal(t, []string{"goal1", "goal2"}, goals)
	})

	t.Run("Goals returns empty slice when nil", func(t *testing.T) {
		discovery := &model.ProjectDiscovery{}

		goals, err := discovery.Goals()

		require.NoError(t, err)
		assert.Empty(t, goals)
	})

	t.Run("SetGoals with nil sets empty array", func(t *testing.T) {
		discovery := &model.ProjectDiscovery{}

		err := discovery.SetGoals(nil)
		require.NoError(t, err)

		goals, err := discovery.Goals()
		require.NoError(t, err)
		assert.Empty(t, goals)
	})
}

func TestProjectDiscovery_ToResponse(t *testing.T) {
	t.Run("converts to response correctly", func(t *testing.T) {
		businessContext := "bakery"
		projectName := "Cake Order Manager"

		discovery := &model.ProjectDiscovery{
			ID:              uuid.New(),
			ProjectID:       uuid.New(),
			Stage:           model.StageProblem,
			BusinessContext: &businessContext,
			ProjectName:     &projectName,
		}
		_ = discovery.SetGoals([]string{"goal1"})

		response, err := discovery.ToResponse()

		require.NoError(t, err)
		assert.Equal(t, discovery.ID, response.ID)
		assert.Equal(t, model.StageProblem, response.Stage)
		assert.Equal(t, 2, response.StageNumber)
		assert.Equal(t, 5, response.TotalStages) // Excludes 'complete' from visible stages
		assert.Equal(t, &businessContext, response.BusinessContext)
		assert.Equal(t, &projectName, response.ProjectName)
		assert.Equal(t, []string{"goal1"}, response.Goals)
	})
}

func TestDiscoveryFeature_IsMVP(t *testing.T) {
	t.Run("returns true for v1", func(t *testing.T) {
		feature := &model.DiscoveryFeature{Version: "v1"}
		assert.True(t, feature.IsMVP())
	})

	t.Run("returns false for other versions", func(t *testing.T) {
		feature := &model.DiscoveryFeature{Version: "v2"}
		assert.False(t, feature.IsMVP())

		feature.Version = "v3"
		assert.False(t, feature.IsMVP())
	})
}
