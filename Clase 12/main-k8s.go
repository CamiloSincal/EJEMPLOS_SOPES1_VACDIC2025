package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Clima struct {
	Municipio   string `json:"municipio"`
	Temperatura int    `json:"temperatura"`
	Humedad     int    `json:"humedad"`
	Clima       string `json:"clima"`
}

var municipios = []string{"Mixco", "Guatemala", "Villa Nueva", "Amatitlán", "Antigua"}
var climas = []string{"Soleado", "Nubslado", "Lluvioso", "Ventoso"}

func generarClima() Clima {
	return Clima{
		Municipio:   municipios[rand.Intn(len(municipios))],
		Temperatura: 15 + rand.Intn(15),
		Humedad:     50 + rand.Intn(50),
		Clima:       climas[rand.Intn(len(climas))],
	}
}

// Enviar a Kafka
func enviarKafka(msg Clima, kafkaBroker string) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaBroker, "clima", 0)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, _ := json.Marshal(msg)
	_, err = conn.WriteMessages(kafka.Message{Value: data})
	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Configuración desde variables de entorno
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	fmt.Printf("Enviando datos de clima...\n")
	fmt.Printf("Kafka broker: %s\n", kafkaBroker)

	for {
		clima := generarClima()
		data, _ := json.MarshalIndent(clima, "", "  ")
		fmt.Println("Nuevo dato:", string(data))

		fmt.Println("Enviando a Kafka")
		if err := enviarKafka(clima, kafkaBroker); err != nil {
			log.Println("Error Kafka:", err)
		} else {
			fmt.Println("Enviado a Kafka exitosamente")
		}

		time.Sleep(3 * time.Second)
	}
}