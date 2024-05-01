package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
)

func (s *service) UpdateRacerInfo(ctx context.Context, racer models.RacerUpdate) error {
	racerUUID, err := uuid.Parse(racer.Id)
	if err != nil {
		return fmt.Errorf("can't parse uuid %w", err)
	}

	r := models.RacerUpdateRepo{
		Id:       racerUUID,
		Username: racer.Username,
		Email:    racer.Email,
	}

	err = s.repo.Racer.UpdateRacer(ctx, r)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
