package http

import (
	"backend/config"
	"backend/internal/contexthelper"
	"backend/internal/language/service"
	"backend/internal/response"
	"backend/pkg/logger"
	"net/http"

	"backend/internal/cfg/dto"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cfg := contexthelper.GetConfig(ctx)
	langService := service.NewService(ctx)

	languages, err := langService.Get(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Get languages failed", "error", err)
		response.InternalServerError(w)
		return
	}

	featuresConfig := config.FeaturesConfig{
		Register:      cfg.Register.Enabled,
		ResetPassword: cfg.ResetPassword.Enabled,
	}

	cfgData := dto.CfgResponseData{
		AppName:         cfg.AppName,
		CommitSha:       cfg.CommitSha,
		Features:        featuresConfig,
		Languages:       languages,
		DefaultLanguage: cfg.DefaultLanguage,
	}
	dto.SetCfgSuccessResponse(w, ctx, &cfgData)
}
