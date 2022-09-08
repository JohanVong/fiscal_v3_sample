package requests

import (
	"github.com/JohanVong/fiscal_v3_sample/internal/models"
)

type AuthResponse struct {
	Token string      `json:"Token"`
	User  models.User `json:"User"`
}
