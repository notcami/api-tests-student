package models

type PhotoUpload struct {
	Name string
	Data []byte
}

// CreateRequest позволяет отдельно проверить отсутствующие и пустые поля формы.
type CreateRequest struct {
	Fields map[string]string
	Photos []PhotoUpload
}
