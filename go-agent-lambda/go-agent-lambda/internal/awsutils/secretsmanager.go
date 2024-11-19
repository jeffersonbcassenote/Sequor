package awsutils

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecret(secretName string) (secretString string, decodedBinarySecret string, err error) {

	region := os.Getenv("REGION")

	myCustomResolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		return endpoints.ResolvedEndpoint{
			URL:           os.Getenv("SECRET_MANAGER_ENDPOINT"),
			SigningRegion: region,
		}, nil
	}

	session := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(region),
		EndpointResolver: endpoints.ResolverFunc(myCustomResolver),
	}))

	svc := secretsmanager.New(session)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				log.Printf("[Error] " + secretsmanager.ErrCodeDecryptionFailure + ": " + aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				log.Printf("[Error] " + secretsmanager.ErrCodeInternalServiceError + ": " + aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				log.Printf("[Error] " + secretsmanager.ErrCodeInvalidParameterException + ": " + aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				log.Printf("[Error] " + secretsmanager.ErrCodeInvalidRequestException + ": " + aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				log.Printf("[Error] " + secretsmanager.ErrCodeResourceNotFoundException + ": " + aerr.Error())
			}
		} else {
			log.Printf("[Error] %s", err.Error())
		}
		return
	}

	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Printf("[Error] Base64 Decode Error: %s", err)
			return "", "", err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
	}

	return
}
