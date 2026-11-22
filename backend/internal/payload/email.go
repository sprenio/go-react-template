package payload

type ConfirmationEmailPayload struct {
	Token string `json:"token"`
}

type EmailQueuePayload struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
