package response

import (
	"net/http"
	"context"
	"backend/internal/apicodes"
)


func SetSettingsSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	SuccessCodeResponse(w, ctx, apicodes.API_Settings_Success)
}

func SettingsErrorResponse(w http.ResponseWriter, ctx context.Context) {
	apiErrorResponse(w, http.StatusBadRequest, apicodes.API_Settings_Failed)
}

