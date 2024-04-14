package Model

type Message struct {
	FromUsername string `json:"from_user"`
	Message      string `json:"message"`
	ToUser       string `json:"to_user"`
	Type         string `json:"type"`
}
