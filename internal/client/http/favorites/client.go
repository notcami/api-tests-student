package favorites

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpPostFavorite отправляет JSON-запрос добавления в избранное.
func HttpPostFavorite(t *testing.T, token string, request any) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").JSON(request).Expect(t).End().Response
}

// HttpGetFavorites запрашивает страницу избранного без фильтрации.
func HttpGetFavorites(t *testing.T, token string, page int) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", strconv.Itoa(page)).Query("page_size", strconv.Itoa(api.PageSize)).
		Expect(t).End().Response
}

// HttpDeleteFavorite удаляет запись избранного по идентификатору объявления.
func HttpDeleteFavorite(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.FavoritesPath + "/" + id).
		Expect(t).End().Response
}
