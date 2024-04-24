package aws

import "mime/multipart"

type CloudService interface {
	UploadOne(fileHeader *multipart.FileHeader) (uploadID string, err error)
	GetOneUrl(uploadID string) (url string, err error)
}
