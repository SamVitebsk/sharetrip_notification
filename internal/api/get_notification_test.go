package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"sharetrip_notification/internal/api/openapi"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNotification(t *testing.T) {
	t.Parallel()

	t.Run("успешное получение уведомления", func(t *testing.T) {
		t.Parallel()

		id := insertNotification(t, "user-get-1", "push", "created")

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/notifications/%s", id), nil)
		require.NoError(t, err)

		resp, err := testApp.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody openapi.NotificationResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)

		require.NotNil(t, respBody.Id)
		require.True(t, time.Since(respBody.CreatedAt).Seconds() < 5, "created at should be very recent")

		expected := openapi.NotificationResponse{
			Id:          id,
			RecipientId: "user-get-1",
			Type:        "push",
			Status:      "created",
			Payload:     map[string]interface{}{"key": "value"},
			CreatedAt:   respBody.CreatedAt,
		}

		assert.Equal(t, expected, respBody)
	})

	t.Run("ошибка если не найдено", func(t *testing.T) {
		t.Parallel()

		randomID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/notifications/%s", randomID), nil)

		resp, err := testApp.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("ошибка при невалидном uuid", func(t *testing.T) {
		t.Parallel()

		req, _ := http.NewRequest(http.MethodGet, "/notifications/invalid-uuid", nil)

		resp, err := testApp.Test(req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func insertNotification(t *testing.T, recipientID, notifType, status string) uuid.UUID {
	t.Helper()

	id := uuid.New()
	now := time.Now().UTC().Truncate(time.Microsecond)
	payload := []byte(`{"key": "value"}`)

	const query = `
		INSERT INTO notifications (id, recipient_id, type, status, payload, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := testPool.Exec(context.Background(), query, id, recipientID, notifType, status, payload, now)
	require.NoError(t, err, "ошибка вставки тестового уведомления")

	return id
}
