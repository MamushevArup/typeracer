package admin

import (
	"context"
	"github.com/MamushevArup/typeracer/adapters/avatar/aws"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"mime/multipart"
)

const limitMax = 10

type Admin interface {
	ShowContentToModerate(ctx context.Context, limit, offset, sort string) ([]models.ModerationServiceResponse, error)
	TextDetails(ctx context.Context, modId string) (models.ModerationTextDetails, error)
	ApproveContent(ctx context.Context, modId string) error
	RejectContent(ctx context.Context, reject models.ModerationRejectToService) error
	AddAvatar(ctx context.Context, fileHeader *multipart.FileHeader) error
}

type service struct {
	repo *repository.Repo
	s3   aws.CloudService
}

func New(repo *repository.Repo, s3 aws.CloudService) Admin {
	return &service{repo: repo, s3: s3}
}
