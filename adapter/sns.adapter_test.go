package adapter

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	mock "github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	var tests = []struct {
		name   string
		out    *sns.PublishOutput
		errRes error
	}{
		{
			name: "publish success",
			out: &sns.PublishOutput{
				MessageId: aws.String("1233313"),
			},
			errRes: nil,
		},
		{
			name:   "publish error",
			out:    nil,
			errRes: fmt.Errorf("aws error"),
		},
	}

	mock.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			apiMock := mock.Mock[snsApi]()
			mock.When(apiMock.Publish(mock.AnyContext(), mock.Any[*sns.PublishInput]())).ThenReturn(tt.out, tt.errRes)

			underTest := NewSnsActions(apiMock)

			id, err := underTest.Publish("arn", aws.String("xxx"), "COMMAND")
			assert.Equal(subtest, tt.errRes, err)
			if err == nil {
				assert.Equal(subtest, *tt.out.MessageId, id)
			}
		})
	}
}
