package apicodes

const (
	API_Login_Success             = 1200
	API_Login_Invalid_Credentials = 1201
)

var loginCodeDescriptions = map[int]string{
	API_Login_Success:             "Login successful",
	API_Login_Invalid_Credentials: "Invalid credentials",
}
