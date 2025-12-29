package response

import (
	"backend/internal/models"
	"net/http"
	"context"
)

type meResponseData struct {
	User  UserResponseData `json:"user"`
}


type UserResponseData struct {
	Id           uint                 `json:"id"`
	Name         string               `json:"name"`
	Email        string               `json:"email"`
	RegisteredAt string               `json:"registered_at"`
	ConfirmedAt  string               `json:"confirmed_at,omitempty"`
	Settings     models.UserSettingsData  `json:"settings,omitempty"`
}

func SetMeSuccessResponse(w http.ResponseWriter, ctx context.Context, user UserResponseData) {
	resp := meResponseData{
		User:  user,
	}
	SuccessDataResponse(w, ctx, resp)
}
