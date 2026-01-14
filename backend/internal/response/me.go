package response

import (
	"backend/internal/models"
	"context"
	"net/http"
)

type meResponseData struct {
	User  models.UserResponseData `json:"user"`
}

func SetMeSuccessResponse(w http.ResponseWriter, ctx context.Context, user models.UserResponseData) {
	resp := meResponseData{
		User:  user,
	}
	SuccessDataResponse(w, ctx, resp)
}
