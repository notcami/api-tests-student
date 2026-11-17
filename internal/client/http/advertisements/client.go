package advertisements

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements ищет объявления по переданному названию.
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	t.Helper()

	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		End().Response
}

// HttpPostAdvertisement отправляет форму создания объявления с файлами фотографий.
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateRequest) *http.Response {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	keys := make([]string, 0, len(request.Fields))
	for key := range request.Fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		require.NoError(t, writer.WriteField(key, request.Fields[key]), "Не удалось записать поле %s", key)
	}
	for _, photo := range request.Photos {
		part, err := writer.CreateFormFile("photos", photo.Name)
		require.NoError(t, err, "Не удалось добавить файл %s", photo.Name)
		_, err = part.Write(photo.Data)
		require.NoError(t, err, "Не удалось записать файл %s", photo.Name)
	}
	require.NoError(t, writer.Close(), "Не удалось завершить multipart-форму")
	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		ContentType(writer.FormDataContentType()).Body(body.String()).Expect(t).End().Response
}

// HttpGetAdvertisement запрашивает объявление по идентификатору.
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление от имени его владельца.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetPhotos запрашивает список фотографий объявления.
func HttpGetPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath + "/" + id + "/photos").
		Expect(t).End().Response
}

// HttpGetPhoto скачивает файл по URL хранилища без передачи токена API.
func HttpGetPhoto(t *testing.T, photoURL string) *http.Response {
	t.Helper()
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, photoURL, nil)
	require.NoError(t, err, "Не удалось создать запрос фотографии")
	client := &http.Client{Timeout: api.RequestTimeout}
	raw, err := client.Do(request)
	require.NoError(t, err, "Не удалось скачать фотографию %s", photoURL)
	return raw
}
