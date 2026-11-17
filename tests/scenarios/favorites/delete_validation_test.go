package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

func TestDeleteFavoriteWithInvalidID(t *testing.T) {
	user := support.Login(t)
	t.Log("[Step] Удаляем из избранного по строке, не являющейся UUID")
	actual := favorites.DeleteFavoriteWithError(t, user.Token, api.InvalidID, http.StatusBadRequest)
	assertions.APIError(t, actual, api.InvalidRequest, "Invalid advertisement ID")
}
