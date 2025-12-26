package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/model"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
)

// Nudge represents a contextual suggestion
type Nudge struct {
	Type    model.NudgeType `json:"type"`
	Title   string          `json:"title"`
	Message string          `json:"message"`
	Action  string          `json:"action"`
	Icon    string          `json:"icon"`
}

// NudgeService manages contextual suggestions
type NudgeService struct {
	repo           repository.AchievementRepository
	achievementSvc *AchievementService
	logger         zerolog.Logger
}

// NewNudgeService creates a new NudgeService
func NewNudgeService(repo repository.AchievementRepository, achievementSvc *AchievementService, logger zerolog.Logger) *NudgeService {
	return &NudgeService{
		repo:           repo,
		achievementSvc: achievementSvc,
		logger:         logger,
	}
}

// GetNextNudge returns the most appropriate nudge for current context
func (s *NudgeService) GetNextNudge(ctx context.Context, projectID uuid.UUID) (*Nudge, error) {
	progress, err := s.achievementSvc.GetProgress(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Nudge selection logic based on progress

	// Haven't viewed any code yet
	if progress.CodeViewsCount == 0 {
		if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeExploreCode); !hasSeenNudge {
			return &Nudge{
				Type:    model.NudgeExploreCode,
				Title:   "Peek Behind the Curtain",
				Message: "Curious what makes your app tick? Tap any file to see how it works.",
				Action:  "Show me",
				Icon:    "eye",
			}, nil
		}
	}

	// Viewed code but never expanded tree
	if progress.CodeViewsCount >= 3 && progress.TreeExpansionsCount == 0 {
		if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeTryTreeView); !hasSeenNudge {
			return &Nudge{
				Type:    model.NudgeTryTreeView,
				Title:   "See the Big Picture",
				Message: "Want to see how your files connect? Try the tree view!",
				Action:  "Show connections",
				Icon:    "git-branch",
			}, nil
		}
	}

	// At level 1 for a while, suggest level up
	if progress.CurrentLevel == model.LevelFunctional && progress.FilesViewedCount >= 5 {
		if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeLevelUp); !hasSeenNudge {
			return &Nudge{
				Type:    model.NudgeLevelUp,
				Title:   "Ready for More?",
				Message: "You're getting the hang of it! Try Level 2 to see more details.",
				Action:  "Level up",
				Icon:    "trending-up",
			}, nil
		}
	}

	// At level 2, suggest viewing relations
	if progress.CurrentLevel == model.LevelTree && progress.TreeExpansionsCount >= 3 {
		if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeViewRelations); !hasSeenNudge {
			return &Nudge{
				Type:    model.NudgeViewRelations,
				Title:   "Explore Dependencies",
				Message: "See how files depend on each other and understand the architecture.",
				Action:  "View relations",
				Icon:    "share-2",
			}, nil
		}
	}

	// At higher levels, suggest export
	if progress.CurrentLevel >= model.LevelTechnical && progress.FilesViewedCount >= 10 {
		if hasSeenNudge, _ := s.repo.HasSeenNudge(ctx, projectID, model.NudgeExport); !hasSeenNudge {
			return &Nudge{
				Type:    model.NudgeExport,
				Title:   "Take It With You",
				Message: "Ready to export your project? Download the complete code package.",
				Action:  "Export project",
				Icon:    "download",
			}, nil
		}
	}

	return nil, nil
}

// RecordNudgeAction records user interaction with a nudge
func (s *NudgeService) RecordNudgeAction(ctx context.Context, projectID uuid.UUID, nudgeType model.NudgeType, action string) error {
	switch action {
	case "shown":
		return s.repo.RecordNudgeShown(ctx, projectID, nudgeType)
	case "dismissed":
		// Find the most recent nudge of this type and mark dismissed
		nudges, err := s.repo.GetRecentNudges(ctx, projectID, 10)
		if err != nil {
			return err
		}
		for _, n := range nudges {
			if n.NudgeType == nudgeType && n.DismissedAt == nil {
				return s.repo.RecordNudgeDismissed(ctx, n.ID)
			}
		}
	case "clicked":
		nudges, err := s.repo.GetRecentNudges(ctx, projectID, 10)
		if err != nil {
			return err
		}
		for _, n := range nudges {
			if n.NudgeType == nudgeType && n.ClickedAt == nil {
				return s.repo.RecordNudgeClicked(ctx, n.ID)
			}
		}
	}
	return nil
}
