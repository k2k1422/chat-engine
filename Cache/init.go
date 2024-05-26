package Cache

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func init() {
	// Connect to Redis server
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_SERVER"),   // Redis server address
		Password: os.Getenv("REDIS_PASSWORD"), // No password set
		DB:       0,                           // Use default DB
	})

	log.Printf("redis:%s, password:%s", os.Getenv("REDIS_SERVER"), os.Getenv("REDIS_PASSWORD"))

	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Error pinging Redis:", err)
		return
	} else {
		log.Println("Redis connected sucessfully")
	}

}

func Set(prifix string, key string, value string) bool {
	// Set a key in Redis with a value
	err := Rdb.Set(context.Background(), prifix+":"+key, value, 0).Err()
	if err != nil {
		log.Println("Redis Error setting key:", err)
		return false
	}
	return true
}

func LPush(prifix string, key string, value string) error {
	// Push the value to the end of the list
	err := Rdb.RPush(context.Background(), key, value).Err()
	if err != nil {
		log.Println("Redis Error pushing value:", err)
		return err
	}
	return nil
}

func LRemove(prifix string, key string, value string) bool {
	// Remove the value from the list
	_, err := Rdb.LRem(context.Background(), key, 0, value).Result()
	if err != nil {
		log.Println("Redis Error deleting value from list:", err)
		return false
	}
	return true
}

func LGet(prifix string, key string) []string {
	list, err := Rdb.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		log.Println("Redis Error retrieving list:", err)
		return []string{}
	}
	return list
}

func LFind(prifix string, key string, value string) bool {
	// Remove the value from the list
	// Use LPOS to find the position of the value in the list
	list, err := Rdb.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		log.Println("Error finding value:", err)
		return false
	}

	log.Printf("Redis Found list of entry for key:%s, list:%+v", key, list)

	for _, v := range list {
		if v == value {
			log.Printf("Redis found value:%s in key:%s", value, key)
			return true
		}
	}
	return false
}
