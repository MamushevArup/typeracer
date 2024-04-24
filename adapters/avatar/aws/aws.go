package aws

import (
	"bytes"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"image/png"
	"mime/multipart"
	"time"
)

type awsService struct {
	service    *s3.S3
	bucketName string
}

const (
	filePreSignDuration = time.Hour * 24
)

func New(cfg *config.Config) (CloudService, error) {

	// create a new session of aws
	newSession, err := session.NewSession(&aws.Config{
		Region:      &cfg.AWS.Region,
		Credentials: credentials.NewStaticCredentials(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session for aws s3 error: %w", err)
	}

	// create new aws service using the aws session
	service := s3.New(newSession)

	return &awsService{
		service:    service,
		bucketName: cfg.AWS.BucketName,
	}, nil
}

func (c *awsService) UploadOne(fileHeader *multipart.FileHeader) (string, error) {

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file \nerror:%w", err)
	}

	// Decode the PNG image
	img, err := png.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG image: %w", err)
	}

	// Resize the image to a manageable size if needed
	resizedImg := imaging.Resize(img, 100, 0, imaging.Lanczos)

	// Create a buffer to store the resized PNG image
	buf := new(bytes.Buffer)
	err = png.Encode(buf, resizedImg)
	if err != nil {
		return "", fmt.Errorf("failed to encode resized image: %w", err)
	}

	awsFileID := uuid.New().String()

	_, err = c.service.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(awsFileID),
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to aws bucket \nerror: %w", err)
	}

	return awsFileID, nil
}

func (c *awsService) GetOneUrl(uploadID string) (url string, err error) {

	req, _ := c.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(uploadID),
	})

	url, err = req.Presign(filePreSignDuration)

	return
}
