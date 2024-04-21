package Model

type Message struct {
	FromUsername string `json:"from_user"`
	Message      string `json:"message" validate:"required"`
	ToUser       string `json:"to_user" validate:"required"`
	Type         string `json:"type" validate:"required,oneof=message typing"`
	ChatID       uint   `json:"chat_id"`
}
