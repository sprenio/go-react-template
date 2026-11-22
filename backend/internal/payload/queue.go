package payload

import "backend/internal/user"

type ResetPasswordQueuePayload struct {
	Email string `json:"email"`
}
type ChangeEmailQueuePayload struct {
	UserId   uint   `json:"user_id"`
	NewEmail string `json:"new_email"`
}

type RegisterQueuePayload struct {
	user.User
	LanguageCode    string `json:"language_code"`
	NewPassword     string `json:"password"`
	NormalizedEmail string `json:"normalized_email"`
}

func (rp RegisterQueuePayload) ToUser() user.User {
	u := rp.User
	u.Password = rp.NewPassword
	u.NormalizedEmail = rp.NormalizedEmail
	return u
}
