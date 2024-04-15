package race

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"time"
)

// data struct hold values need to pass to the repo layer
type data struct {
	textID    uuid.UUID
	createdAt time.Time
}

func (s *service) CurrentSpeed(racer *models.RacerCurrentWpm, textLen int) (models.RacerSpeed, error) {

	if racer.Email == "" {
		return models.RacerSpeed{}, errors.New("email is empty")
	}
	if racer.Duration <= 0 {
		return models.RacerSpeed{}, errors.New("duration is less than or equal to 0")
	}
	if racer.Index < 0 || racer.Index > textLen {
		return models.RacerSpeed{}, errors.New("current symbol is outbound the ranges")

	}

	var racerSpeed models.RacerSpeed

	wpm := countWPM(racer.Index, racer.Duration)

	racerSpeed.Email = racer.Email
	racerSpeed.Wpm = int(wpm)

	return racerSpeed, nil
}

func (s *service) EndRace(raceReq models.RaceEndRequest, link, id string) (models.RaceResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var raceCalc models.RaceResult

	textID := s.d.textID
	if textID == uuid.Nil {
		return raceCalc, errors.New("text id is empty")
	}
	if err := validateEndRace(raceReq); err != nil {
		return raceCalc, err
	}

	linkUUID, err := uuid.Parse(link)
	if err != nil {
		return raceCalc, err
	}

	raceCalc.Accuracy = int(calcAccuracy(raceReq.Errors, raceReq.Length))
	raceCalc.Email = raceReq.Email
	raceCalc.Time = 120
	raceCalc.Wpm = int(countWPM(raceReq.Length, raceReq.Duration))

	s.mu.Lock()
	s.finishedRacers[raceReq.Email] = len(s.finishedRacers) + 1
	s.mu.Unlock()

	raceCalc.Place = s.finishedRacers[raceReq.Email]

	uid, err := uuid.Parse(id)
	if err != nil {
		return raceCalc, nil
	}

	raceCalc.RacerId = uid

	repoRacer := s.convertToRepo(ctx, &raceCalc, linkUUID)

	err = s.repo.Multiple.InsertSession(ctx, repoRacer)
	if err != nil {
		return raceCalc, err
	}

	err = s.repo.Multiple.UpdateRacerHistory(ctx, raceCalc.Wpm, uid, textID)
	if err != nil {
		return raceCalc, err
	}

	return raceCalc, nil
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
		Accuracy:      int(r.Accuracy),
		StartTime:     time.Time{},
		Winner:        winner,
		Place:         s.finishedRacers[r.Email],
		TrackSize:     len(s.finishedRacers),
	}
}

func validateEndRace(raceReq models.RaceEndRequest) error {
	if raceReq.Email == "" {
		return errors.New("email is empty")
	}
	if raceReq.Errors < 0 {
		return errors.New("errors is less than 0")
	}
	if raceReq.Length < 0 {
		return errors.New("length is less than 0")
	}
	return nil
}
