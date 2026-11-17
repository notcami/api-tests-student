package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/support"
)

func TestAddNonexistentAdvertisementToFavorites(t *testing.T) {
	owner, user := support.LoginUsers(t)
	id := support.DeletedAdvertisementID(t, owner.Token)
	t.Log("[Step] Добавляем в избранное заведомо отсутствующее объявление")
	actual := favorites.AddFavoriteWithError(t, user.Token, models.AddRequest{AdvertisementID: id}, http.StatusNotFound)
	assertions.APIError(t, actual, api.NotFound, api.AdvertisementNotFound)
	assertions.FavoriteCount(t, favorites.GetFavorites(t, user.Token), id, 0)
}
