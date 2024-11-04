package model

type EventProps struct {
	QueueName string
	TopicArn  string
}

type BucketProps struct {
	Name    string
	MaxSize int32
}
