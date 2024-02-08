package contribute

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"log"
	"time"
)

type Contributor interface {
	ContributeText(ctr models.ContributeText) error
	RacerExists(ctx context.Context, racerId uuid.UUID) (bool, error)
}

type service struct {
	repo *repository.Repo
}

func NewContribute(repo *repository.Repo) Contributor {
	return &service{
		repo: repo,
	}
}

func (s *service) ContributeText(ctr models.ContributeText) error {
	ctr.Length = len(ctr.Content)
	ctr.TextID = uuid.New()
	ctr.SentAt = time.Now()
	if ctr.RacerID == [16]byte{} {
		return errors.New("uuid is empty")
	}
	return s.repo.Contributor.Contribute(context.Background(), ctr)
}

func (s *service) RacerExists(ctx context.Context, racerId uuid.UUID) (bool, error) {
	exist, err := s.repo.Contributor.RacerExist(ctx, racerId)
	if err != nil {
		log.Println(err)
		return false, errors.New("user doesn't exist")
	}
	return exist, nil
}
