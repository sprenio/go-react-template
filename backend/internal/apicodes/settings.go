package apicodes

const (
	API_Settings_Success                  = 1400
	API_Settings_Failed                   = 1401
)

var settingsCodeDescriptions = map[int]string{
	API_Settings_Success:                  "Settings updated successfully",
	API_Settings_Failed:                   "Failed to update settings",
}
