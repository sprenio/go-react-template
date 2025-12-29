package payload

import "backend/internal/models"

type RegisterPayload struct {
	models.User
	LanguageId  uint8  `json:"language_id"`
	NewPassword string `json:"password"`
}

type EmailChangePayload struct {
	NewEmail  string `json:"new_email"`
}

type PasswordChangePayload struct {
	NewPassword string `json:"new_password"`
}


func (rp RegisterPayload) ToUser() models.User {
	u := rp.User
	u.Password = rp.NewPassword
	return u
}
