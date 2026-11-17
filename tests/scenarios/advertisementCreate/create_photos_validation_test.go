package advertisementCreate

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/support"
)

func TestCreateAdvertisementWithInvalidPhotoCount(t *testing.T) {
	cases := []struct {
		name    string
		count   int
		message string
	}{
		{name: "no_photos", count: 0, message: "at least 1 photo is required"},
		{name: "too_many_photos", count: api.MaxPhotos + 1, message: "maximum 3 photos allowed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			owner := support.Login(t)
			request := support.ValidAdvertisement(t)
			request.Photos = make([]models.PhotoUpload, 0, tc.count)
			for index := 0; index < tc.count; index++ {
				request.Photos = append(request.Photos, support.Photo(t, index))
			}
			t.Logf("[Step] Создаём объявление с количеством фотографий: %d", tc.count)
			actual := advertisements.CreateAdvertisementWithError(t, owner.Token, request, http.StatusBadRequest)
			assertions.APIError(t, actual, api.InvalidRequest, tc.message)
		})
	}
}
