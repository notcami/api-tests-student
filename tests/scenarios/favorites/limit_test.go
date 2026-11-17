package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/support"
)

func TestFavoritesLimit(t *testing.T) {
	owner, user := support.LoginUsers(t)
	before := favorites.GetFavorites(t, user.Token)
	require.LessOrEqual(t, before.Total, api.MaxFavorites, "Исходный список уже превышает лимит избранного")
	t.Logf("[Precondition] В избранном %d записей; заполняем оставшиеся места до %d", before.Total, api.MaxFavorites)
	for count := before.Total; count < api.MaxFavorites; count++ {
		t.Logf("[Precondition] Заполняем избранное: %d/%d", count+1, api.MaxFavorites)
		created := advertisements.CreateAdvertisement(t, owner.Token, support.ValidAdvertisement(t))
		favorites.AddFavorite(t, user.Token, created.ID)
	}
	full := favorites.GetFavorites(t, user.Token)
	require.Equal(t, api.MaxFavorites, full.Total, "Не удалось подготовить избранное с максимальным количеством записей")

	t.Log("[Step] Пытаемся добавить новое объявление сверх лимита")
	extra := advertisements.CreateAdvertisement(t, owner.Token, support.ValidAdvertisement(t))
	actual := favorites.AddFavoriteWithError(t, user.Token, models.AddRequest{AdvertisementID: extra.ID}, http.StatusBadRequest)
	assertions.APIError(t, actual, api.LimitReached, api.FavoritesLimitMessage)
	after := favorites.GetFavorites(t, user.Token)
	require.Equal(t, api.MaxFavorites, after.Total, "Отклонённое добавление изменило количество избранных объявлений")
	assertions.FavoriteCount(t, after, extra.ID, 0)
	t.Log("[Check] Лимит соблюдён; при очистке удаляются только записи и объявления этого теста")
}
