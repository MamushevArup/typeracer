package race

import (
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
)

func (s *service) CurrentSpeed(racer *models.RacerCurrentWpm, textLen int) (models.RacerSpeed, error) {

	err := inputValidation(racer, textLen)
	if err != nil {
		return models.RacerSpeed{}, fmt.Errorf("input validation failed %w", err)
	}

	var racerSpeed models.RacerSpeed

	wpm := countWPM(racer.Index, racer.Duration)

	racerSpeed.Email = racer.Email
	racerSpeed.Wpm = int(wpm)

	return racerSpeed, nil
}
