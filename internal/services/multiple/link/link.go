package link

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"log"
	"time"
)

type Checker interface {
	Create(id string) (uuid.UUID, error)
	Check(ctx context.Context, link string) error
	Kill(ticker *time.Ticker)
}

type link struct {
	repo *repository.Repo
}

func NewLink(repo *repository.Repo) Checker {
	return &link{
		repo: repo,
	}
}

func (l *link) Kill(ticker *time.Ticker) {
	for range ticker.C {
		err := l.repo.Link.Remove(context.TODO(), time.Now())
		if err != nil {
			log.Printf("error cleaning expired links %v\n", err)
		}
	}
}

func (l *link) Check(ctx context.Context, link string) error {
	lnk, err := uuid.Parse(link)
	if err != nil {
		return err
	}
	ex, err := l.repo.Link.Check(ctx, lnk)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("link doesn't exist")
	}

	return nil
}
func (l *link) Create(id string) (uuid.UUID, error) {
	trackId := uuid.New()

	err := l.repo.Link.Add(context.TODO(), trackId, id, time.Now())
	if err != nil {
		return [16]byte{}, err
	}
	return trackId, nil
}
