package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
)

func (s *service) TextDetails(ctx context.Context, modId string) (models.ModerationTextDetails, error) {
	if modId == "" {
		return models.ModerationTextDetails{}, fmt.Errorf("moderation id is empty")
	}
	modUUID, err := uuid.Parse(modId)
	if err != nil {
		return models.ModerationTextDetails{}, fmt.Errorf("moderation id is not valid")
	}

	details, err := s.repo.Admin.ContentDetails(ctx, modUUID)
	if err != nil {
		return models.ModerationTextDetails{}, fmt.Errorf("%w", err)
	}

	details.ModerationID = modUUID

	return details, nil
}

