package single

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

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

func convertStrToUUID(id string) (uuid.UUID, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}
	return userUUID, err
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

	accuracy := (1 - float64(errors)/float64(length)) * 100.0

	accuracyStr := fmt.Sprintf("%.1f", accuracy)

	final, err := strconv.ParseFloat(accuracyStr, 64)
	if err != nil {
		return 0.0
	}

	return final
}

func calculateStartTime(endTime time.Time, durationInSeconds int) time.Time {
	// Subtract the duration from the end time to get the start time
	startTime := endTime.Add(-time.Second * time.Duration(durationInSeconds))
	return startTime
}
