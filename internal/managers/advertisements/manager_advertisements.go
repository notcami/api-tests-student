package advertisements

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
)

// GetAdvertisementsList возвращает результаты поиска и проверяет статус 200.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()

	raw := advertisements.HttpGetAdvertisements(t, token, text)
	return response.Decode[models.AdvertisementsList](t, response.Read(t, raw, http.StatusOK))
}

// CreateAdvertisement создаёт объявление и регистрирует его удаление после теста.
func CreateAdvertisement(t *testing.T, token string, request models.CreateRequest) models.Advertisement {
	t.Helper()
	body := createAdvertisement(t, token, request, http.StatusCreated)
	result := response.Decode[models.Advertisement](t, body)
	require.NotEmpty(t, result.ID, "В ответе создания отсутствует id объявления")
	return result
}

// CreateAdvertisementWithError проверяет статус ошибочного создания и возвращает ошибку API.
func CreateAdvertisementWithError(t *testing.T, token string, request models.CreateRequest, status int) response.Error {
	t.Helper()
	return response.Decode[response.Error](t, createAdvertisement(t, token, request, status))
}

// createAdvertisement ставит успешное создание на очистку до проверки ожидаемого статуса.
func createAdvertisement(t *testing.T, token string, request models.CreateRequest, status int) []byte {
	t.Helper()
	raw := advertisements.HttpPostAdvertisement(t, token, request)
	body := response.Read(t, raw)
	if raw.StatusCode == http.StatusCreated {
		created := response.Decode[models.Advertisement](t, body)
		if created.ID != "" {
			t.Cleanup(func() {
				cleanupAdvertisement(t, token, created.ID)
			})
		}
	}
	response.RequireStatus(t, raw, body, status)
	return body
}

// cleanupAdvertisement удаляет своё объявление, проверяя отсутствие ресурса при повторном удалении с 403.
func cleanupAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	t.Logf("[Cleanup] Удаляем объявление %s", id)
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	body := response.Read(t, raw)
	if raw.StatusCode == http.StatusForbidden {
		// Стенд возвращает 403 и для уже удалённого объявления: подтверждаем отсутствие через GET.
		check := advertisements.HttpGetAdvertisement(t, id)
		response.Read(t, check, http.StatusNotFound)
		t.Logf("[Cleanup] Объявление %s уже удалено", id)
		return
	}
	response.RequireStatus(t, raw, body, http.StatusNoContent, http.StatusNotFound)
}

// GetAdvertisement возвращает объявление и проверяет статус 200.
func GetAdvertisement(t *testing.T, id string) models.Advertisement {
	t.Helper()
	raw := advertisements.HttpGetAdvertisement(t, id)
	return response.Decode[models.Advertisement](t, response.Read(t, raw, http.StatusOK))
}

// GetAdvertisementWithError проверяет ошибку чтения объявления.
func GetAdvertisementWithError(t *testing.T, id string, status int) response.Error {
	t.Helper()
	raw := advertisements.HttpGetAdvertisement(t, id)
	return response.Decode[response.Error](t, response.Read(t, raw, status))
}

// DeleteAdvertisement удаляет объявление и проверяет пустой ответ 204.
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	require.Empty(t, response.Read(t, raw, http.StatusNoContent), "Ответ DELETE должен быть пустым")
}

// GetPhotos возвращает фотографии объявления и проверяет статус 200.
func GetPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	raw := advertisements.HttpGetPhotos(t, token, id)
	return response.Decode[[]models.Photo](t, response.Read(t, raw, http.StatusOK))
}

// DownloadPhoto проверяет доступность файла и возвращает его содержимое.
func DownloadPhoto(t *testing.T, photoURL string) []byte {
	t.Helper()
	return response.Read(t, advertisements.HttpGetPhoto(t, photoURL), http.StatusOK)
}
