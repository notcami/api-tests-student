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

func TestAddOwnAdvertisementToFavorites(t *testing.T) {
	owner := support.Login(t)
	created := advertisements.CreateAdvertisement(t, owner.Token, support.ValidAdvertisement(t))
	before := favorites.GetFavorites(t, owner.Token)

	t.Log("[Step] Владелец пытается добавить собственное объявление в избранное")
	actual := favorites.AddFavoriteWithError(t, owner.Token, models.AddRequest{AdvertisementID: created.ID}, http.StatusForbidden)
	assertions.APIError(t, actual, api.Forbidden, api.OwnAdvertisement)
	list := favorites.GetFavorites(t, owner.Token)
	assertions.FavoriteCount(t, list, created.ID, 0)
	require.Equal(t, before.Total, list.Total, "Отклонённый запрос изменил избранное владельца")
}
