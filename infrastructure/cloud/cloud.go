package cloud

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	configuration "product-service-api/infrastructure/config"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func UploadImageToS3(image *multipart.FileHeader) (string, error) {
	cfg, err := configuration.LoadConfig()
	if err != nil {
		logrus.Error("failed to load configuration:", err)
		return "", err
	}

	awsAccessKeyID := cfg.CLOUDSTORAGE.AWS_ACCESS_KEY_ID
	awsSecretAccessKey := cfg.CLOUDSTORAGE.AWS_SECRET_ACCESS_KEY
	awsRegion := cfg.CLOUDSTORAGE.AWS_REGION
	bucketName := cfg.CLOUDSTORAGE.AWS_BUCKET_NAME

	maxUploadSize := int64(10 * 1024 * 1024)
	if image.Size > maxUploadSize {
		return "", errors.New("file size exceeds the maximum allowed size of 10MB")
	}

	extension := strings.ToLower(filepath.Ext(image.Filename))
	allowedExtensions := map[string]bool{".jpg": true, ".png": true, ".jpeg": true}
	if !allowedExtensions[extension] {
		return "", errors.New("invalid image file format. supported formats: .jpg, .jpeg, .png")
	}

	imagePath := uuid.New().String() + extension

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
	)
	if err != nil {
		logrus.Error("failed to create AWS config:", err)
		return "", err
	}

	svc := s3.NewFromConfig(awsCfg)

	file, err := image.Open()
	if err != nil {
		logrus.Error("failed to open file:", err)
		return "", err
	}
	defer file.Close()

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imagePath),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		logrus.Error("failed to upload file to S3:", err)
		return "", err
	}

	imageURL := "https://" + bucketName + ".s3." + awsRegion + ".amazonaws.com/" + imagePath
	return imageURL, nil
}

func UploadImageBytesToS3(imageBytes []byte, originalFilename string) (string, error) {
	cfg, err := configuration.LoadConfig()
	if err != nil {
		logrus.Error("failed to load configuration:", err)
		return "", err
	}

	awsAccessKeyID := cfg.CLOUDSTORAGE.AWS_ACCESS_KEY_ID
	awsSecretAccessKey := cfg.CLOUDSTORAGE.AWS_SECRET_ACCESS_KEY
	awsRegion := cfg.CLOUDSTORAGE.AWS_REGION
	bucketName := cfg.CLOUDSTORAGE.AWS_BUCKET_NAME

	maxUploadSize := int64(10 * 1024 * 1024)
	if int64(len(imageBytes)) > maxUploadSize {
		return "", errors.New("file size exceeds the maximum allowed size of 10MB")
	}

	extension := strings.ToLower(filepath.Ext(originalFilename))
	allowedExtensions := map[string]bool{".jpg": true, ".png": true, ".jpeg": true}
	if !allowedExtensions[extension] {
		return "", errors.New("invalid image file format. supported formats: .jpg, .jpeg, .png")
	}

	imagePath := uuid.New().String() + extension

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretAccessKey, "")),
	)
	if err != nil {
		logrus.Error("failed to create AWS config:", err)
		return "", err
	}

	svc := s3.NewFromConfig(awsCfg)

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imagePath),
		Body:   bytes.NewReader(imageBytes),
	})
	if err != nil {
		logrus.Error("failed to upload file to S3:", err)
		return "", err
	}

	imageURL := "https://" + bucketName + ".s3." + awsRegion + ".amazonaws.com/" + imagePath
	return imageURL, nil
}
