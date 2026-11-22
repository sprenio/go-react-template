package http

import (
	"backend/internal/helper"
	"backend/internal/response"
	"backend/internal/settings"
	"backend/internal/settings/dto"
	"backend/internal/settings/service"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) SetSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := helper.GetRequestData[settings.UserSettingsData](r, w)
	if err != nil {
		response.InvalidInputResponse(w, err)
		return
	}
	service := service.NewService(ctx)
	err = service.Update(ctx, req)
	if err != nil {
		logger.ErrorCtx(ctx, "Update user settings failed", "error", err)
		dto.SettingsErrorResponse(w, ctx)
		return
	}
	dto.SetSettingsSuccessResponse(w, r.Context())
}
