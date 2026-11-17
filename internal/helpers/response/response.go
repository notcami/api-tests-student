package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Error struct {
	Error   string          `json:"error"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}

// Read читает и закрывает тело ответа, при необходимости проверяет статус.
func Read(t *testing.T, raw *http.Response, expectedStatuses ...int) []byte {
	t.Helper()
	require.NotNil(t, raw, "Сервер не вернул HTTP-ответ")
	defer func() {
		assert.NoError(t, raw.Body.Close(), "Не удалось закрыть тело HTTP-ответа")
	}()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать HTTP-ответ")
	RequireStatus(t, raw, body, expectedStatuses...)
	return body
}

// RequireStatus проверяет статус и добавляет тело ответа в сообщение об ошибке.
func RequireStatus(t *testing.T, raw *http.Response, body []byte, expectedStatuses ...int) {
	t.Helper()
	if len(expectedStatuses) > 0 {
		t.Logf("[Result] HTTP %d; ожидаемые статусы: %v", raw.StatusCode, expectedStatuses)
		require.Contains(t, expectedStatuses, raw.StatusCode,
			"Неожиданный HTTP-статус; тело ответа: %s", body)
	}
}

// Decode разбирает JSON-ответ в модель и показывает тело при ошибке разбора.
func Decode[T any](t *testing.T, body []byte) T {
	t.Helper()
	var result T
	require.NoError(t, json.Unmarshal(body, &result), "Ответ не соответствует модели: %s", body)
	return result
}
