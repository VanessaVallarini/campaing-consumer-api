package listener

import (
	"context"
	"encoding/json"
	"sync"

	"campaing-comsumer-service/internal/model"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/lockp111/go-easyzap"
)

type AwsClient interface {
	SendMessage(ctx context.Context, data interface{}, queueUrl *string) error
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

type CampaingService interface {
	CampaingHandler(ctx context.Context, campaing *model.Event) error
}

func EventTrackingListener(ctx context.Context, awsClient AwsClient, service CampaingService, queueUrl string) {
	waitGroup := &sync.WaitGroup{}
	for {
		sqsMessage, queueErr := awsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueUrl,
			MaxNumberOfMessages: 10,
		})
		if queueErr != nil {
			easyzap.Error(ctx, queueErr, "Status on receiving message. Check the status of event tracking queue.")
			continue
		}
		if sqsMessage.Messages != nil && len(sqsMessage.Messages) > 0 {
			waitGroup.Add(len(sqsMessage.Messages))
			for index := range sqsMessage.Messages {
				func(message types.Message) {
					var eventMessage *model.Event
					easyzap.Debug(ctx, "[Event tracking] Processing message id %v:", message.MessageId)
					defer waitGroup.Done()
					if err := json.Unmarshal([]byte(*message.Body), &eventMessage); err != nil {
						easyzap.Error(ctx, err, "[Event tracking] Error to parse message from queue. [Message Body: %v]", *message.Body)
						return
					}
					if eventMessage != nil {
						if err := service.CampaingHandler(ctx, eventMessage); err != nil {
							easyzap.Error(ctx, err, "[Event tracking] Failed to process event tracking message. [Error: %v]", err)
							return
						}
						easyzap.Debug(ctx, "[Event tracking] Deleting message.")
						if _, errorToDeleteMessage := awsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
							QueueUrl:      &queueUrl,
							ReceiptHandle: message.ReceiptHandle,
						}); errorToDeleteMessage != nil {
							easyzap.Debug(ctx, "[Event tracking] Failed to delete message.")
						}
						return
					}
				}(sqsMessage.Messages[index])
			}
			waitGroup.Wait()
		}
	}
}
