package support

import (
	"os"
	"sync"
	"testing"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
)

type sessionKey struct {
	apiURL   string
	login    string
	password string
}

var (
	sessionsMu sync.Mutex
	sessions   = make(map[sessionKey]models.LoginOkResponse)
)

// loginUser переиспользует авторизацию в рамках пакета, чтобы не перегружать ручку входа.
func loginUser(t *testing.T, login, password string) models.LoginOkResponse {
	t.Helper()
	key := sessionKey{apiURL: os.Getenv("API_URL"), login: login, password: password}
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	if session, ok := sessions[key]; ok {
		t.Log("[Precondition] Используем токен, полученный для этого пакета")
		return session
	}
	session := auth.Login(t, login, password)
	sessions[key] = session
	return session
}
