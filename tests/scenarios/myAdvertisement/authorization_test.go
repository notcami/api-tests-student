package myAdvertisement

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/support"
)

func TestGetMyAdvertisementsWithoutAuthorization(t *testing.T) {
	utils.LoadEnv()
	for _, tc := range support.AuthorizationCases() {
		t.Run(tc.Name, func(t *testing.T) {
			t.Logf("[Step] Получаем собственные объявления: %s", tc.Name)
			body := myAdvertisements.GetMyAdvertisements(t, tc.Token, http.StatusUnauthorized, api.PageSize, "")
			actual := response.Decode[response.Error](t, []byte(body))
			assertions.APIError(t, actual, api.Unauthorized, tc.Message)
		})
	}
}
