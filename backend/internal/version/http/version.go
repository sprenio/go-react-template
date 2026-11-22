package http

import (
	"backend/internal/contexthelper"
	"backend/internal/version/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func Register(r chi.Router) {
	h := NewHandler()

	r.Get("/version", h.Version)
}

func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cfg := contexthelper.GetConfig(ctx)
	version := dto.VersionResponse{
		Sha:     cfg.CommitSha,
		Version: cfg.GitTag,
		Name:    cfg.AppName,
		Env:     cfg.AppEnv,
	}
	dto.SetVersionResponse(w, ctx, version)
}
