package Message

import (
	"encoding/json"
	"log"
	"messaging/Cache"
	"messaging/Channel"
	"messaging/Database"
	"messaging/KafkaEvent"
	"messaging/Model"
	"time"

	"github.com/gorilla/websocket"
)

func HandleUnicastProducerMessage() {
	for {
		msg := <-Channel.Broadcast

		_, err := Database.UserRepo().GetUser(msg.ToUser)
		if err != nil {
			log.Println("user not found")
			msg.Message = "User not found: " + msg.ToUser
			Channel.WriteJSONWS(msg, msg.FromUsername)
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

		msg.ChatID = chat.ID

		servers := Cache.LGet("topic", "master")
		log.Printf("list of servers exist:%+v", servers)
		for _, server := range servers {

			// send kafka event to all the to user which is obiovs also send to from users
			if Cache.LFind("connection", server, msg.ToUser) || Cache.LFind("connection", server, msg.FromUsername) {
				log.Printf("send kafka message to server:%s", server)
				jsonBytes, err := json.Marshal(msg)
				err = KafkaEvent.ProduceMessage(server, []byte(string(jsonBytes)), nil)
				if err != nil {
					log.Printf("Failed to send the message into the server:%s", server)
				} else {
					log.Printf("Sent message sucessfully to the server:%s", server)
				}
			} else {
				log.Printf("For user:%s, there are no active web scoket connection", msg.ToUser)
			}
		}
	}
}

func HandleUnicastConsumerMessage() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-Channel.ConsumerUnicast

		_, err := Database.UserRepo().GetUser(msg.ToUser)
		if err != nil {
			log.Println("user not found")
			msg.Message = "User not found: " + msg.ToUser
			Channel.WriteJSONWS(msg, msg.FromUsername)
		}

		chat, err := Database.ChatRepo().GetChat(msg.ChatID)

		if err != nil {
			log.Println("error while writing to database:", err)
		} else {
			log.Printf("Sucessfully while fetching the user chat:%+v", chat)
		}

		// Send it out to every client that is currently connected
		_, isKey := Channel.Clients[msg.ToUser]
		if isKey {
			err := Channel.WriteJSONWS(msg, msg.ToUser)
			if err != nil {
				log.Printf("error: %v", err)
				log.Println("Connection has been closed for user's web scoket:", msg.ToUser)
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

		// send messages to self
		_, isKey = Channel.Clients[msg.FromUsername]
		if isKey {
			err := Channel.WriteJSONWS(msg, msg.FromUsername)
			if err != nil {
				log.Printf("error: %v", err)
				log.Println("Connection has been closed for user's web scoket:", msg.ToUser)
			} else {

				if err != nil {
					log.Printf("Failed to update the chat:{}, err:{}", chat, err)
				}
			}
		} else {
			log.Println("Not active connection to user:", msg.ToUser)
		}
	}
}

func DeleteKeyCacheIfNotConnected() {
	for {
		time.Sleep(2 * time.Second)
		sessions := Cache.LGet("connection", KafkaEvent.TopicName)
		for _, session := range sessions {
			_, isKey := Channel.Clients[session]
			if !isKey {
				Cache.LRemove("connection", KafkaEvent.TopicName, session)
				log.Printf("Removed the key has the connection expired:%s", session)
			}
		}

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
			ChatID:       chat.ID,
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
