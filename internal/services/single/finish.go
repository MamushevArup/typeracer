package single

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"time"
)

func (s *service) EndRace(ctx context.Context, req models.ReqEndSingle) (models.RespEndSingle, error) {

	var resp models.RespEndSingle

	if req.Length < req.Errors {
		return resp, fmt.Errorf("errors greater than text length -> errors : %v, length :%v", req.Errors, req.Length)
	}

	resp.RacerId = req.RacerId
	resp.Wpm = int(countWPM(req.Length, req.Duration))
	resp.Accuracy = calcAccuracy(req.Errors, req.Length)
	resp.Duration = req.Duration
	resp.StartedTime = calculateStartTime(time.Now(), resp.Duration)
	resp.RaceId = s.ids.raceUUID
	resp.TextId = s.ids.textUUID

	err := s.repo.Starter.EndSingleRace(ctx, resp)
	if err != nil {
		return resp, fmt.Errorf("fail to finish race user_id=%v error=%w", resp.RacerId, err)
	}

	return resp, nil
}

func (s *service) GuestCalculate(ctx context.Context, duration, errors, length int) (models.RespEndSingle, error) {
	var resp models.RespEndSingle

	if length < errors {
		return resp, fmt.Errorf("errors greater than text length -> errors : %v, length :%v", errors, length)
	}

	resp.Wpm = int(countWPM(length, duration))
	resp.Accuracy = calcAccuracy(errors, length)
	resp.StartedTime = calculateStartTime(time.Now(), duration)

	return resp, nil
}
