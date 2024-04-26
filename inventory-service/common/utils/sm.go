package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"inventory-system/common/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
)

type SecretManager struct {
	Region        string
	SecretName    string
	ACCESS_KEY    string
	ACCESS_SECRET string
}

func NewSecretManager(region, SecretName string, AccessKey, AccessSecret string) *SecretManager {
	return &SecretManager{
		Region:        region,
		SecretName:    SecretName,
		ACCESS_KEY:    AccessKey,
		ACCESS_SECRET: AccessSecret,
	}
}

func (sm *SecretManager) GetSecrets() (string, error) {

	log := logger.GetLogger()
	log.Info("Inside GetSecret")
	creds := credentials.NewStaticCredentialsProvider(sm.ACCESS_KEY, sm.ACCESS_SECRET, "")

	// Load default config with custom credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(sm.Region),
		config.WithCredentialsProvider(creds),
	)

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(sm.SecretName), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {

		log.Info("E", err)
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "", err
	}

	var secretString string = *result.SecretString

	return secretString, nil

}
