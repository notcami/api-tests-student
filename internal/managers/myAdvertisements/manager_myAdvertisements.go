package myAdvertisements

import (
	"testing"

	"api-tests-template/internal/client/http/myAdvertisements"
	"api-tests-template/internal/helpers/response"
)

// GetMyAdvertisements возвращает список объявлений авторизованного пользователя и проверяет ожидаемый статус-код
func GetMyAdvertisements(t *testing.T, token string, expectedStatusCode int, limit int, search string) string {
	t.Helper()
	raw := myAdvertisements.HttpGetMyAdvertisements(t, token, limit, search)
	return string(response.Read(t, raw, expectedStatusCode))
}
