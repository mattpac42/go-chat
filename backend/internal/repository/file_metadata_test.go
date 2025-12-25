package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
)

func TestMockFileMetadataRepository_Create(t *testing.T) {
	t.Run("creates metadata successfully", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()

		meta, err := repo.Create(ctx, fileID, "Short desc", "Long description", "Backend Services")

		require.NoError(t, err)
		assert.NotEmpty(t, meta.ID)
		assert.Equal(t, fileID, meta.FileID)
		assert.Equal(t, "Short desc", meta.ShortDescription)
		assert.Equal(t, "Long description", meta.LongDescription)
		assert.Equal(t, "Backend Services", meta.FunctionalGroup)
		assert.NotZero(t, meta.CreatedAt)
		assert.NotZero(t, meta.UpdatedAt)
	})

	t.Run("fails when metadata already exists for file", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()

		_, err := repo.Create(ctx, fileID, "First", "First desc", "Group1")
		require.NoError(t, err)

		_, err = repo.Create(ctx, fileID, "Second", "Second desc", "Group2")
		assert.Error(t, err)
	})
}

func TestMockFileMetadataRepository_GetByFileID(t *testing.T) {
	t.Run("returns metadata when exists", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()
		created, _ := repo.Create(ctx, fileID, "Short", "Long", "Group")

		meta, err := repo.GetByFileID(ctx, fileID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, meta.ID)
		assert.Equal(t, fileID, meta.FileID)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()

		_, err := repo.GetByFileID(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockFileMetadataRepository_GetByID(t *testing.T) {
	t.Run("returns metadata when exists", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()
		created, _ := repo.Create(ctx, fileID, "Short", "Long", "Group")

		meta, err := repo.GetByID(ctx, created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, meta.ID)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()

		_, err := repo.GetByID(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockFileMetadataRepository_Update(t *testing.T) {
	t.Run("updates all fields when provided", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()
		created, _ := repo.Create(ctx, fileID, "Short", "Long", "Group")

		newShort := "New Short"
		newLong := "New Long"
		newGroup := "New Group"
		meta, err := repo.Update(ctx, created.ID, &newShort, &newLong, &newGroup)

		require.NoError(t, err)
		assert.Equal(t, "New Short", meta.ShortDescription)
		assert.Equal(t, "New Long", meta.LongDescription)
		assert.Equal(t, "New Group", meta.FunctionalGroup)
	})

	t.Run("updates only specified fields", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()
		created, _ := repo.Create(ctx, fileID, "Short", "Long", "Group")

		newShort := "New Short"
		meta, err := repo.Update(ctx, created.ID, &newShort, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, "New Short", meta.ShortDescription)
		assert.Equal(t, "Long", meta.LongDescription)
		assert.Equal(t, "Group", meta.FunctionalGroup)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()

		newShort := "New Short"
		_, err := repo.Update(ctx, uuid.New(), &newShort, nil, nil)

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockFileMetadataRepository_Upsert(t *testing.T) {
	t.Run("creates new metadata when none exists", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()

		meta, err := repo.Upsert(ctx, fileID, "Short", "Long", "Group")

		require.NoError(t, err)
		assert.NotEmpty(t, meta.ID)
		assert.Equal(t, fileID, meta.FileID)
		assert.Equal(t, "Short", meta.ShortDescription)
	})

	t.Run("updates existing metadata", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()

		first, _ := repo.Upsert(ctx, fileID, "First", "First Long", "Group1")
		second, err := repo.Upsert(ctx, fileID, "Second", "Second Long", "Group2")

		require.NoError(t, err)
		assert.Equal(t, first.ID, second.ID) // Same metadata record
		assert.Equal(t, "Second", second.ShortDescription)
		assert.Equal(t, "Second Long", second.LongDescription)
		assert.Equal(t, "Group2", second.FunctionalGroup)
	})
}

func TestMockFileMetadataRepository_Delete(t *testing.T) {
	t.Run("deletes existing metadata", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()
		created, _ := repo.Create(ctx, fileID, "Short", "Long", "Group")

		err := repo.Delete(ctx, created.ID)

		require.NoError(t, err)
		_, err = repo.GetByID(ctx, created.ID)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()

		err := repo.Delete(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockFileMetadataRepository_DeleteByFileID(t *testing.T) {
	t.Run("deletes metadata by file ID", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		fileID := uuid.New()
		_, _ = repo.Create(ctx, fileID, "Short", "Long", "Group")

		err := repo.DeleteByFileID(ctx, fileID)

		require.NoError(t, err)
		_, err = repo.GetByFileID(ctx, fileID)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()

		err := repo.DeleteByFileID(ctx, uuid.New())

		assert.Equal(t, ErrNotFound, err)
	})
}

func TestMockFileMetadataRepository_GetFilesWithMetadata(t *testing.T) {
	t.Run("returns files with their metadata", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		projectID := uuid.New()

		// Add files
		file1 := &model.File{
			ID:        uuid.New(),
			ProjectID: projectID,
			Path:      "src/main.go",
			Filename:  "main.go",
			Language:  "go",
			CreatedAt: time.Now().UTC(),
		}
		file2 := &model.File{
			ID:        uuid.New(),
			ProjectID: projectID,
			Path:      "src/handler.go",
			Filename:  "handler.go",
			Language:  "go",
			CreatedAt: time.Now().UTC(),
		}
		repo.AddFile(file1)
		repo.AddFile(file2)

		// Add metadata for file1 only
		_, _ = repo.Create(ctx, file1.ID, "Entry point", "Main application entry", "Backend Services")

		files, err := repo.GetFilesWithMetadata(ctx, projectID)

		require.NoError(t, err)
		assert.Len(t, files, 2)

		// Find file1 in results
		var foundFile1 bool
		for _, f := range files {
			if f.ID == file1.ID {
				foundFile1 = true
				assert.Equal(t, "Entry point", f.ShortDescription)
				assert.Equal(t, "Backend Services", f.FunctionalGroup)
			}
		}
		assert.True(t, foundFile1)
	})

	t.Run("returns empty list when no files", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()

		files, err := repo.GetFilesWithMetadata(ctx, uuid.New())

		require.NoError(t, err)
		assert.Empty(t, files)
	})
}

func TestMockFileMetadataRepository_GetFilesByFunctionalGroup(t *testing.T) {
	t.Run("returns only files in specified group", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		projectID := uuid.New()

		// Add files
		file1 := &model.File{
			ID:        uuid.New(),
			ProjectID: projectID,
			Path:      "src/main.go",
			Filename:  "main.go",
			Language:  "go",
			CreatedAt: time.Now().UTC(),
		}
		file2 := &model.File{
			ID:        uuid.New(),
			ProjectID: projectID,
			Path:      "src/handler.go",
			Filename:  "handler.go",
			Language:  "go",
			CreatedAt: time.Now().UTC(),
		}
		repo.AddFile(file1)
		repo.AddFile(file2)

		// Add metadata with different groups
		_, _ = repo.Create(ctx, file1.ID, "Entry point", "Main application entry", "Backend Services")
		_, _ = repo.Create(ctx, file2.ID, "Handler", "HTTP handler", "API Layer")

		files, err := repo.GetFilesByFunctionalGroup(ctx, projectID, "Backend Services")

		require.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, file1.ID, files[0].ID)
		assert.Equal(t, "Backend Services", files[0].FunctionalGroup)
	})
}

func TestMockFileMetadataRepository_GetFunctionalGroups(t *testing.T) {
	t.Run("returns groups with file counts", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		projectID := uuid.New()

		// Add files
		file1 := &model.File{ID: uuid.New(), ProjectID: projectID, Path: "a.go", Filename: "a.go", CreatedAt: time.Now().UTC()}
		file2 := &model.File{ID: uuid.New(), ProjectID: projectID, Path: "b.go", Filename: "b.go", CreatedAt: time.Now().UTC()}
		file3 := &model.File{ID: uuid.New(), ProjectID: projectID, Path: "c.go", Filename: "c.go", CreatedAt: time.Now().UTC()}
		repo.AddFile(file1)
		repo.AddFile(file2)
		repo.AddFile(file3)

		// Add metadata
		_, _ = repo.Create(ctx, file1.ID, "Short1", "Long1", "Backend Services")
		_, _ = repo.Create(ctx, file2.ID, "Short2", "Long2", "Backend Services")
		_, _ = repo.Create(ctx, file3.ID, "Short3", "Long3", "Frontend")

		groups, err := repo.GetFunctionalGroups(ctx, projectID)

		require.NoError(t, err)
		assert.Len(t, groups, 2)

		// Find Backend Services group
		var backendCount, frontendCount int
		for _, g := range groups {
			if g.Name == "Backend Services" {
				backendCount = g.FileCount
			}
			if g.Name == "Frontend" {
				frontendCount = g.FileCount
			}
		}
		assert.Equal(t, 2, backendCount)
		assert.Equal(t, 1, frontendCount)
	})

	t.Run("excludes empty functional groups", func(t *testing.T) {
		repo := NewMockFileMetadataRepository()
		ctx := context.Background()
		projectID := uuid.New()

		file1 := &model.File{ID: uuid.New(), ProjectID: projectID, Path: "a.go", Filename: "a.go", CreatedAt: time.Now().UTC()}
		repo.AddFile(file1)
		_, _ = repo.Create(ctx, file1.ID, "Short", "Long", "") // Empty group

		groups, err := repo.GetFunctionalGroups(ctx, projectID)

		require.NoError(t, err)
		assert.Empty(t, groups)
	})
}
