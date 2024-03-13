package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog/log"
	"github.com/xuxant/Darwin_API/pkg/models"
)

func ValidateCredentials(accessKey, secret string) (*string, error) {
	ctx := context.TODO()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secret,
			},
		}))
	if err != nil {
		return nil, err
	}

	client := sts.NewFromConfig(awsCfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("Caller Identity: %s", *identity.Account)
	return identity.Account, nil
}

func GetAllBuckets(accessKey, secretKey string) (models.Buckets, error) {
	ctx := context.TODO()
	var buckets []string
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		}))
	if err != nil {
		return models.Buckets{}, err
	}
	client := s3.NewFromConfig(awsCfg)

	result, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return models.Buckets{}, err
	}

	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}
	return models.Buckets{Buckets: buckets}, nil

}

func GetAllObjects(accessKey, secretKey, bucketName string) (models.Objects, error) {
	ctx := context.TODO()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		}))
	if err != nil {
		return models.Objects{}, err
	}
	client := s3.NewFromConfig(awsCfg)

	results, err := client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return models.Objects{}, err
	}
	var objects []string
	for _, object := range results.Contents {
		objects = append(objects, *object.Key)
	}
	return models.Objects{
		BucketName: bucketName,
		Objects:    objects,
	}, nil
}
