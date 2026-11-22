package dto

import (
	"backend/internal/apicodes"
	"backend/internal/response"
	"context"
	"net/http"
)

func SetSettingsSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessNoContentResponse(w, ctx)
}

func SettingsErrorResponse(w http.ResponseWriter, ctx context.Context) {
	response.ErrorCodeResponse(w, http.StatusBadRequest, apicodes.API_Settings_Failed)
}
