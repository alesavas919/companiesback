package security

import (
	"companies/models"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func SecretString(value int) (string, error) {
	var secretName = string(os.Getenv("SSL_SECRET_SL")) //os.Getenv("SSL_SECRET_SL")
	region := "us-east-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "Configuration not found" + os.Getenv("TEST_ENV_A"), err
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
		return "Response not found" + os.Getenv("TEST_ENV_A") + err.Error(), err
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	var dataResponse models.InfoResp
	err = json.Unmarshal([]byte(secretString), &dataResponse)
	if err != nil {
		return "", err
	}
	if value == 0 {
		return dataResponse.MCo, nil
	}
	if value == 1 {
		return dataResponse.OrdenPr, nil
	}
	if value == 2 {
		return dataResponse.UList, nil
	}
	// Your code goes here.
	return "", errors.New("NOT FOUND DATA")
}
