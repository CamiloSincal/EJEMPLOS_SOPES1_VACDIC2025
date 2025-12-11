package main

import (
	"context"
	"log"
	"net"

	pb "grpc/proto" // Importa el paquete generado para el servicio gRPC

	"google.golang.org/grpc" // Importa la librería gRPC de Google
)

// Implementación del servicio TweetService
type tweetServer struct {
	pb.UnimplementedTweetServiceServer // Estructura base para implementar el servicio
}

// Método para manejar la solicitud SendTweet
func (s *tweetServer) SendTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
	// Log de la información recibida en la solicitud
	log.Printf("Nuevo Tweet: %s desde %s con clima %s", req.Description, req.Country, req.Weather)
	// Respuesta al cliente con un mensaje de confirmación
	return &pb.TweetResponse{Status: "Tweet recibido correctamente"}, nil
}

func main() {
	// Escuchar conexiones en el puerto 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al abrir el puerto: %v", err) // Manejo de error si no se puede abrir el puerto
	}

	// Crear un servidor gRPC
	grpcServer := grpc.NewServer()
	// Registrar el servicio TweetService en el servidor
	pb.RegisterTweetServiceServer(grpcServer, &tweetServer{})

	// Log para indicar que el servidor está corriendo
	log.Println("Servidor gRPC corriendo en el puerto 50051")
	// Iniciar el servidor y aceptar conexiones
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error en el servidor: %v", err) // Manejo de error si el servidor falla
	}
}