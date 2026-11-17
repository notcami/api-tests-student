package advertisementCreate

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/api"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/support"
)

func TestCreateAdvertisementWithInvalidFields(t *testing.T) {
	cases := []struct {
		name    string
		field   string
		value   string
		omit    bool
		message string
	}{
		{name: "missing_title", field: "title", omit: true, message: api.TitleRequired},
		{name: "empty_title", field: "title", message: api.TitleRequired},
		{name: "long_title", field: "title", value: strings.Repeat("a", api.MaxTitleLength+1),
			message: fmt.Sprintf("title must not exceed %d characters provided %d", api.MaxTitleLength, api.MaxTitleLength+1)},
		{name: "missing_description", field: "description", omit: true, message: api.DescriptionRequired},
		{name: "empty_description", field: "description", message: api.DescriptionRequired},
		{name: "long_description", field: "description", value: strings.Repeat("d", api.MaxDescriptionLength+1),
			message: fmt.Sprintf("description must be between 1 and %d characters", api.MaxDescriptionLength)},
		{name: "negative_price", field: "price", value: "-1", message: api.PriceOutOfRange},
		{name: "price_above_maximum", field: "price", value: "1000000.01", message: api.PriceOutOfRange},
		{name: "non_numeric_price", field: "price", value: "abc", message: "cannot convert string value abc to float"},
		{name: "zero_quantity", field: "quantity", value: "0", message: api.QuantityOutOfRange},
		{name: "negative_quantity", field: "quantity", value: "-1", message: api.QuantityOutOfRange},
		{name: "quantity_above_maximum", field: "quantity", value: strconv.Itoa(api.MaxQuantity + 1), message: api.QuantityOutOfRange},
		{name: "non_numeric_quantity", field: "quantity", value: "abc", message: "cannot convert string value to int"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			owner := support.Login(t)
			request := support.ValidAdvertisement(t)
			if tc.omit {
				delete(request.Fields, tc.field)
			} else {
				request.Fields[tc.field] = tc.value
			}
			t.Logf("[Step] Отправляем форму с некорректным полем %s (%s)", tc.field, tc.name)
			actual := advertisements.CreateAdvertisementWithError(t, owner.Token, request, http.StatusBadRequest)
			assertions.APIError(t, actual, api.InvalidRequest, tc.message)
		})
	}
}
