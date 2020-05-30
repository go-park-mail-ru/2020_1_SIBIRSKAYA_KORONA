package repository

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	// "image"
	// _ "image/jpeg"
	// _ "image/png"
)

type S3Store struct {
	sessionS3 *session.Session
	bucket    string
}

func CreateS3Repository(sessS3 *session.Session, bucket_ string) attach.FileRepository {
	return &S3Store{sessionS3: sessS3, bucket: bucket_}
}

func (s3Store *S3Store) UploadFile(attachFile *multipart.FileHeader, attach *models.AttachedFile) (string, error) {
	file, err := attachFile.Open()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	defer file.Close()

	filename := attach.FileKey
	uploader := s3manager.NewUploader(s3Store.sessionS3)
	manager, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(s3Store.bucket),
		Key:                aws.String(filename),
		Body:               file,
		ContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", attachFile.Filename)),
	})

	if err != nil {
		logger.Error(err)
		return "", errors.ErrBadFileUploadS3
	}

	return manager.Location, nil
}

func (s3Store *S3Store) DeleteFile(filenameKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Store.bucket),
		Key:    aws.String(filenameKey),
	}

	deleteManager := s3.New(s3Store.sessionS3)
	result, err := deleteManager.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				logger.Error(aerr)
				return errors.ErrBadFileDeleteS3
			}
		} else {
			logger.Error(err)
			return errors.ErrBadFileDeleteS3
		}
	}

	fmt.Println(result)
	return nil
}
