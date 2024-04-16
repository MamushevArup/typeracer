package single

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
)

func (s *service) StartRace(ctx context.Context, userId string) (models.SingleResponse, error) {

	var sgl models.SingleResponse
	var userUUID uuid.UUID
	var err error

	if userId == guest {
		sgl.RacerInfo.Avatar = "random_string_from_aws"
		sgl.RacerInfo.Username = guest
	} else {

		userUUID, err = convertStrToUUID(userId)
		if err != nil {
			return sgl, fmt.Errorf("convert to uuid %v:%w", userUUID, err)
		}

		racerInfo, err := s.repo.Starter.RacerInfo(ctx, userUUID)
		if err != nil {
			return sgl, fmt.Errorf("%v:%w", racerInfo, err)
		}

		sgl.RacerInfo = racerInfo
	}

	newRaceID, err := uuid.NewUUID()
	if err != nil {
		return sgl, fmt.Errorf("unable to start single race session due to %w", err)
	}

	text, err := s.repo.Starter.TextInfo(ctx)
	if err != nil {
		return sgl, fmt.Errorf("%w", err)
	}

	s.ids.raceUUID = newRaceID
	s.ids.racerUUID = userUUID
	s.ids.textUUID = text.TextID

	sgl.TextInfo = text

	return sgl, nil
}
