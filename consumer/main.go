package main

import (
	"github.com/IBM/sarama"
	"log"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9094"}, config)
	if err != nil {
		log.Fatal("Could not create consumer: ", err.Error())
	}

	log.Println("start consuming")
	subscribe("actor", consumer)

	// block the process
	select {}
}

func subscribe(topic string, consumer sarama.Consumer) {
	partitionList, err := consumer.Partitions(topic) // get all partitions on the given topic
	if err != nil {
		log.Fatal("Error retrieving partitionList ", err.Error())
	}
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		go func() {
			if err != nil {
				log.Printf("could not create partition consumer, err: %s", err.Error())
			}

			for {
				select {
				case err := <-pc.Errors():
					log.Printf("could not process message, err: %s", err.Error())
				case message := <-pc.Messages():
					// Store message somewhere (database, S3, etc.)
					log.Printf("received message \"%s\" %v\n", string(message.Key), string(message.Value))
				}
			}
		}()
	}
}
