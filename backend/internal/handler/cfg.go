package handler

import (
	"backend/config"
	"backend/internal/contexthelper"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) CfgHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cfg := contexthelper.GetConfig(ctx)
	db := contexthelper.GetDb(ctx)

	langRepo := repository.NewLanguageRepository(db)
	langService := service.NewLanguageService(langRepo)

	languages, err := langService.GetLanguages(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Get languages failed: %v", err)
		response.InternalServerError(w)
		return
	}

	featuresConfig := config.FeaturesConfig{
		Register:      cfg.Register.Enabled,
		ResetPassword: cfg.ResetPassword.Enabled,
	}

	cfgData := response.CfgResponseData{
		AppName:         cfg.AppName,
		Features:        featuresConfig,
		Languages:       languages,
		DefaultLanguage: cfg.DefaultLanguage,
	}
	response.SetCfgSuccessResponse(w, ctx, &cfgData)
}
