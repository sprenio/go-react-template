package dto

import (
	"context"
	"net/http"

	"backend/internal/apicodes"
	"backend/internal/response"
)

func SetRegisterSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessNoContentResponse(w, ctx)
}

func RegisterErrorUserNameTaken(w http.ResponseWriter) {
	response.ErrorCodeResponse(w, http.StatusConflict, apicodes.API_Register_User_Name_Taken)
}
