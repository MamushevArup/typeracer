package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"strconv"
	"time"
)

const limitMax = 10

func (s *service) RejectContent(ctx context.Context, reject models.ModerationRejectToService) error {
	if reject.ModerationID == "" {
		return fmt.Errorf("moderation id is empty")
	}
	modUUID, err := uuid.Parse(reject.ModerationID)
	if err != nil {
		return fmt.Errorf("moderation id is not valid")
	}

	if reject.Reason == "" {
		return fmt.Errorf("reason for reject must be provided")
	}

	rejectRepo := models.ModerationRejectToRepo{
		ModerationID: modUUID,
		Reason:       reject.Reason,
	}

	err = s.repo.Admin.RejectContent(ctx, rejectRepo)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (s *service) ShowContentToModerate(ctx context.Context, limit, offset, sort string) ([]models.ModerationServiceResponse, error) {
	localLimit := limitMax
	localOffset := 0
	var err error
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

	if sort != "desc" {
		sort = "asc"
	}

	details, err := s.repo.Admin.SelectModeration(ctx, localLimit, localOffset, sort)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	convertedDetails := s.convertToService(details)

	return convertedDetails, nil
}

func (s *service) convertToService(d []models.ModerationRepoResponse) []models.ModerationServiceResponse {
	resp := make([]models.ModerationServiceResponse, len(d))

	for i := range d {
		resp[i].ModerationID = d[i].ModerationID
		resp[i].ContributorName = d[i].Username
		resp[i].SentAt = d[i].SentAt.Format("01.02 15:04")
	}

	return resp
}

func (s *service) TextDetails(ctx context.Context, modId string) (models.ModerationTextDetails, error) {
	if modId == "" {
		return models.ModerationTextDetails{}, fmt.Errorf("moderation id is empty")
	}
	modUUID, err := uuid.Parse(modId)
	if err != nil {
		return models.ModerationTextDetails{}, fmt.Errorf("moderation id is not valid")
	}

	details, err := s.repo.Admin.ContentDetails(ctx, modUUID)
	if err != nil {
		return models.ModerationTextDetails{}, fmt.Errorf("%w", err)
	}

	details.ModerationID = modUUID

	return details, nil
}

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
