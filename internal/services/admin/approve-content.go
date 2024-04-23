package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"time"
)

func (s *service) ApproveContent(ctx context.Context, modId string) error {
	if modId == "" {
		return fmt.Errorf("moderation id is empty")
	}
	modUUID, err := uuid.Parse(modId)
	if err != nil {
		return fmt.Errorf("moderation id is not valid")
	}

	info, err := s.repo.Admin.ModerationById(ctx, modUUID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	textDetails, err := s.convertToText(info)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	contributor := s.convertToContributor(info, textDetails.TextID)

	transaction := models.TextAcceptTransaction{Text: textDetails, Contributor: contributor}

	err = s.repo.Admin.ApproveContent(ctx, transaction)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = s.repo.Admin.DeleteModerationText(ctx, modUUID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (s *service) convertToText(info models.ModerationApprove) (models.ApproveToText, error) {
	textId, err := uuid.NewUUID()
	if err != nil {
		return models.ApproveToText{}, fmt.Errorf("can't generate text id")
	}
	return models.ApproveToText{
		TextID:        textId,
		ContributorID: info.RacerID,
		Content:       info.Content,
		Author:        info.Author,
		AcceptedAt:    time.Now(),
		Length:        len(info.Content),
		Source:        info.Source,
		SourceTitle:   info.SourceTitle,
	}, nil
}

func (s *service) convertToContributor(info models.ModerationApprove, textUUID uuid.UUID) models.ApproveToContributor {
	return models.ApproveToContributor{
		ContributorID: info.RacerID,
		SentAt:        info.SentAt,
		TextID:        textUUID,
	}
}
