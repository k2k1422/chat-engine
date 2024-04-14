package Database

import (
	"log"
)

func Ping() {
	err := Connection.Ping()
	if err != nil {
		log.Println("failed to connect")
		panic("failed to connect")
	}
	log.Println("ping sucessfull")
}
