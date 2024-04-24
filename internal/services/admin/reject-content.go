package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
)

func (s *service) RejectContent(ctx context.Context, reject models.ModerationRejectToService) error {
	if reject.ModerationID == "" {
		return fmt.Errorf("moderation id is empty")
	}
	modUUID, err := uuid.Parse(reject.ModerationID)
	if err != nil {
		return fmt.Errorf("moderation id is not valid")
	}

	if reject.Reason == "" {
		return fmt.Errorf("reason for reject must be provided")
	}

	rejectRepo := models.ModerationRejectToRepo{
		ModerationID: modUUID,
		Reason:       reject.Reason,
	}

	err = s.repo.Admin.RejectContent(ctx, rejectRepo)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
