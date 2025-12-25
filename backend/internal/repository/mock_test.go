package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

func TestMockProjectRepository_Create(t *testing.T) {
	t.Run("creates project successfully", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()

		project, err := repo.Create(ctx, "Test Project")

		require.NoError(t, err)
		assert.Equal(t, "Test Project", project.Title)
		assert.NotEmpty(t, project.ID)
		assert.NotZero(t, project.CreatedAt)
		assert.NotZero(t, project.UpdatedAt)
	})
}

func TestMockProjectRepository_GetByID(t *testing.T) {
	t.Run("returns project when exists", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()
		created, _ := repo.Create(ctx, "Test Project")

		project, err := repo.GetByID(ctx, created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, project.ID)
		assert.Equal(t, "Test Project", project.Title)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()

		_, err := repo.GetByID(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockProjectRepository_List(t *testing.T) {
	t.Run("returns empty list when no projects", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()

		projects, err := repo.List(ctx)

		require.NoError(t, err)
		assert.Empty(t, projects)
	})

	t.Run("returns all projects", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()
		_, _ = repo.Create(ctx, "Project 1")
		_, _ = repo.Create(ctx, "Project 2")

		projects, err := repo.List(ctx)

		require.NoError(t, err)
		assert.Len(t, projects, 2)
	})
}

func TestMockProjectRepository_Delete(t *testing.T) {
	t.Run("deletes existing project", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()
		project, _ := repo.Create(ctx, "Test Project")

		err := repo.Delete(ctx, project.ID)

		require.NoError(t, err)
		_, err = repo.GetByID(ctx, project.ID)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()

		err := repo.Delete(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockProjectRepository_Messages(t *testing.T) {
	t.Run("creates and retrieves messages", func(t *testing.T) {
		repo := NewMockProjectRepository()
		ctx := context.Background()
		project, _ := repo.Create(ctx, "Test Project")

		msg1, err := repo.CreateMessage(ctx, project.ID, model.RoleUser, "Hello")
		require.NoError(t, err)
		assert.Equal(t, model.RoleUser, msg1.Role)
		assert.Equal(t, "Hello", msg1.Content)

		msg2, err := repo.CreateMessage(ctx, project.ID, model.RoleAssistant, "Hi!")
		require.NoError(t, err)

		messages, err := repo.GetMessages(ctx, project.ID)
		require.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, msg1.ID, messages[0].ID)
		assert.Equal(t, msg2.ID, messages[1].ID)
	})
}
