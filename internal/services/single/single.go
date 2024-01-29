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
	EndRace(ctx context.Context, req *models.ReqEndSingle) (*models.RespEndSingle, error)
	RacerExists(ctx context.Context, racerId uuid.UUID) (bool, error)
	RealTimeCalc(ctx context.Context, currentSymbol, duration int) (int, error)
}

type service struct {
	repo *repository.Repo
}

func NewPracticeY(repo *repository.Repo) PracticeY {
	return &service{
		repo: repo,
	}
}

var (
	lenErr     = errors.New("invalid body error cannot be greater than text length")
	startRace  = errors.New("you can't start a race")
	fetchText  = errors.New("can't fetch text")
	finishRace = errors.New("cannot finish single race")
	// TODO replace for actual error
	blah           = errors.New("blah blah")
	raceNotStarted = errors.New("race is not started")
)

func (s *service) StartRace(ctx context.Context, userId uuid.UUID) (*models.Single, error) {
	//TODO implement me
	single, err := s.repo.Starter.StartSingle(ctx, userId)
	if err != nil {
		return nil, startRace
	}
	return single, nil
}

func (s *service) EndRace(ctx context.Context, req *models.ReqEndSingle) (*models.RespEndSingle, error) {
	var resp models.RespEndSingle
	length, err := s.repo.Starter.GetTextLen(ctx)
	if length < req.Errors {
		return &resp, lenErr
	}
	resp.RacerId = req.RacerID
	resp.Wpm = int(countWPM(length, req.Duration))
	resp.Accuracy = calcAccuracy(req.Errors, length)
	resp.Duration = req.Duration
	resp.StartedTime = calculateStartTime(time.Now(), resp.Duration)
	if err != nil {
		return nil, fetchText
	}
	err = s.repo.Starter.EndSingleRace(ctx, &resp)
	if err != nil {
		log.Println("Can't end race in services")
		return &resp, finishRace
	}
	return &resp, nil
}

func (s *service) RealTimeCalc(ctx context.Context, currentSymbol, duration int) (int, error) {
	var wpm int
	textLen, err := s.repo.Starter.GetTextLen(ctx)
	if err != nil {
		return 0, raceNotStarted
	}
	if textLen == 0 {
		return wpm, raceNotStarted
	}
	if currentSymbol >= textLen {
		return wpm, errors.New("index out of bound")
	}
	return int(countWPM(currentSymbol, duration)), nil
}

func (s *service) RacerExists(ctx context.Context, racerId uuid.UUID) (bool, error) {
	exist, err := s.repo.Starter.RacerExist(ctx, racerId)
	if err != nil {
		log.Println(err)
		return false, blah
	}
	return exist, nil
}

func countWPM(length, duration int) float64 {
	const averageWordLength = 5 // Assuming an average word length of 5 characters

	// Calculate total words in the text
	totalWords := length / averageWordLength

	// Calculate WPM
	wpm := float64(totalWords) / (float64(duration) / 60.0)

	return wpm
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
