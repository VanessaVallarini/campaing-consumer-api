package client

import (
	"campaing-comsumer-service/internal/metrics"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/lockp111/go-easyzap"
)

type Aws struct {
	client  *sqs.Client
	metrics *metrics.Metrics
}

func NewAwsClient(metrics *metrics.Metrics, awsURL, region string) *Aws {
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
		client:  sqs.NewFromConfig(cfg),
		metrics: metrics,
	}
}

func (a *Aws) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	msg, err := a.client.ReceiveMessage(ctx, params)
	if err != nil {
		mv := []string{"error", "receiving"}
		a.metrics.EventTrackingListener.WithLabelValues(mv...).Inc()
		return nil, err
	}
	return msg, nil
}

func (a *Aws) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	msg, err := a.client.DeleteMessage(ctx, params)
	if err != nil {
		mv := []string{"error", "delete"}
		a.metrics.EventTrackingListener.WithLabelValues(mv...).Inc()
		return nil, err
	}
	return msg, nil
}
