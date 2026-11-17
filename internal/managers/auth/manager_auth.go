package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	t.Helper()
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	body := response.Read(t, loginRawResponse, http.StatusOK)
	loginResponse := response.Decode[models.LoginOkResponse](t, body)
	require.NotEmpty(t, loginResponse.Token, "В ответе авторизации отсутствует токен")
	require.NotEmpty(t, loginResponse.User.Id, "В ответе авторизации отсутствует id пользователя")

	return loginResponse
}
