package link

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (l *link) Check(ctx context.Context, link string) error {

	lnk, err := uuid.Parse(link)
	if err != nil {
		return fmt.Errorf("fail to parse link %v : err=%w", link, err)
	}

	ex, err := l.repo.Link.Check(ctx, lnk)
	if err != nil {
		return fmt.Errorf("fail to check link %v : err=%w", link, err)
	}

	if !ex {
		return fmt.Errorf("link %v doesn't exist", link)
	}

	return nil
}
