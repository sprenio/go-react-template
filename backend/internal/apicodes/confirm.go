package apicodes

const (
	API_Confirm_Success       = 1300
	API_Confirm_Failure       = 1301
	API_Confirm_Invalid_Token = 1302
)

var confirmCodeDescriptions = map[int]string{
	API_Confirm_Success:       "Confirmation successful",
	API_Confirm_Failure:       "Confirmation failed",
	API_Confirm_Invalid_Token: "Invalid confirmation token",
}
