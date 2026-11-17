package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

func TestDeletingAdvertisementRemovesFavorite(t *testing.T) {
	owner, user := support.LoginUsers(t)
	created := advertisements.CreateAdvertisement(t, owner.Token, support.ValidAdvertisement(t))
	before := favorites.GetFavorites(t, user.Token)
	favorites.AddFavorite(t, user.Token, created.ID)
	assertions.FavoriteCount(t, favorites.GetFavorites(t, user.Token), created.ID, 1)

	t.Log("[Step] Владелец удаляет объявление, добавленное другим пользователем в избранное")
	advertisements.DeleteAdvertisement(t, owner.Token, created.ID)
	actual := advertisements.GetAdvertisementWithError(t, created.ID, http.StatusNotFound)
	assertions.APIError(t, actual, api.NotFound, api.AdvertisementNotFound)

	t.Log("[Step] Проверяем избранное после удаления объявления")
	list := favorites.GetFavorites(t, user.Token)
	assertions.FavoriteCount(t, list, created.ID, 0)
	require.Equal(t, before.Total, list.Total, "Удалённое объявление продолжает учитываться в total избранного")
}
