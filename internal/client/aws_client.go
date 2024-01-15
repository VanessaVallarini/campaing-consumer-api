package client

import (
	"campaing-comsumer-service/internal/util"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/lockp111/go-easyzap"
)

type IAwsClient interface {
	SendMessage(ctx context.Context, data interface{}, queueUrl *string) error
}

type AwsClient struct {
	client *sqs.Client
}

func NewAwsClient(ctx context.Context, awsURL, region string) *AwsClient {
	// customResolver is required here since we use localstack and need to point the aws url to localhost.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           awsURL,
			SigningRegion: region,
		}, nil

	})

	// load the default aws config along with custom resolver.
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	return &AwsClient{
		client: sqs.NewFromConfig(cfg),
	}
}

func (a *AwsClient) SendMessage(ctx context.Context, data interface{}, queue *string) error {
	stringData, er := util.ParseToString(data)
	if er == nil {
		_, err := a.client.SendMessage(ctx, &sqs.SendMessageInput{
			MessageBody: &stringData,
			QueueUrl:    queue,
		})
		if err != nil {
			easyzap.Error(ctx, err, "could not send message")
			return err
		}
	}

	return nil
}

func (a *AwsClient) ReceiveMessage(ctx context.Context, params sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return a.client.ReceiveMessage(ctx, &params)
}

func (a *AwsClient) DeleteMessage(ctx context.Context, params sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return a.client.DeleteMessage(ctx, &params)
}
