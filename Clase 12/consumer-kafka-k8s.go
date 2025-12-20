package main

import (
	"context"
	"fmt"
	"log"
	"os"
	time "time"

	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   "clima",
		GroupID: "clima-consumer-group",
	})

	fmt.Println("Esperando mensajes de Kafka...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error leyendo mensaje:", err)
			continue
		}
		fmt.Printf("[Kafka] Mensaje recibido: %s\n", string(m.Value))
		time.Sleep(2 * time.Second)
	}
}