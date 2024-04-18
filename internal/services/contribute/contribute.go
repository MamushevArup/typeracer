package contribute

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"time"
)

func (s *service) ContributeText(ctx context.Context, contributeDTO models.ContributeHandlerRequest) error {
	textInfo, err := s.convertHandlerToService(contributeDTO)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	isDuplicate, err := s.repo.Contributor.DuplicateContent(ctx, textInfo.RacerID, textInfo.ContentHash)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if isDuplicate {
		return fmt.Errorf("this text already exists, please try another one")
	}

	err = s.repo.Contributor.InsertToModeration(ctx, textInfo)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (s *service) convertHandlerToService(c models.ContributeHandlerRequest) (models.ContributeServiceRequest, error) {
	moderID, err := uuid.NewUUID()
	if err != nil {
		return models.ContributeServiceRequest{}, fmt.Errorf("can't generate uuid: %w", err)
	}

	racerUUID, err := uuid.Parse(c.RacerID)
	if err != nil {
		return models.ContributeServiceRequest{}, fmt.Errorf("can't parse uuid: %w", err)
	}

	return models.ContributeServiceRequest{
		ModerationID: moderID,
		RacerID:      racerUUID,
		Content:      c.Content,
		Author:       c.Author,
		Source:       c.Source,
		SourceTitle:  c.SourceTitle,
		SentAt:       time.Now(),
		Status:       0,
		Length:       len(c.Content),
		ContentHash:  generateHash(c.Content),
	}, nil
}

func generateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}
