package response

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/apicodes"
	"backend/internal/contexthelper"
	"backend/internal/cookie"
	"backend/internal/repository"
	"backend/pkg/logger"
)

type msgResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	msgResponse
	Description string `json:"description,omitempty"`
}
type successResponse struct {
	msgResponse
	Data any `json:"data,omitempty"`
}
type invalidInputValueResponse struct {
	errorResponse
	InvalidField string `json:"invalid_field,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, ctx context.Context) {
	SuccessCodeResponse(w, ctx, apicodes.API_General_Success)
}

func SuccessCodeResponse(w http.ResponseWriter, ctx context.Context, code int) {
	response := &successResponse{
		msgResponse: getMsgResponse(code),
	}
	apiSuccessResponse(w, ctx, response)
}
func SuccessMessageResponse(w http.ResponseWriter, ctx context.Context, message string) {
	response := &msgResponse{
		Code:    apicodes.API_General_Success,
		Message: message,
	}
	apiSuccessResponse(w, ctx, response)
}
func SuccessDataResponse(w http.ResponseWriter, ctx context.Context, data any) {
	response := &successResponse{
		msgResponse: getMsgResponse(apicodes.API_General_Success),
		Data:        data,
	}
	apiSuccessResponse(w, ctx, response)
}
func SuccessDataCodeResponse(w http.ResponseWriter, ctx context.Context, data any, code int) {
	response := &successResponse{
		msgResponse: getMsgResponse(code),
		Data:        data,
	}
	apiSuccessResponse(w, ctx, response)
}

// error responses
func InvalidJsonErrorResponse(w http.ResponseWriter) {
	apiErrorResponse(w, http.StatusBadRequest, apicodes.API_General_Invalid_JSON)
}

func InvalidInputValueErrorResponse(w http.ResponseWriter, field string, desc string) {
	response := invalidInputValueResponse{
		errorResponse: errorResponse{
			msgResponse: getMsgResponse(apicodes.API_General_Invalid_Input_Value),
			Description: desc,
		},
		InvalidField: field,
	}
	apiResponse(w, http.StatusBadRequest, response)
}

func UnauthorizedErrorResponse(w http.ResponseWriter, description string) {
	apiErrorWithDescriptionResponse(w, http.StatusUnauthorized, http.StatusUnauthorized, description)
}

// private
func apiResponse(w http.ResponseWriter, status int, response any) {
	setJsonResponseHeaders(w, status)
	json.NewEncoder(w).Encode(response)
}

func apiSuccessResponse(w http.ResponseWriter, ctx context.Context, response any) {
	setUserAccessCookies(ctx, w)
	apiResponse(w, http.StatusOK, response)
}
func apiErrorResponse(w http.ResponseWriter, status int, code int) {
	response := getErrorResponseByCode(code, "")
	apiResponse(w, status, response)
}

func apiErrorWithDescriptionResponse(w http.ResponseWriter, status int, code int, description string) {
	response := getErrorResponseByCode(code, description)
	apiResponse(w, status, response)
}

func getErrorResponseByCode(code int, description string) errorResponse {
	return errorResponse{
		msgResponse: getMsgResponse(code),
		Description: description,
	}
}
func setJsonResponseHeaders(w http.ResponseWriter, status int) {
	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func getMsgResponse(code int) msgResponse {
	return msgResponse{
		Code:    code,
		Message: apicodes.GetCodeDescription(code),
	}
}

func setUserAccessCookies(ctx context.Context, w http.ResponseWriter) {
	accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
	logger.DebugCtx(ctx, "AccessTokenData: %v", accessTokenData)
	if !accessTokenData.SetCookies {
		return
	}
	if accessTokenData.UserId > 0 {
		cookie.SetAccessToken(ctx, w, accessTokenData.UserId)
	}
	if accessTokenData.RefreshToken != "" {
		db := contexthelper.GetDb(ctx)
		sessionRepo := repository.NewUserSessionsRepository(db)

		// Calculate TTL for refresh token
		cfg := contexthelper.GetConfig(ctx)
		ttl := time.Hour * 24 * time.Duration(int64(cfg.Token.RefreshTokenTtlDays))
		expiresAt := time.Now().Add(ttl)
		ttlSeconds := int(ttl / time.Second)

		// Refresh token in database
		err := sessionRepo.RefreshToken(ctx, accessTokenData.UserId, accessTokenData.RefreshToken, expiresAt)
		if err != nil {
			logger.ErrorCtx(ctx, "Refresh token failed: %v", err)
			cookie.DeleteRefreshToken(ctx, w)
			return
		}

		// Set refresh token cookie
		cookie.SetRefreshToken(w, ctx, accessTokenData.RefreshToken, ttlSeconds)
	}

}
