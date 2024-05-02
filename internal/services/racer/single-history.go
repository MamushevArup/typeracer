package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"strconv"
)

var limitMax = 10

func (s *service) SingleHistory(ctx context.Context, id, limit, offset string) ([]models.SingleHistoryHandler, error) {

	racerUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id, err=%v", err)
	}

	localLimit := limitMax
	localOffset := 0
	if limit != "" {
		localLimit, err = strconv.Atoi(limit)
		if err != nil {
			return nil, fmt.Errorf("limit must be a number")
		}
		if localLimit == 0 || localLimit > limitMax {
			localLimit = limitMax
		}
	}

	if offset != "" {
		localOffset, err = strconv.Atoi(offset)
		if err != nil {
			return nil, fmt.Errorf("offset must be a number")
		}
		if localOffset < 0 {
			localOffset = 0
		}
	}

	history, err := s.repo.Racer.SingleHistoryRows(ctx, racerUUID, localLimit, localOffset)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	historyHandler := s.convertHistoryToHandler(history)

	return historyHandler, nil
}

func (s *service) convertHistoryToHandler(r []models.SingleHistory) []models.SingleHistoryHandler {
	resp := make([]models.SingleHistoryHandler, len(r))
	for i := range r {
		resp[i].SingleID = r[i].SingleID
		resp[i].Speed = r[i].Speed
		resp[i].Duration = r[i].Duration
		resp[i].Accuracy = r[i].Accuracy
		resp[i].StartedAt = r[i].StartedAt.Format("02:01:06")
	}
	return resp
}
