package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/support"
)

func TestFavoritesWithoutAuthorization(t *testing.T) {
	for _, tc := range support.AuthorizationCases() {
		for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodDelete} {
			t.Run(tc.Name+"/"+method, func(t *testing.T) {
				owner, user := support.LoginUsers(t)
				created := advertisements.CreateAdvertisement(t, owner.Token, support.ValidAdvertisement(t))
				expectedCount := 0
				if method == http.MethodDelete {
					favorites.AddFavorite(t, user.Token, created.ID)
					expectedCount = 1
				}
				before := favorites.GetFavorites(t, user.Token)
				t.Logf("[Step] Выполняем %s для избранного: %s", method, tc.Name)
				var actual response.Error
				switch method {
				case http.MethodGet:
					actual = favorites.GetFavoritesWithError(t, tc.Token, http.StatusUnauthorized)
				case http.MethodPost:
					actual = favorites.AddFavoriteWithError(t, tc.Token, models.AddRequest{AdvertisementID: created.ID}, http.StatusUnauthorized)
				case http.MethodDelete:
					actual = favorites.DeleteFavoriteWithError(t, tc.Token, created.ID, http.StatusUnauthorized)
				}
				assertions.APIError(t, actual, api.Unauthorized, tc.Message)
				after := favorites.GetFavorites(t, user.Token)
				assertions.FavoriteCount(t, after, created.ID, expectedCount)
				require.Equal(t, before.Total, after.Total, "Неавторизованный запрос изменил избранное пользователя")
			})
		}
	}
}
