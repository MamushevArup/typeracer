package race

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

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

	formattedRes := fmt.Sprintf("%.2f", (1-float64(errors)/float64(length))*100.0)

	result, err := strconv.ParseFloat(formattedRes, 64)
	if err != nil {
		return 0.0
	}

	return result
}

func randomize(ids []uuid.UUID) uuid.UUID {
	if len(ids) == 0 {
		return uuid.Nil
	}

	return ids[rand.Intn(len(ids))] // Select a random UUID from the slice
}

func inputValidation(racer *models.RacerCurrentWpm, textLen int) error {
	if racer.Email == "" {
		return fmt.Errorf("email is empty got %v", racer.Email)
	}
	if racer.Duration <= 0 {
		return fmt.Errorf("duration is less than or equal to 0 got %v", racer.Duration)
	}
	if racer.Index < 0 || racer.Index > textLen {
		return fmt.Errorf("current symbol is outbound the ranges index=%v, actual text length=%v", racer.Index, textLen)
	}
	return nil
}

func (s *service) racePopulate(raceCalc models.RaceResult, raceReq models.RaceEndRequest) models.RaceResult {
	raceCalc.Accuracy = calcAccuracy(raceReq.Errors, raceReq.Length)
	raceCalc.Email = raceReq.Email
	raceCalc.Time = 120
	raceCalc.Wpm = int(countWPM(raceReq.Length, raceReq.Duration))
	raceCalc.Place = s.finishedRacers[raceReq.Email]
	return raceCalc
}

func (s *service) convertToRepo(ctx context.Context, r *models.RaceResult, link uuid.UUID) *models.RacerRepoM {
	var winner string

	for k, v := range s.finishedRacers {
		if v == 1 {
			winner = k
		}
	}

	// here call to the db to find user with such email. If success return uuid else winner was guest
	winnerID, err := s.repo.Multiple.RacerID(ctx, winner)
	if err != nil {
		return nil
	}
	if err != nil && winnerID != uuid.Nil {
		winner = winnerID.String()
	}

	return &models.RacerRepoM{
		GeneratedLink: link,
		RacerId:       r.RacerId,
		Duration:      r.Time,
		Wpm:           r.Wpm,
		Accuracy:      r.Accuracy,
		StartTime:     time.Time{},
		Winner:        winner,
		Place:         s.finishedRacers[r.Email],
		TrackSize:     len(s.finishedRacers),
	}
}

func validateEndRace(raceReq models.RaceEndRequest) error {
	if raceReq.Email == "" {
		return fmt.Errorf("email is empty got %v", raceReq.Email)
	}
	if raceReq.Errors < 0 {
		return fmt.Errorf("errors is less than zero")
	}
	if raceReq.Length < 0 {
		return fmt.Errorf("length is less than zero")
	}
	return nil
}
