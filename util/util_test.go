package util

import (
	m "github.com/fabian99m/cqrsdemo/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
