package single

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"log"
	"time"
)

// SERVICE LAYER
// this package stands for implement the practice yourself section only

type PracticeY interface {
	StartRace(ctx context.Context, userId uuid.UUID) (*models.Single, error)
	EndRace(ctx context.Context, req models.ReqEndSingle) (*models.RespEndSingle, error)
}

type service struct {
	repo *repository.Repo
}

func NewPracticeY(repo *repository.Repo) PracticeY {
	return &service{
		repo: repo,
	}
}

func (s *service) StartRace(ctx context.Context, userId uuid.UUID) (*models.Single, error) {
	//TODO implement me
	single, err := s.repo.Starter.StartSingle(ctx, userId)
	if err != nil {
		return nil, errors.New("you can't start a race")
	}
	return single, nil
}

func (s *service) EndRace(ctx context.Context, req models.ReqEndSingle) (*models.RespEndSingle, error) {
	var resp models.RespEndSingle
	length, err := s.repo.Starter.GetTextLen(ctx)
	resp.Wpm = countWPM(length, req.Duration)
	resp.Accuracy = calcAccuracy(req.Errors, length)
	resp.Duration = req.Duration
	resp.RacerId = req.RacerID
	resp.StartedTime = calculateStartTime(time.Now(), resp.Duration)
	if err != nil {
		return nil, errors.New("can't fetch text")
	}
	err = s.repo.Starter.EndSingleRace(ctx, resp)
	if err != nil {
		log.Println("Can't end race in services")
		return &resp, errors.New("cannot finish single race")
	}
	return &resp, nil
}

func countWPM(length, duration int) int {
	if duration == 0 {
		return 0
	}
	durInMinute := int(float64(duration) / 60)
	return length / durInMinute
}
func calcAccuracy(errors, length int) float64 {
	if length == 0 {
		return 0.0
	}

	return (1 - float64(errors)/float64(length)) * 100.0
}
func calculateStartTime(endTime time.Time, durationInSeconds int) time.Time {
	// Subtract the duration from the end time to get the start time
	startTime := endTime.Add(-time.Second * time.Duration(durationInSeconds))
	return startTime
}
