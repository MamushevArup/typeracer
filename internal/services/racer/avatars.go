package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
)

func (s *service) Avatars(ctx context.Context) ([]models.Avatar, error) {
	avatars, err := s.repo.Racer.SelectAvatar(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return avatars, nil
}
