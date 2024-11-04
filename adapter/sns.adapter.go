package adapter

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"log/slog"
)

type snsActions struct {
	snsClient snsApi
}

type snsApi interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func NewSnsActions(snsClient snsApi) SnsOperations {
	return &snsActions{snsClient: snsClient}
}

func (r snsActions) Publish(topicArn string, message *string, typeMessage string) (string, error) {
	publishInput := sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  message,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"typeMessage": {
				DataType:    aws.String("String"),
				StringValue: aws.String(typeMessage),
			},
		},
	}

	p, err := r.snsClient.Publish(context.Background(), &publishInput)
	if err != nil {
		slog.Error("couldn't publish message to topic %v. error: %v", topicArn, err)
		return "", err
	}

	return *p.MessageId, err
}
