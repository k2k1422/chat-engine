package KafkaEvent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"messaging/Cache"
	"messaging/Channel"
	"messaging/Model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var KafkaBrokers string
var TopicName string
var Producer *kafka.Producer
var Consumer *kafka.Consumer

func init() {
	KafkaBrokers = os.Getenv("BOOTSTRAP_SERVER")
	TopicName = os.Getenv("TOPIC_NAME")
	log.Printf("the env:%s\n", KafkaBrokers)

	Cache.LRemove("topic", "master", TopicName)
	Cache.LPush("topic", "master", TopicName)
	log.Printf("push to redis sucessfully")

	config := &kafka.ConfigMap{
		"bootstrap.servers": KafkaBrokers,
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	}

	var err error

	// Create Kafka consumer
	Consumer, err = kafka.NewConsumer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	log.Printf("consumber group created sucessfully")

	err = Consumer.SubscribeTopics([]string{TopicName}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error subscribing to topic: %v\n", err)
		os.Exit(1)
	}

	Producer, err = kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Failed to create producer: %v\n", err)
		return
	}

	log.Printf("producer group created sucessfully")

	go ConsumeMessage()

}

func ProduceMessage(topic string, message []byte, headers []kafka.Header) error {
	// Produce a message to the specified topic
	deliveryChan := make(chan kafka.Event)
	err := Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Headers:        headers,
	}, deliveryChan)
	if err != nil {
		return err
	}

	// Wait for the delivery report to be received
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func ConsumeMessage() {
	run := true
	for run {
		ev := Consumer.Poll(1000)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))
			var msg Model.Message
			err := json.Unmarshal([]byte(e.Value), &msg)
			if err != nil {
				log.Println("Error:", err)
				return
			}
			log.Printf("message unmarshal sucessfully")
			Channel.ConsumerUnicast <- msg
			log.Printf("message sent sucessfully to the unicast channel")
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			run = false
		}

	}
}
