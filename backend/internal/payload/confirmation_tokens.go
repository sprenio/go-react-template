package payload

type EmailChangePayload struct {
	NewEmail string `json:"new_email"`
}

type PasswordChangePayload struct {
	NewPassword string `json:"new_password"`
}
