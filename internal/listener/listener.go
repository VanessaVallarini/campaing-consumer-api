package listener

import (
	"context"
	"encoding/json"
	"sync"

	"campaing-comsumer-service/internal/metrics"
	"campaing-comsumer-service/internal/model"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/lockp111/go-easyzap"
)

type AwsClient interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

type CampaingService interface {
	Handler(ctx context.Context, campaing *model.Event) error
}

func EventTrackingListener(ctx context.Context, metrics *metrics.Metrics, awsClient AwsClient, service CampaingService, queueUrl string) {
	waitGroup := &sync.WaitGroup{}
	for {
		sqsMessage, queueErr := awsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueUrl,
			MaxNumberOfMessages: 10,
		})
		if queueErr != nil {
			easyzap.Error(ctx, queueErr, "status on receiving message - check the status of event tracking queue.")

			continue
		}
		if sqsMessage.Messages != nil && len(sqsMessage.Messages) > 0 {
			waitGroup.Add(len(sqsMessage.Messages))
			for index := range sqsMessage.Messages {
				ctx := context.Background()
				func(message types.Message) {
					var eventMessage *model.Event
					defer waitGroup.Done()
					if err := json.Unmarshal([]byte(*message.Body), &eventMessage); err != nil {
						easyzap.Error(ctx, err, "[event tracking] error to parse message from queue. [message body: %v].", *message.Body)

						return
					}
					if eventMessage != nil {
						if err := service.Handler(ctx, eventMessage); err != nil {
							easyzap.Error(ctx, err, "[event tracking] error to process event tracking message.", err)

							return
						}

						if _, errorToDeleteMessage := awsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
							QueueUrl:      &queueUrl,
							ReceiptHandle: message.ReceiptHandle,
						}); errorToDeleteMessage != nil {
							easyzap.Error(ctx, errorToDeleteMessage, "[event tracking] error to delete message.")
						}

						mv := []string{"success", ""}
						metrics.EventTrackingListener.WithLabelValues(mv...).Inc()

						return
					}
				}(sqsMessage.Messages[index])
			}
			waitGroup.Wait()
		}
	}
}
