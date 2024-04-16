package single

import (
	"context"
	"fmt"
)

func (s *service) RealTimeCalc(ctx context.Context, currentSymbol, duration int) (int, error) {
	var wpm int

	textLen, err := s.repo.Starter.GetTextLen(ctx, s.ids.textUUID)
	if err != nil {
		return -1, fmt.Errorf("%w", err)
	}

	if currentSymbol >= textLen {
		return wpm, fmt.Errorf("current index greater than text length  index : %v, length : %v", currentSymbol, textLen)
	}

	return int(countWPM(currentSymbol, duration)), nil
}
