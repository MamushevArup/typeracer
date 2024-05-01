package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
)

func (s *service) Details(ctx context.Context, racerId string) (models.RacerHandler, error) {
	var racer models.RacerHandler
	racerUUID, err := uuid.Parse(racerId)
	if err != nil {
		return racer, fmt.Errorf("failed to parse racer id: %w", err)
	}

	racerRepo, err := s.repo.Racer.Info(ctx, racerUUID)
	if err != nil {
		return models.RacerHandler{}, fmt.Errorf("failed to get racer info: %w", err)
	}

	return convertRacerToHandler(racerRepo), nil
}

func convertRacerToHandler(r models.RacerRepository) models.RacerHandler {

	timeFormat := r.CreatedAt.Format("01:06")

	return models.RacerHandler{
		Username:      r.Username,
		Avatar:        r.Avatar,
		CreatedAt:     timeFormat,
		AvgSpeed:      r.AvgSpeed,
		LastRaceSpeed: r.LastRaceSpeed,
		BestSpeed:     r.BestSpeed,
		Races:         r.Races,
	}
}
