package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"strconv"
)

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
