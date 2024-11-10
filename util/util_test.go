package util

import (
	"fmt"
	"testing"

	e "github.com/fabian99m/cqrsdemo/errors"
	m "github.com/fabian99m/cqrsdemo/model"

	"github.com/stretchr/testify/assert"
)

func TestIsType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		errorIn error
		ok      bool
	}{
		{
			name: "succes",
			errorIn: e.RequestError{
				StatusCode: 400,
				Status: e.Status{
					Code: 88,
				},
			},
			ok: true,
		},
		{
			name: "succes",
			errorIn: e.ApiError{
				Status: e.Status{Code: 200},
			},
			ok: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(sbt *testing.T) {
			sbt.Log("running test: ", tt.name)
			assert.Equal(t, tt.ok, IsType[e.RequestError](tt.errorIn))
		})
	}
}

func TestAs(t *testing.T) {
	ex := e.RequestError{
		StatusCode: 400,
		Status: e.Status{
			Code: 88,
		},
	}

	res, ok := As[e.RequestError](ex)

	assert.True(t, ok)
	assert.Equal(t, ex.StatusCode, res.StatusCode)
	assert.Equal(t, ex.Status.Code, res.Status.Code)
}

func TestAsFail(t *testing.T) {
	ex := e.ApiError{
		Status: e.Status{Code: 200},
	}

	_, ok := As[e.RequestError](ex)

	assert.False(t, ok)
}

func TestUnmarshalBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload string
		success bool
	}{
		{
			name:    "sucess",
			payload: "{\"id\":\"123\",\"type\":\"CC\"}",
			success: true,
		},
		{
			name:    "error",
			payload: "{\"id\":\"123\",\"type\":\"CC\",}",
			success: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(sbt *testing.T) {
			sbt.Log("running test: ", tt.name)

			dni, err := UnmarshalTo[m.Dni]([]byte(tt.payload))

			assert.Equal(sbt, tt.success, err == nil)
			if err == nil {
				assert.Equal(sbt, "123", dni.Id)
			}
		})
	}
}

type testStruct struct {
	Name string `validate:"required"`
}

func TestValidateStruct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		testStruct testStruct
		success    bool
	}{
		{
			name: "error",
			testStruct: testStruct{
				Name: "",
			},
			success: false,
		},
		{
			name: "success",
			testStruct: testStruct{
				Name: "dasda",
			},
			success: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(sbt *testing.T) {
			sbt.Log("running test: ", tt.name)
			assert.Equal(sbt, tt.success, ValidateStruct(tt.testStruct) == nil)
		})
	}
}

func TestGetValidaiton(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		err     error
		success bool
	}{
		{
			name:    "success",
			err:     ValidateStruct(testStruct{Name: ""}),
			success: true,
		},
		{
			name:    "invalid error",
			err:     fmt.Errorf("error"),
			success: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(sbt *testing.T) {
			sbt.Log("running test: ", tt.name)
			assert.Equal(sbt, tt.success, len(GetValidations(tt.err)) > 0)
		})
	}
}
