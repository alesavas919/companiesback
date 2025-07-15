package security

import (
	"companies/models"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func SecretString(secretValue string) any {
	var secretName = string(os.Getenv("SSL_SECRET_SL")) //os.Getenv("SSL_SECRET_SL")
	region := "us-east-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "Configuration not found" + os.Getenv("TEST_ENV_A")
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "Response not found" + os.Getenv("TEST_ENV_A") + err.Error()
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	var dataResponse models.InfoResp
	err = json.Unmarshal([]byte(secretString), &dataResponse)
	if err != nil {
		return err
	}
	// Your code goes here.
	return dataResponse
}
