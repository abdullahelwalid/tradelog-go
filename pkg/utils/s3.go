package utils

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"context"
	"log"
	"time"
)

type S3Presign struct {
	PresignClient *s3.PresignClient
}


func (presigner S3Presign) GeneratePutURL(ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (string, error) {
	request, err := presigner.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return "", err
	}
	return request.URL, nil
}
