package support

import "api-tests-template/internal/constants/api"

type AuthorizationCase struct {
	Name    string
	Token   string
	Message string
}

// AuthorizationCases возвращает варианты отсутствующей и некорректной авторизации.
func AuthorizationCases() []AuthorizationCase {
	return []AuthorizationCase{
		{Name: "missing_token", Message: api.AuthorizationRequired},
		{Name: "invalid_token", Token: api.InvalidToken, Message: api.InvalidAuthorizationMessage},
	}
}
