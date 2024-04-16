package single

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"log"
	"time"
)

// SERVICE LAYER
// this package stands for implement the practice yourself section only

type PracticeY interface {
	StartRace(ctx context.Context, userId string) (models.SingleResponse, error)
	EndRace(ctx context.Context, req *models.ReqEndSingle) (*models.RespEndSingle, error)
	RacerExists(ctx context.Context, racerId string) (bool, error)
	RealTimeCalc(ctx context.Context, currentSymbol, duration int) (int, error)
}

type service struct {
	repo *repository.Repo
	ids  identifiers
}

type identifiers struct {
	textUUID  uuid.UUID
	racerUUID uuid.UUID
	raceUUID  uuid.UUID
}

func NewPracticeY(repo *repository.Repo) PracticeY {
	return &service{
		repo: repo,
		ids:  identifiers{},
	}
}

var (
	lenErr         = errors.New("invalid body error cannot be greater than text length")
	startRace      = errors.New("you can't start a race")
	fetchText      = errors.New("can't fetch text")
	finishRace     = errors.New("cannot finish single race")
	raceNotStarted = errors.New("race is not started")
	guest          = "guest"
)

func (s *service) StartRace(ctx context.Context, userId string) (models.SingleResponse, error) {

	var sgl models.SingleResponse
	var userUUID uuid.UUID
	var err error

	if userId == guest {
		sgl.RacerInfo.Avatar = "random_string_from_aws"
		sgl.RacerInfo.Username = guest
	} else {

		userUUID, err = convertStrToUUID(userId)
		if err != nil {
			return sgl, fmt.Errorf("convert to uuid %v:%w", userUUID, err)
		}

		racerInfo, err := s.repo.Starter.RacerInfo(ctx, userUUID)
		if err != nil {
			return sgl, fmt.Errorf("%v:%w", racerInfo, err)
		}

		sgl.RacerInfo = racerInfo
	}

	newRaceID, err := uuid.NewUUID()
	if err != nil {
		return sgl, fmt.Errorf("unable to start single race session due to %w", err)
	}

	text, err := s.repo.Starter.TextInfo(ctx)
	if err != nil {
		return sgl, fmt.Errorf("%w", err)
	}

	s.ids.raceUUID = newRaceID
	s.ids.racerUUID = userUUID
	s.ids.textUUID = text.TextID

	sgl.TextInfo = text

	return sgl, nil
}

func (s *service) EndRace(ctx context.Context, req *models.ReqEndSingle) (*models.RespEndSingle, error) {
	var resp models.RespEndSingle
	length, err := s.repo.Starter.GetTextLen(ctx)
	if length < req.Errors {
		return &resp, lenErr
	}
	resp.RacerId = req.RacerId
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

func (s *service) RacerExists(ctx context.Context, racerId string) (bool, error) {
	if racerId == "guest" {
		return true, nil
	}

	racerUUID, err := convertStrToUUID(racerId)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	exist, err := s.repo.Starter.RacerExist(ctx, racerUUID)
	if err != nil {
		return false, fmt.Errorf("%w", err)
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

func convertStrToUUID(id string) (uuid.UUID, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}
	return userUUID, err
}
