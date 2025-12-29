package handler

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

func (h *Handler) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req models.UserSettingsData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.InvalidJsonErrorResponse(w)
		return
	}
	logger.DebugCtx(ctx, "SettingsHandler called with data: %v", req)


	t := reflect.TypeOf(req)
	if req.AppOpt2 != nil {
		switch *req.AppOpt2 {
		case models.APP_OPTION_A, models.APP_OPTION_B, "":
			// OK
		default:
			inputField, _ := t.FieldByName("AppOpt2")
			jsonTag := inputField.Tag.Get("json")
			fieldName, _, _ := strings.Cut(jsonTag, ",")
			response.InvalidInputValueErrorResponse(w, fieldName, "")
			return
		}
	}


	usRepo := repository.NewUserSettingsRepository(h.db)
	langRepo := repository.NewLanguageRepository(h.db)

	service := service.NewUserSettingsService(usRepo, langRepo)
	err := service.Update(ctx, req)
	if err != nil {
		logger.ErrorCtx(ctx, "Update user settings failed: %v", err)
		response.InternalServerError(w)
		return
	}
	response.SetSettingsSuccessResponse(w, r.Context())
}
