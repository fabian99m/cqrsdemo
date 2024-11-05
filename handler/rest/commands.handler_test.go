package handler

import (
	"github.com/fabian99m/cqrsdemo/adapter"
	e "github.com/fabian99m/cqrsdemo/errors"
	handler "github.com/fabian99m/cqrsdemo/handler/messages"
	"github.com/fabian99m/cqrsdemo/model"
	"fmt"
	"strings"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mk "github.com/ovechkin-dm/mockio/mock"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		body         string
		statusCode   int
		errorCode    int
		serviceErr   error
		serviceOut   string
		handlerError error
	}{
		{
			name:       "json invalid",
			body:       "{,}",
			errorCode:  e.GenericError.Code,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "invalid command name",
			body:       "{\"name\":\"\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}",
			errorCode:  e.MissingCommandName.Code,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "command name not found",
			body:       "{\"name\":\"test\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}",
			errorCode:  e.CommandNotRegistered.Code,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "command process ok",
			body:       "{\"name\":\"holacommand\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}",
			serviceOut: "id12131",
			statusCode: http.StatusAccepted,
		},
		{
			name:       "poublis event erropr",
			body:       "{\"name\":\"holacommand\",\"payload\":{\"id\":\"121213123\",\"type\":\"CE\"}}",
			serviceOut: "",
			serviceErr: fmt.Errorf("aws error"),
			errorCode:  e.GenericError.Code,
			statusCode: http.StatusInternalServerError,
		},
	}

	mk.SetUp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(subtest *testing.T) {
			t.Log("running test: ", tt.name)

			mockSns := mk.Mock[adapter.SnsOperations]()
			mk.When(mockSns.Publish(mk.AnyString(), mk.NotNil[*string](), mk.AnyString())).ThenReturn(tt.serviceOut, tt.serviceErr)

			commands := handler.CmdMapper{
				"holacommand": nil,
			}

			underTest := NewCommandRestHandler(mockSns, commands, &model.EventProps{QueueName: "test", TopicArn: "test"})

			recorder := httptest.NewRecorder()
			err := underTest.Process(recorder, httptest.NewRequest("POST", "/commands", strings.NewReader(tt.body)))

			if err != nil {
				var reqErr e.RequestError
				if errors.As(err, &reqErr) {
					assert.Equal(subtest, tt.statusCode, reqErr.StatusCode)
					assert.Equal(subtest, tt.errorCode, reqErr.Status.Code)
				}
			} else {
				assert.Equal(subtest, tt.statusCode, recorder.Code)
			}
		})
	}
}
