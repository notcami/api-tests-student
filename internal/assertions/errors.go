package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/helpers/response"
)

// APIError проверяет точные значения кода и сообщения ошибки API.
func APIError(t *testing.T, actual response.Error, code, message string) {
	t.Helper()
	require.Equal(t, code, actual.Error, "Некорректный код ошибки API: %+v", actual)
	require.Equal(t, message, actual.Message, "Некорректное сообщение ошибки API")
	t.Logf("[Check] Получена ожидаемая ошибка %s: %s", code, message)
}
