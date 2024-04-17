package race

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"time"
)

// data struct hold values need to pass to the repo layer
type data struct {
	textID    uuid.UUID
	createdAt time.Time
}

func (s *service) EndRace(raceReq models.RaceEndRequest, link, id string) (models.RaceResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout+(1*time.Second))
	defer cancel()

	var raceCalc models.RaceResult

	textID := s.d.textID

	if textID == uuid.Nil {
		return raceCalc, fmt.Errorf("text id is empty link=%v, user=%v", link, id)
	}

	if err := validateEndRace(raceReq); err != nil {
		return raceCalc, fmt.Errorf("can't parse request link=%v, user=%v, err=%w", link, id, err)
	}

	linkUUID, err := uuid.Parse(link)
	if err != nil {
		return raceCalc, fmt.Errorf("can't parse link to uuid link=%v, user=%v, err=%w", link, id, err)
	}

	s.mu.Lock()
	s.finishedRacers[raceReq.Email] = len(s.finishedRacers) + 1
	s.mu.Unlock()

	raceCalc = s.racePopulate(raceCalc, raceReq)

	uid, err := uuid.Parse(id)
	if err != nil {
		return raceCalc, nil
	}

	raceCalc.RacerId = uid

	repoRacer := s.convertToRepo(ctx, &raceCalc, linkUUID)

	err = s.repo.Multiple.InsertSession(ctx, repoRacer)
	if err != nil {
		return raceCalc, fmt.Errorf("fail to insert session email=%v due to %w", raceReq.Email, err)
	}

	err = s.repo.Multiple.UpdateRacerHistory(ctx, raceCalc.Wpm, uid, textID)
	if err != nil {
		return raceCalc, fmt.Errorf("fail to update racer stats %w", err)
	}

	return raceCalc, nil
}
