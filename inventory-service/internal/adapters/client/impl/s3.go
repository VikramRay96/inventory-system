package impl

import (
	"bitbucket.org/kodnest/go-common-libraries/logger"
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type S3 struct {
	AwsSession *session.Session
}

func NewS3() *S3 {
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("aws.region"))})
	if err != nil {
		return nil
	}
	return &S3{
		AwsSession: awsSession,
	}
}

func (s S3) UploadFile(serviceName string, serviceType string, topicName string, topicId string, resource []byte, fileRequestId string, contentType string, fileExtension string, fType string) error {
	methodName := "UploadVideoFile"
	log := logger.New(logger.Info)

	log.Info("")

	// Upload the file to S3 and make it publicly accessible
	_, putObjErr := s3.New(s.AwsSession).PutObject(&s3.PutObjectInput{
		Body:          bytes.NewReader(resource),
		Bucket:        aws.String(viper.GetString("S3BucketName")),
		Key:           aws.String("kod/" + serviceName + "/" + serviceType + "/" + topicId + "/" + topicName + "_" + topicId + "/" + fType + "/" + fType + "_" + fileRequestId + fileExtension),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(resource))),
	})
	if putObjErr != nil {
		log.Error("Inside " + methodName + " error: " + putObjErr.Error() + " while uploading resource for topic: " + topicId)
		return putObjErr
	}
	return nil
}
