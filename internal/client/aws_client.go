package client

import (
	"campaing-comsumer-service/internal/util"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/lockp111/go-easyzap"
)

type Aws struct {
	client *sqs.Client
}

func NewAwsClient(awsURL, region string) *Aws {
	// customResolver is required here since we use localstack and need to point the aws url to localhost.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           awsURL,
			SigningRegion: region,
		}, nil

	})

	// load the default aws config along with custom resolver.
	cfg, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		easyzap.Panicf("configuration error: %v", err)
	}

	return &Aws{
		client: sqs.NewFromConfig(cfg),
	}
}

func (a *Aws) SendMessage(ctx context.Context, data interface{}, queue *string) error {
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

func (a *Aws) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return a.client.ReceiveMessage(ctx, params)
}

func (a *Aws) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return a.client.DeleteMessage(ctx, params)
}
