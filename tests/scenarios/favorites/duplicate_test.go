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

func TestAddFavoriteTwice(t *testing.T) {
	owner, user := support.LoginUsers(t)
	created := advertisements.CreateAdvertisement(t, owner.Token, support.ValidAdvertisement(t))
	favorites.AddFavorite(t, user.Token, created.ID)
	before := favorites.GetFavorites(t, user.Token)

	t.Log("[Step] Повторно добавляем то же объявление в избранное")
	actual := favorites.AddFavoriteWithError(t, user.Token, models.AddRequest{AdvertisementID: created.ID}, http.StatusConflict)
	assertions.APIError(t, actual, api.AlreadyExists, api.DuplicateFavorite)
	list := favorites.GetFavorites(t, user.Token)
	assertions.FavoriteCount(t, list, created.ID, 1)
	require.Equal(t, before.Total, list.Total, "Повторное добавление изменило количество избранных объявлений")
}
