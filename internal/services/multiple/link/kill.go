package link

import (
	"context"
	"log"
	"time"
)

func (l *link) Kill(ticker *time.Ticker) {
	for range ticker.C {
		err := l.repo.Link.Remove(context.TODO(), time.Now())
		if err != nil {
			log.Printf("error cleaning expired links %v\n", err)
		}
	}
}
