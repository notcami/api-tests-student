package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление, проверяет ответ и регистрирует очистку записи.
func AddFavorite(t *testing.T, token, id string) models.Favorite {
	t.Helper()
	request := models.AddRequest{AdvertisementID: id}
	body := addFavorite(t, token, request, http.StatusCreated)
	result := response.Decode[models.Favorite](t, body)
	require.Equal(t, id, result.Advertisement.ID, "В избранное добавлено другое объявление")
	require.NotEmpty(t, result.FavoriteID, "Не заполнен favorite_id")
	require.NotEmpty(t, result.AddedAt, "Не заполнено время добавления в избранное")
	return result
}

// addFavorite регистрирует очистку любой созданной записи до проверки ожидаемого статуса.
func addFavorite(t *testing.T, token string, request any, status int) []byte {
	t.Helper()
	raw := client.HttpPostFavorite(t, token, request)
	body := response.Read(t, raw)
	if raw.StatusCode == http.StatusCreated {
		created := response.Decode[models.Favorite](t, body)
		if created.Advertisement.ID != "" {
			t.Cleanup(func() {
				t.Logf("[Cleanup] Удаляем из избранного объявление %s", created.Advertisement.ID)
				cleanup := client.HttpDeleteFavorite(t, token, created.Advertisement.ID)
				response.Read(t, cleanup, http.StatusNoContent, http.StatusNotFound)
			})
		}
	}
	response.RequireStatus(t, raw, body, status)
	return body
}

// AddFavoriteWithError проверяет отказ добавления в избранное.
func AddFavoriteWithError(t *testing.T, token string, request any, status int) response.Error {
	t.Helper()
	return response.Decode[response.Error](t, addFavorite(t, token, request, status))
}

// GetFavorites возвращает все страницы избранного, чтобы проверка не зависела от позиции записи.
func GetFavorites(t *testing.T, token string) models.List {
	t.Helper()
	raw := client.HttpGetFavorites(t, token, 1)
	result := response.Decode[models.List](t, response.Read(t, raw, http.StatusOK))
	for page := 2; page <= result.TotalPages; page++ {
		raw = client.HttpGetFavorites(t, token, page)
		next := response.Decode[models.List](t, response.Read(t, raw, http.StatusOK))
		result.Items = append(result.Items, next.Items...)
	}
	require.Len(t, result.Items, result.Total, "Количество полученных записей избранного не совпадает с total")
	return result
}

// GetFavoritesWithError проверяет отказ чтения избранного.
func GetFavoritesWithError(t *testing.T, token string, status int) response.Error {
	t.Helper()
	raw := client.HttpGetFavorites(t, token, 1)
	return response.Decode[response.Error](t, response.Read(t, raw, status))
}

// DeleteFavorite удаляет запись и проверяет пустой ответ 204.
func DeleteFavorite(t *testing.T, token, id string) {
	t.Helper()
	raw := client.HttpDeleteFavorite(t, token, id)
	require.Empty(t, response.Read(t, raw, http.StatusNoContent), "Ответ DELETE должен быть пустым")
}

// DeleteFavoriteWithError проверяет отказ удаления из избранного.
func DeleteFavoriteWithError(t *testing.T, token, id string, status int) response.Error {
	t.Helper()
	raw := client.HttpDeleteFavorite(t, token, id)
	return response.Decode[response.Error](t, response.Read(t, raw, status))
}
