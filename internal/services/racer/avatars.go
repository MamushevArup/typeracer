package racer

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
)

var ErrNegativeId = errors.New("id can't be negative")

func (s *service) Avatars(ctx context.Context) ([]models.Avatar, error) {
	avatars, err := s.repo.Racer.SelectAvatar(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return avatars, nil
}

func (s *service) UpdateAvatar(ctx context.Context, avatar models.AvatarUpdate) error {
	if avatar.Id <= 0 {
		return ErrNegativeId
	}

	racerUUID, err := uuid.Parse(avatar.RacerId)
	if err != nil {
		return fmt.Errorf("can't parse uuid %w", err)
	}

	avatarRepo := convertAvatarToRepo(avatar.Id, racerUUID)

	err = s.repo.Racer.UpdateAvatar(ctx, avatarRepo)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func convertAvatarToRepo(id int, racerId uuid.UUID) models.AvatarUpdateRepo {
	return models.AvatarUpdateRepo{
		Id:      id,
		RacerID: racerId,
	}
}
