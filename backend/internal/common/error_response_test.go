package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponse_JSONSerialization(t *testing.T) {
	t.Run("with details", func(t *testing.T) {
		errResp := ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "The provided date format is invalid",
			Details: map[string]string{"field": "date"},
		}

		bytes, err := json.Marshal(errResp)
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(bytes, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "INVALID_INPUT", unmarshaled["code"])
		assert.Equal(t, "The provided date format is invalid", unmarshaled["message"])
		assert.NotNil(t, unmarshaled["details"])
	})

	t.Run("without details (omitempty)", func(t *testing.T) {
		errResp := ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Resource not found",
		}

		bytes, err := json.Marshal(errResp)
		require.NoError(t, err)

		jsonStr := string(bytes)
		assert.Contains(t, jsonStr, `"code":"NOT_FOUND"`)
		assert.Contains(t, jsonStr, `"message":"Resource not found"`)
		assert.NotContains(t, jsonStr, `"details"`)
	})

	t.Run("zero values / empty string", func(t *testing.T) {
		errResp := ErrorResponse{}

		bytes, err := json.Marshal(errResp)
		require.NoError(t, err)

		jsonStr := string(bytes)
		assert.Contains(t, jsonStr, `"code":""`)
		assert.Contains(t, jsonStr, `"message":""`)
		assert.NotContains(t, jsonStr, `"details"`)
	})
}
