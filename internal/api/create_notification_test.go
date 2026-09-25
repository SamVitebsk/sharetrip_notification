package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"sharetrip_notification/internal/api/openapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNotification(t *testing.T) {
	t.Parallel()

	t.Run("успешное создание уведомления", func(t *testing.T) {
		t.Parallel()

		payload := map[string]interface{}{
			"message": "hello test",
		}

		reqBody := openapi.CreateNotificationRequest{
			RecipientId: "user-test-1",
			Type:        "email",
			Payload:     payload,
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respBody openapi.NotificationResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)

		require.NotNil(t, respBody.Id)
		require.True(t, time.Since(respBody.CreatedAt).Seconds() < 5, "created at should be very recent")

		expected := openapi.NotificationResponse{
			Id:          respBody.Id,
			RecipientId: "user-test-1",
			Type:        "email",
			Status:      "created",
			Payload:     payload,
			CreatedAt:   respBody.CreatedAt,
		}

		assert.Equal(t, expected, respBody)
	})

	t.Run("ошибка при пустом получателе", func(t *testing.T) {
		t.Parallel()

		reqBody := openapi.CreateNotificationRequest{
			RecipientId: "",
			Type:        "email",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
