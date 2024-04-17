package link

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (l *link) Create(ctx context.Context, id string) (uuid.UUID, error) {
	// link creation
	trackId, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to generate link user:%v : err=%w", id, err)
	}
	// TODO send link using channel read only to the multiple

	err = l.repo.Link.Add(ctx, trackId, id, time.Now())
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to add link link: %v, user : %v, error :%w", trackId, id, err)
	}

	return trackId, nil
}
