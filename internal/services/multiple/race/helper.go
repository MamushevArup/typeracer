package race

import (
	"github.com/google/uuid"
	"math/rand"
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

	return (1 - float64(errors)/float64(length)) * 100.0
}

func randomize(ids []uuid.UUID) uuid.UUID {
	if len(ids) == 0 {
		return uuid.Nil
	}
	return ids[rand.Intn(len(ids))] // Select a random UUID from the slice
}
