package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
)

func (s *service) HistorySingleText(ctx context.Context, racerId string, singleID uuid.UUID) (models.SingleHistoryText, error) {
	racerUUID, err := uuid.Parse(racerId)
	if err != nil {
		return models.SingleHistoryText{}, fmt.Errorf("invalid id, err=%v", err)
	}

	text, err := s.repo.Racer.SelectHistoryText(ctx, racerUUID, singleID)
	if err != nil {
		return text, fmt.Errorf("%w", err)
	}

	return text, nil
}
