package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

func TestRemoveNonexistentAdvertisementFromFavorites(t *testing.T) {
	owner, user := support.LoginUsers(t)
	id := support.DeletedAdvertisementID(t, owner.Token)
	t.Log("[Step] Удаляем из избранного заведомо отсутствующее объявление")
	actual := favorites.DeleteFavoriteWithError(t, user.Token, id, http.StatusNotFound)
	assertions.APIError(t, actual, api.NotFound, api.FavoriteNotFound)
	assertions.FavoriteCount(t, favorites.GetFavorites(t, user.Token), id, 0)
}
