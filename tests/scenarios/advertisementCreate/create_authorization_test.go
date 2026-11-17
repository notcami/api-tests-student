package advertisementCreate

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/support"
)

func TestCreateAdvertisementWithoutAuthorization(t *testing.T) {
	utils.LoadEnv()
	for _, tc := range support.AuthorizationCases() {
		t.Run(tc.Name, func(t *testing.T) {
			t.Log("[Precondition] Подготавливаем корректную форму объявления")
			request := support.ValidAdvertisement(t)
			t.Logf("[Step] Создаём объявление: %s", tc.Name)
			actual := advertisements.CreateAdvertisementWithError(t, tc.Token, request, http.StatusUnauthorized)
			assertions.APIError(t, actual, api.Unauthorized, tc.Message)
		})
	}
}
