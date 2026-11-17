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

func TestAddFavoriteWithInvalidID(t *testing.T) {
	cases := []struct {
		name    string
		request any
		message string
	}{
		{name: "missing_id", request: map[string]string{},
			message: "Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag"},
		{name: "empty_id", request: models.AddRequest{}, message: "invalid UUID length: 0"},
		{name: "invalid_id", request: models.AddRequest{AdvertisementID: api.InvalidID}, message: "invalid UUID length: 10"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user := support.Login(t)
			t.Logf("[Step] Добавляем в избранное с некорректным advertisement_id: %s", tc.name)
			actual := favorites.AddFavoriteWithError(t, user.Token, tc.request, http.StatusBadRequest)
			assertions.APIError(t, actual, api.InvalidRequest, tc.message)
		})
	}
}
