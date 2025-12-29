package response

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend/config"
	"backend/internal/apicodes"
	"backend/internal/contexthelper"
	"backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

type apiTokenResponse interface {
	GenerateToken(ctx context.Context)
}

type msgResponse struct {
	Token   string `json:"token,omitempty"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *msgResponse) GenerateToken(ctx context.Context) {

	userID, ok := contexthelper.GetUserId(ctx)

	if !ok || userID == 0 {
		if !ok {
			logger.WarnCtx(ctx, "Failed to get user id from context")
		}
		r.Token = ""
		return
	}
	token, err := generateJWT(userID)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to generate JWT token: %v", err)
		r.Token = ""
	} else {
		r.Token = token
	}
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

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
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

func apiSuccessResponse(w http.ResponseWriter, ctx context.Context, response apiTokenResponse) {

	response.GenerateToken(ctx)
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

func generateJWT(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret)
}
