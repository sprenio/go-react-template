package response

import (
	"net/http"
)

func httpErrorResponse(w http.ResponseWriter, status int) {
	apiErrorResponse(w, status, status)
}

func InternalServerError(w http.ResponseWriter) {
	httpErrorResponse(w, http.StatusInternalServerError)
}
func NotFoundErrorResponse(w http.ResponseWriter) {
	httpErrorResponse(w, http.StatusNotFound)
}
func MethodNotAllowedErrorResponse(w http.ResponseWriter) {
	httpErrorResponse(w, http.StatusMethodNotAllowed)
}
func BadRequestErrorResponse(w http.ResponseWriter) {
	httpErrorResponse(w, http.StatusBadRequest)
}
