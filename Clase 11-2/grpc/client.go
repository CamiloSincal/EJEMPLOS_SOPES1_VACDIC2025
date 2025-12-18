package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"
	pb "grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Estructura para recibir datos de Rust
type DatosClima struct {
	Name        string `json:"name"`
	Temperatura int32  `json:"temperatura"`
	Humedad     int32  `json:"humedad"`
	Clima       string `json:"clima"`
}

// Cliente gRPC global
var grpcClient pb.TweetServiceClient

func main() {
	// Configurar conexión gRPC
	grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")
	if grpcServerAddr == "" {
		grpcServerAddr = "localhost:50051"
	}

	conn, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient = pb.NewTweetServiceClient(conn)
	log.Printf("Conectado al servidor gRPC en %s", grpcServerAddr)

	// Configurar servidor HTTP REST
	http.HandleFunc("/clima", recibirClimaHandler)
	http.HandleFunc("/health", healthCheckHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor REST corriendo en puerto %s", port)
	log.Printf("Esperando datos de Rust...")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error al iniciar servidor HTTP: %v", err)
	}
}

// Handler para recibir datos de clima desde Rust
func recibirClimaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var datos DatosClima
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		log.Printf("Error al decodificar JSON: %v", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	log.Printf("\n--- Datos recibidos de Rust ---")
	log.Printf("Lugar: %s", datos.Name)
	log.Printf("Temperatura: %d°C", datos.Temperatura)
	log.Printf("Humedad: %d%%", datos.Humedad)
	log.Printf("Clima: %s", datos.Clima)

	// Enviar datos al servidor gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	description := formatDescription(datos)
	
	resp, err := grpcClient.SendTweet(ctx, &pb.TweetRequest{
		Description: description,
		Country:     datos.Name,
		Weather:     datos.Clima,
	})

	if err != nil {
		log.Printf("Error al enviar a gRPC: %v", err)
		http.Error(w, "Error al procesar datos", http.StatusInternalServerError)
		return
	}

	log.Printf("Respuesta gRPC: %s", resp.Status)

	// Responder a Rust
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Datos procesados y enviados a gRPC",
	})
}

// Health check para Kubernetes
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Formatear descripción para el tweet
func formatDescription(datos DatosClima) string {
	return fmt.Sprintf("Clima en %s: %s, Temperatura: %d°C, Humedad: %d%%",
		datos.Name, datos.Clima, datos.Temperatura, datos.Humedad)
}