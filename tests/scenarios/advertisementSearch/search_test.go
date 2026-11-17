package advertisementSearch

import (
	"fmt"
	"strings"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/support"
)

func TestSearchByTitleSubstring(t *testing.T) {
	cases := []struct {
		name  string
		query func(string) string
	}{
		{name: "beginning", query: func(id string) string { return "some" + id }},
		{name: "middle", query: func(id string) string { return "middle" + id }},
		{name: "end", query: func(id string) string { return id + "thing" }},
		{name: "lower_case", query: func(id string) string { return strings.ToLower("middle" + id) }},
		{name: "upper_case", query: func(id string) string { return strings.ToUpper("middle" + id) }},
		{name: "mixed_case", query: func(id string) string { return utils.MixCase("middle" + id) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			owner := support.Login(t)
			request := support.ValidAdvertisement(t)
			suffix := utils.RandomString(8)
			request.Fields["title"] = fmt.Sprintf("some%s middle%s %sthing", suffix, suffix, suffix)
			created := advertisements.CreateAdvertisement(t, owner.Token, request)
			query := tc.query(suffix)
			t.Logf("[Step] Ищем созданное объявление по подстроке %q", query)
			list := advertisements.GetAdvertisementsList(t, owner.Token, query)
			assertions.ContainsAdvertisementByTitle(t, list.Items, created.Title)
			t.Log("[Check] Созданное объявление найдено")
		})
	}
}
