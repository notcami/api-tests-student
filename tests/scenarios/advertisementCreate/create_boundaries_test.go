package advertisementCreate

import (
	"strconv"
	"strings"
	"testing"

	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/support"
)

func TestCreateAdvertisementAtValidBoundaries(t *testing.T) {
	for _, price := range []string{"0", strconv.Itoa(api.MaxPrice)} {
		t.Run("price_"+price, func(t *testing.T) {
			owner := support.Login(t)
			request := support.ValidAdvertisement(t)
			suffix := utils.RandomString(16)
			request.Fields["title"] = strings.Repeat("a", api.MaxTitleLength-len(suffix)) + suffix
			request.Fields["description"] = strings.Repeat("d", api.MaxDescriptionLength)
			request.Fields["price"] = price
			request.Fields["quantity"] = strconv.Itoa(api.MaxQuantity)
			t.Log("[Step] Создаём объявление с граничными допустимыми значениями")
			created := advertisements.CreateAdvertisement(t, owner.Token, request)
			actual := advertisements.GetAdvertisement(t, created.ID)
			support.RequireAdvertisementFields(t, actual, request, owner.User.Id)
		})
	}
}
