package admin

import (
	"context"
	"fmt"
	"mime/multipart"
)

func (s *service) AddAvatar(ctx context.Context, fileHeader *multipart.FileHeader) error {
	url, err := s.s3.UploadOne(fileHeader)
	if err != nil {
		return fmt.Errorf("failed to upload avatar to s3: %w", err)
	}

	url, err = s.s3.GetOneUrl(url)
	if err != nil {
		return fmt.Errorf("failed to get avatar url: %w", err)
	}

	added, err := s.repo.Admin.InsertAvatarURL(ctx, url)
	if err != nil {
		return fmt.Errorf("error occur during insert avatar err=%w", err)
	}

	if !added {
		return fmt.Errorf("failed to insert avatar url")
	}

	return nil
}
