package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"backend/pkg/apperror"
	"backend/pkg/logger"

	"errors"

	"github.com/go-chi/chi/v5"
)

// Interfejs dla struktur, które potrafią same się zwalidować/sformatować
type Validatable interface {
	Validate(ctx context.Context) error
}

var (
	ErrPayloadTooLarge = errors.New("payload too large")
	ErrInvalidJSON     = errors.New("invalid json format")
)

func GetDataFromRequest[T any](r *http.Request, w http.ResponseWriter, target *T) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	if err := json.NewDecoder(tee).Decode(target); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			logger.ErrorCtx(r.Context(), "Payload too large", "error", err)
			return ErrPayloadTooLarge
		}
		logger.ErrorCtx(r.Context(), "Invalid Json", "body", buf.String(), "error", err)
		return ErrInvalidJSON
	}
	return nil
}

func GetIdFromRequest(r *http.Request) (uint, error) {
	idParam := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id64 < 1 {
		ctx := r.Context()
		logger.ErrorCtx(ctx, "Error parsing id parameter", "id", idParam, "error", err)
		return 0, apperror.Wrap(err, "Error parsing id parameter")
	}
	return uint(id64), nil
}

func GetStringIdFromRequest(r *http.Request) (string, uint, error) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		logger.ErrorCtx(ctx, "Error parsing id parameter: empty")
		return "", 0, errors.New("Error parsing id parameter")
	}
	parts := strings.Split(idParam, "_")
	if len(parts) != 2 {
		logger.ErrorCtx(ctx, "Error parsing id parameter: invalid format")
		return "", 0, errors.New("Error parsing id parameter")
	}
	val, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		logger.ErrorCtx(ctx, "Error parsing id parameter: invalid format")
		return "", 0, errors.New("Error parsing id parameter")
	}

	return parts[0], uint(val), nil
}

func GetBoolFromRequestFormValue(r *http.Request, param string) (bool, error) {
	return strconv.ParseBool(r.FormValue(param))
}

// T musi być typem, który po wskaźniku (*T) implementuje Validatable
func GetRequestData[T any, PT interface {
	*T
	Validatable
}](r *http.Request, w http.ResponseWriter) (T, error) {
	ctx := r.Context()
	var req T

	if err := GetDataFromRequest(r, w, &req); err != nil {
		logger.ErrorCtx(ctx, "Error getting request data", "error", err)
		return req, err
	}

	// Wywołujemy walidację (i ewentualne formatowanie) dokładnie raz:
	err := PT(&req).Validate(ctx)

	return req, err
}
