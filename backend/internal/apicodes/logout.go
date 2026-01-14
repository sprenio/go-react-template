package apicodes

const (
	API_Logout_Success = 1800
	API_Logout_Failed  = 1801
)

var logoutCodeDescriptions = map[int]string{
	API_Logout_Success: "Logout successful",
	API_Logout_Failed:  "Logout failed",
}
