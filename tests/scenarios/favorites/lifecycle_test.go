package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/tests/support"
)

func TestFavoritesAddAndRemove(t *testing.T) {
	owner, user := support.LoginUsers(t)
	request := support.ValidAdvertisement(t)
	created := advertisements.CreateAdvertisement(t, owner.Token, request)
	before := favorites.GetFavorites(t, user.Token)
	assertions.FavoriteCount(t, before, created.ID, 0)

	t.Log("[Step] Добавляем чужое объявление в избранное")
	added := favorites.AddFavorite(t, user.Token, created.ID)
	support.RequireAdvertisementFields(t, added.Advertisement, request, owner.User.Id)
	list := favorites.GetFavorites(t, user.Token)
	assertions.FavoriteCount(t, list, created.ID, 1)
	require.Equal(t, before.Total+1, list.Total, "После добавления total должен увеличиться на один")

	t.Log("[Step] Удаляем объявление из избранного")
	favorites.DeleteFavorite(t, user.Token, created.ID)
	list = favorites.GetFavorites(t, user.Token)
	assertions.FavoriteCount(t, list, created.ID, 0)
	require.Equal(t, before.Total, list.Total, "После удаления total должен вернуться к исходному значению")

	t.Log("[Step] Повторно добавляем удалённую запись")
	favorites.AddFavorite(t, user.Token, created.ID)
	assertions.FavoriteCount(t, favorites.GetFavorites(t, user.Token), created.ID, 1)
}
