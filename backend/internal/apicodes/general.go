package apicodes

import (
	"maps"
	"net/http"
)

const (
	API_General_Success       = 1000
	API_General_Unknown_Error = 1001
	API_General_Invalid_JSON  = 1002

	API_General_Custom_Error        = 1051
	API_General_Invalid_Input_Value = 1052
)

// Mapa kodów odpowiedzi do opisów
var generalCodeDescriptions = map[int]string{
	API_General_Success:       "Success",
	API_General_Unknown_Error: "Unknown error",
	API_General_Invalid_JSON:  "Invalid JSON",

	API_General_Custom_Error:        "Custom error",
	API_General_Invalid_Input_Value: "Invalid input value",
}

var codeDescriptions = make(map[int]string)

func init() {
	codeMaps := []map[int]string{
		registerCodeDescriptions,
		loginCodeDescriptions,
		confirmCodeDescriptions,
		settingsCodeDescriptions,
		emailChangeCodeDescriptions,
		passwordResetCodeDescriptions,
		passwordChangeCodeDescriptions,
		// Add other code maps here

		// general errors at the end to override any duplicates
		generalCodeDescriptions,
	}

	codeDescriptions = make(map[int]string)
	for _, m := range codeMaps {
		maps.Copy(codeDescriptions, m)
	}
}

// Funkcja zwracająca opis na podstawie kodu
func GetCodeDescription(code int) string {
	if code == http.StatusOK {
		code = API_General_Success
	}
	if code < 1000 {
		// For standard HTTP status codes, use the http package to get the description

		desc := http.StatusText(code)
		if desc != "" {
			return desc
		}
	} else if desc, ok := codeDescriptions[code]; ok {
		return desc
	}
	return codeDescriptions[API_General_Unknown_Error]
}
