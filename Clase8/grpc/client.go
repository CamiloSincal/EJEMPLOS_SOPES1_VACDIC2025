package main

import (
	"context"
	"log"
	"time"
	"os"

	pb "grpc/proto" // Importa el paquete generado para el servicio gRPC

	"google.golang.org/grpc" // Importa la librería gRPC de Google
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddr := os.Getenv("GRPC_SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051" // Valor por defecto para desarrollo local
	}
	
	// Conectar con el servidor gRPC en localhost en el puerto 50051
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close() // Asegura que la conexión se cierre al finalizar

	// Crear un cliente para el servicio TweetService definido en el proto
	client := pb.NewTweetServiceClient(conn)

	// Crear un contexto con un tiempo límite de 1 segundo para la solicitud
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Asegura que el contexto se cancele al finalizar

	// Enviar una solicitud al servidor con los datos del tweet
	resp, err := client.SendTweet(ctx, &pb.TweetRequest{
		Description: "Hola desde Go con gRPC!", // Descripción del tweet
		Country:     "Guatemala",               // País desde donde se envía
		Weather:     "Soleado",                 // Clima asociado al tweet
	})
	if err != nil {
		log.Fatalf("Error al enviar tweet: %v", err) // Manejo de error si la solicitud falla
	}

	// Imprimir la respuesta del servidor
	log.Printf("Respuesta del servidor: %s", resp.Status)
}