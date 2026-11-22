package dto

import (
	"context"
	"net/http"

	"backend/config"
	"backend/internal/language"
	"backend/internal/response"
)

type CfgResponseData struct {
	AppName         string
	CommitSha       string
	Features        config.FeaturesConfig
	Languages       []language.Language
	DefaultLanguage string
}

func SetCfgSuccessResponse(w http.ResponseWriter, ctx context.Context, cfg *CfgResponseData) {
	response.SuccessDataResponse(w, ctx, cfg)
}
