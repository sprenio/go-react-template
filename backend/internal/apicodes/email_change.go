package apicodes

const (
	API_Email_Change_Success             = 1500
	API_Email_Change_Email_Already_Used  = 1501
	API_Email_Change_Email_Has_No_Change = 1502
	API_Email_Change_Same_Email          = 1503
)

var emailChangeCodeDescriptions = map[int]string{
	API_Email_Change_Success:             "Email change request successful",
	API_Email_Change_Email_Already_Used:  "Email address is already in use",
	API_Email_Change_Email_Has_No_Change: "Email address has no change",
	API_Email_Change_Same_Email:          "New email is the same as the current one",
}
