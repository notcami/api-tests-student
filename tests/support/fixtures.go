package support

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	defaultPrice       = "1250.50"
	defaultQuantity    = "3"
	defaultDescription = "Настольная лампа, исправна. В комплекте кабель и инструкция."
	suffixLength       = 16
)

// Login загружает настройки и авторизует основного тестового пользователя.
func Login(t *testing.T) authModels.LoginOkResponse {
	t.Helper()
	utils.LoadEnv()
	t.Log("[Precondition] Авторизуем основного пользователя")
	return loginUser(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

// LoginUsers авторизует владельца объявлений и другого пользователя для избранного.
func LoginUsers(t *testing.T) (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	t.Helper()
	owner := Login(t)
	require.NotEmpty(t, os.Getenv("TEST_SECOND_LOGIN"), "Задайте TEST_SECOND_LOGIN в .env.override: нужен второй пользователь TestBoard")
	require.NotEmpty(t, os.Getenv("TEST_SECOND_PASSWORD"), "Задайте TEST_SECOND_PASSWORD в .env.override")
	t.Log("[Precondition] Авторизуем второго пользователя")
	second := loginUser(t, os.Getenv("TEST_SECOND_LOGIN"), os.Getenv("TEST_SECOND_PASSWORD"))
	require.NotEqual(t, owner.User.Id, second.User.Id, "Для избранного нужны два разных пользователя")
	return owner, second
}

// ValidAdvertisement возвращает новую форму со всеми полями и одной PNG-фотографией.
func ValidAdvertisement(t *testing.T) models.CreateRequest {
	t.Helper()
	return models.CreateRequest{
		Fields: map[string]string{
			"title":       "Лампа " + utils.RandomString(suffixLength),
			"description": defaultDescription,
			"price":       defaultPrice,
			"quantity":    defaultQuantity,
		},
		Photos: []models.PhotoUpload{Photo(t, 0)},
	}
}

// Photo создаёт небольшой валидный PNG; индекс задаёт отличающийся цвет файла.
func Photo(t *testing.T, index int) models.PhotoUpload {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	palette := []color.RGBA{
		{R: 0, G: 120, B: 200, A: 255},
		{R: 60, G: 120, B: 200, A: 255},
		{R: 120, G: 120, B: 200, A: 255},
		{R: 180, G: 120, B: 200, A: 255},
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, palette[index%len(palette)])
		}
	}
	var buffer bytes.Buffer
	require.NoError(t, png.Encode(&buffer, img), "Не удалось подготовить PNG-фотографию")
	return models.PhotoUpload{Name: fmt.Sprintf("lamp-%d.png", index), Data: buffer.Bytes()}
}

// DeletedAdvertisementID возвращает валидный id созданного и затем удалённого объявления.
func DeletedAdvertisementID(t *testing.T, token string) string {
	t.Helper()
	t.Log("[Precondition] Создаём и удаляем объявление для проверки отсутствующего ad_id")
	created := advertisements.CreateAdvertisement(t, token, ValidAdvertisement(t))
	advertisements.DeleteAdvertisement(t, token, created.ID)
	actual := advertisements.GetAdvertisementWithError(t, created.ID, http.StatusNotFound)
	assertions.APIError(t, actual, api.NotFound, api.AdvertisementNotFound)
	return created.ID
}

// RequireAdvertisementFields сравнивает сохранённые поля с отправленной формой.
func RequireAdvertisementFields(t *testing.T, actual models.Advertisement, request models.CreateRequest, ownerID string) {
	t.Helper()
	price, err := strconv.ParseFloat(request.Fields["price"], 64)
	require.NoError(t, err, "Некорректная цена в тестовых данных")
	quantity, err := strconv.Atoi(request.Fields["quantity"])
	require.NoError(t, err, "Некорректное количество в тестовых данных")
	require.NotEmpty(t, actual.ID, "Не заполнен id объявления")
	require.Equal(t, ownerID, actual.UserID, "Объявление принадлежит другому пользователю")
	require.Equal(t, request.Fields["title"], actual.Title, "Сохранено другое название")
	require.Equal(t, request.Fields["description"], actual.Description, "Сохранено другое описание")
	require.Equal(t, price, actual.Price, "Сохранена другая цена")
	require.Equal(t, quantity, actual.Quantity, "Сохранено другое количество")
	require.NotEmpty(t, actual.CreatedAt, "Не заполнено время создания")
	require.NotEmpty(t, actual.UpdatedAt, "Не заполнено время обновления")
	require.Len(t, actual.Photos, len(request.Photos), "Не совпадает количество фотографий")
	t.Logf("[Check] Все поля объявления %s совпадают с отправленными", actual.ID)
}
