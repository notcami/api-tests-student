package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// FavoriteCount проверяет количество вхождений конкретного объявления во всём избранном.
func FavoriteCount(t *testing.T, list models.List, id string, expected int) {
	t.Helper()
	count := 0
	for _, item := range list.Items {
		if item.Advertisement.ID == id {
			count++
		}
	}
	require.Equal(t, expected, count, "Объявление %s: неверное количество записей в избранном", id)
	t.Logf("[Check] Объявление %s встречается в избранном %d раз", id, count)
}
