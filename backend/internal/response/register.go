package response

import (
	"net/http"
	"context"

	"backend/internal/apicodes"
)

func SetRegisterSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	SuccessCodeResponse(w, ctx, apicodes.API_Register_Success)
}

func RegisterErrorUserNameOrEmailTaken(w http.ResponseWriter) {
	apiErrorResponse(w, http.StatusConflict, apicodes.API_Register_User_Name_Or_Email_Taken)
}
