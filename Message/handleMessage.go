package Message

import (
	"log"
	"messaging/Channel"
	"messaging/Database"
	"messaging/Model"
	"time"

	"github.com/gorilla/websocket"
)

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-Channel.Broadcast

		_, err := Database.UserRepo().GetUser(msg.ToUser)
		if err != nil {
			log.Println("user not found")
			msg.Message = "User not found: " + msg.ToUser
			FromUser := Channel.Clients[msg.FromUsername]
			FromUser.WriteJSON(msg)
		}

		chat := Model.Chat{
			Message:   msg.Message,
			FromUser:  msg.FromUsername,
			Username:  msg.ToUser,
			Time:      time.Now(),
			Read:      false,
			Delivered: false,
		}

		chat, err = Database.ChatRepo().CreateChat(chat)

		if err != nil {
			log.Println("error while writing to database:", err)
		} else {
			log.Println("Sucessfully wrote the chat into db for user:", msg.FromUsername)
		}

		// Send it out to every client that is currently connected
		Channel.ClientsMutex.Lock()
		toUser, isKey := Channel.Clients[msg.ToUser]
		if isKey {
			err := toUser.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				log.Println("Connection has been closed for user's web scoket:", msg.ToUser)

				// If an error occurs while writing, close the connection
				toUser.Close()
			} else {
				// if message sent sucessfully the write into the database
				chat.Delivered = true
				chat.DeliveredTime = time.Now()
				chat, err := Database.ChatRepo().UpdateChat(chat)

				if err != nil {
					log.Printf("Failed to update the chat:{}, err:{}", chat, err)
				}
			}
		} else {
			log.Println("Not active connection to user:", msg.ToUser)
		}
		Channel.ClientsMutex.Unlock()
	}
}

func FindAndSendTheUndelivedChat(ws *websocket.Conn, username string) {
	chats, err := Database.ChatRepo().GetChatDelivered(username, false)
	if err != nil {
		log.Printf("Failed to update the username:%s, err:%v", username, err)
	}

	log.Printf("the number of undelivred char:%d", len(chats))

	for _, chat := range chats {
		err := ws.WriteJSON(Model.Message{
			FromUsername: chat.FromUser,
			Message:      chat.Message,
			ToUser:       username,
			Type:         "message",
		})

		if err != nil {
			log.Printf("error: %v", err)
			log.Println("Connection has been closed for user's web scoket:", username)
			// If an error occurs while writing, close the connection
			ws.Close()
		} else {
			// if message sent sucessfully the write into the database

			// update the delived time and value
			chat.Delivered = true
			chat.DeliveredTime = time.Now()
			chat, err := Database.ChatRepo().UpdateChat(chat)

			if err != nil {
				log.Printf("Failed to update the chat:%+v, err:%v", chat, err)
			}
		}

	}

}
