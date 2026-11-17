package myAdvertisement

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/tests/support"
)

func TestGetMyAdvertisements(t *testing.T) {
	owner := support.Login(t)
	request := support.ValidAdvertisement(t)
	created := advertisements.CreateAdvertisement(t, owner.Token, request)
	t.Log("[Step] Получаем собственные объявления по уникальному названию")
	body := myAdvertisements.GetMyAdvertisements(t, owner.Token, http.StatusOK, api.PageSize, created.Title)
	list := response.Decode[models.AdvertisementsList](t, []byte(body))
	require.Len(t, list.Items, 1, "В собственных объявлениях не найдено созданное объявление")
	require.Equal(t, created.ID, list.Items[0].ID, "В списке собственных объявлений другое объявление")
	support.RequireAdvertisementFields(t, list.Items[0], request, owner.User.Id)
}
