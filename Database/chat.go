package Database

import (
	"errors"
	"log"
	"messaging/Model"

	"gorm.io/gorm"
)

type chatRepo struct{}

func ChatRepo() *chatRepo {
	return &chatRepo{}
}

func (*chatRepo) CreateChat(chat Model.Chat) (Model.Chat, error) {
	if err := DBConnection.Create(&chat).Error; err != nil {
		return Model.Chat{}, err
	}
	return chat, nil
}

func (*chatRepo) GetChat(ID uint) (Model.Chat, error) {
	// Query to retrieve user by name or email
	var chat Model.Chat
	if err := DBConnection.Where("id = ?", ID).First(&chat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Model.Chat{}, errors.New("no record found")
		} else {
			return Model.Chat{}, err
		}
	}

	log.Printf("Chats found for username:%s\n", chat.Username)
	return chat, nil
}

func (*chatRepo) GetChatDelivered(username string, delivered bool) ([]Model.Chat, error) {
	// Query to retrieve user by name or email
	var chats []Model.Chat
	if err := DBConnection.Where("delivered = ? AND username = ?", delivered, username).Order("time asc").Find(&chats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []Model.Chat{}, errors.New("no record found")
		} else {
			return []Model.Chat{}, err
		}
	}

	log.Printf("Chats found for username:%s\n", username)
	return chats, nil
}

func (*chatRepo) UpdateChat(chat Model.Chat) (Model.Chat, error) {

	// Query to update the chat
	if err := DBConnection.Save(&chat).Error; err != nil {
		log.Printf("Error updating chat with ID %d: %v\n", chat.ID, err)
		return chat, err
	}

	log.Printf("Updated the chat sucessfully for ID:%d, username:%s\n", chat.ID, chat.Username)
	return chat, nil
}
