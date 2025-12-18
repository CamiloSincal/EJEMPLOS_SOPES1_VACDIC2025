package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Implementación del servicio TweetService
type tweetServer struct {
	pb.UnimplementedTweetServiceServer
}

// Método para manejar la solicitud SendTweet
func (s *tweetServer) SendTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
	// Validar que los datos no estén vacíos
	if req.Description == "" || req.Country == "" || req.Weather == "" {
		log.Printf("Datos incompletos recibidos")
		return &pb.TweetResponse{
			Status: "Error: Datos incompletos",
		}, nil
	}

	// Log detallado de la información recibida
	log.Printf("\n========================================")
	log.Printf("NUEVO TWEET RECIBIDO")
	log.Printf("========================================")
	log.Printf("País/Lugar: %s", req.Country)
	log.Printf("Clima: %s", req.Weather)
	log.Printf("Descripción: %s", req.Description)
	log.Printf("========================================\n")


	// Respuesta al cliente con un mensaje de confirmación
	responseMessage := fmt.Sprintf("Tweet de %s recibido y procesado correctamente", req.Country)
	
	return &pb.TweetResponse{
		Status: responseMessage,
	}, nil
}

func main() {
	// Obtener puerto desde variable de entorno o usar 50051 por defecto
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	// Escuchar conexiones en el puerto configurado
	address := fmt.Sprintf("0.0.0.0:%s", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error al abrir el puerto %s: %v", port, err)
	}

	// Crear un servidor gRPC con opciones
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10 * 1024 * 1024), // Tamaño máximo de mensajes recibidos = 10MB
		grpc.MaxSendMsgSize(10 * 1024 * 1024), // Tamaño máximo de mensajes enviados = 10MB
	)

	// Registrar el servicio TweetService en el servidor
	pb.RegisterTweetServiceServer(grpcServer, &tweetServer{})

	// Registrar health check service para Kubernetes
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Registrar reflection service para debugging
	reflection.Register(grpcServer)

	// Log de inicio
	log.Printf("Servidor gRPC iniciado")
	log.Printf("Escuchando en: %s", address)
	log.Printf("Health checks habilitados")
	log.Printf("Reflection habilitado para debugging")
	log.Println("Esperando conexiones...")

	// Manejo de señales para shutdown graceful
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		
		log.Println("\nSeñal de apagado recibida, cerrando servidor...")
		grpcServer.GracefulStop()
		log.Println("Servidor cerrado correctamente")
	}()

	// Iniciar el servidor y aceptar conexiones
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}