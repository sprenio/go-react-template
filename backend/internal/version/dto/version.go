package dto

import (
	"context"
	"net/http"

	"backend/internal/response"
)

type VersionResponse struct {
	Sha     string `json:"sha"`
	Version string `json:"version"`
	Name    string `json:"name"`
	Env     string `json:"env"`
}

func SetVersionResponse(w http.ResponseWriter, ctx context.Context, res VersionResponse) {
	response.SuccessDataResponse(w, ctx, res)
}
