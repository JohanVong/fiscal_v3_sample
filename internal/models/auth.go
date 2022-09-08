package models

// swagger:parameters apiAuth
type LoginForm struct {
	// Входные данные для получения токена
	// in: body
	Body struct {
		// Required: true
		Username string `json:"Login"`
		// Required: true
		Password string `json:"Password"`
	}
}
