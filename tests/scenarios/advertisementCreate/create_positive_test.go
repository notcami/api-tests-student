package advertisementCreate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/support"
)

func TestCreateAdvertisementWithAllParameters(t *testing.T) {
	owner := support.Login(t)
	request := support.ValidAdvertisement(t)
	request.Photos = nil
	for index := 0; index < api.MaxPhotos; index++ {
		request.Photos = append(request.Photos, support.Photo(t, index))
	}
	t.Log("[Step] Создаём объявление со всеми полями и тремя фотографиями")
	created := advertisements.CreateAdvertisement(t, owner.Token, request)
	support.RequireAdvertisementFields(t, created, request, owner.User.Id)

	t.Log("[Step] Читаем созданное объявление по id")
	actual := advertisements.GetAdvertisement(t, created.ID)
	require.Equal(t, created.ID, actual.ID, "GET вернул другое объявление")
	support.RequireAdvertisementFields(t, actual, request, owner.User.Id)
	require.Equal(t, created.Photos, actual.Photos, "Фотографии изменились после создания")

	t.Log("[Step] Ищем объявление по полному названию")
	list := advertisements.GetAdvertisementsList(t, owner.Token, request.Fields["title"])
	var found []models.Advertisement
	for _, item := range list.Items {
		if item.ID == created.ID {
			found = append(found, item)
		}
	}
	require.Len(t, found, 1, "Поиск по полному названию не вернул созданное объявление ровно один раз")
	support.RequireAdvertisementFields(t, found[0], request, owner.User.Id)

	t.Log("[Step] Получаем список фотографий и скачиваем каждый файл")
	photos := advertisements.GetPhotos(t, owner.Token, created.ID)
	require.ElementsMatch(t, created.Photos, photos, "Ручка фотографий вернула другой набор файлов")
	for index, photo := range created.Photos {
		require.NotEmpty(t, photo.ID, "У фотографии %d не заполнен id", index)
		require.NotEmpty(t, photo.URL, "У фотографии %d не заполнен URL", index)
		require.Equal(t, index, photo.SortOrder, "Нарушен порядок фотографий")
		data := advertisements.DownloadPhoto(t, photo.URL)
		require.Equal(t, request.Photos[index].Data, data, "Загруженная фотография %d отличается от отправленной", index)
	}
	t.Log("[Check] Объявление доступно по id, находится в поиске; все три фотографии сохранены без изменений")
}
