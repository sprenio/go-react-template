package handler

import (
	"encoding/json"
	"net/http"
	"reflect"

	"backend/internal/apperrors"
	"backend/internal/contexthelper"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Language string `json:"language"`
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.MethodNotAllowedErrorResponse(w)
		return
	}
	ctx := r.Context()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.InvalidJsonErrorResponse(w)
		return
	}
	logger.InfoCtx(ctx, "RegisterHandler called with username: %s, email: %s, language: %s", req.Username, req.Email, req.Language)

	t := reflect.TypeOf(req)
	v := reflect.ValueOf(req)
	requiredFields := []string{"Username", "Email", "Password"}
	for _, fieldName := range requiredFields {
		field := v.FieldByName(fieldName)
		if field.String() == "" {
			inputField, _ := t.FieldByName(fieldName)
			inputFieldName := inputField.Tag.Get("json")
			response.InvalidInputValueErrorResponse(w, inputFieldName, inputFieldName+" field is required")
			return
		}
	}

	db := contexthelper.GetDb(ctx)
	uRepo := repository.NewUserRepository(db)
	ctRepo := repository.NewConfirmationTokenRepository(db)
	langRepo := repository.NewLanguageRepository(db)

	service := service.NewRegisterService(ctRepo, uRepo, langRepo)
	err := service.RegisterUser(ctx, req.Username, req.Email, req.Password, req.Language)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to register user: %v", err)
		if apperrors.IsRegisterUserNameOrEmailTakenError(err) {
			logger.InfoCtx(ctx, "Username or email already taken: %s, %s", req.Username, req.Email)
			response.RegisterErrorUserNameOrEmailTaken(w)
		} else if apperrors.IsAppInvalidInputError(err) {
			logger.InfoCtx(ctx, "Invalid input: %v", err)
			inputFieldName := ""
			field, found := t.FieldByName(err.(*apperrors.AppInvalidInputError).Field)
			if found {
				inputFieldName = field.Tag.Get("json")
			}
			response.InvalidInputValueErrorResponse(w, inputFieldName, err.Error())
		} else {
			response.InternalServerError(w)
		}
		return
	}
	logger.InfoCtx(ctx, "User %s (%s) registered successfully", req.Username, req.Email)
	response.SetRegisterSuccessResponse(w, ctx)
}
