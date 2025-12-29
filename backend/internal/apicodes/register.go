package apicodes

const (
	API_Register_Success                  = 1100
	API_Register_User_Name_Or_Email_Taken = 1101
)

var registerCodeDescriptions = map[int]string{
	API_Register_Success:                  "Registration successful",
	API_Register_User_Name_Or_Email_Taken: "Username or email already taken",
}
