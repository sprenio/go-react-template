package response

import (
	"net/http"
	"context"

	"backend/config"
	"backend/internal/models"
)

type CfgResponseData struct {
	AppName        string
	Features       config.FeaturesConfig
	Languages      []models.Language
	DefaultLanguage string
}

func SetCfgSuccessResponse(w http.ResponseWriter, ctx context.Context, cfg *CfgResponseData) {
	SuccessDataResponse(w,ctx,  cfg)
}
