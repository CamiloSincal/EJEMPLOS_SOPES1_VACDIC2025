// main.rs
use axum::{
    routing::{get, post},
    Json,
    Router,
    http::StatusCode,
    response::IntoResponse,
};
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;

// ======================================================================
// ESTRUCTURAS DE DATOS
// ======================================================================

// Estructura que representa los datos de clima que envía Locust
#[derive(Serialize, Deserialize, Debug, Clone)]
struct DatosClima {
    name: String,           // Nombre del lugar (guatemala, mexico, etc.)
    temperatura: i32,       // Temperatura en grados (18-28)
    humedad: i32,          // Porcentaje de humedad (40-80)
    clima: String,         // Tipo de clima (soleado, nublado, lluvioso)
}

// Respuesta que enviamos de vuelta a Locust
#[derive(Serialize)]
struct Respuesta {
    mensaje: String,
    datos_recibidos: DatosClima,
}

// ======================================================================
// FUNCIÓN PRINCIPAL
// ======================================================================

#[tokio::main]
async fn main() {
    // Configuración del enrutador
    let app = Router::new()
        .route("/clima", post(recibir_y_reenviar_clima))
        .route("/health", get(health_check));

    // Configuración de la dirección del servidor
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    println!("Servidor de clima corriendo en http://{}", addr);
    println!("Endpoint disponible: POST /clima");
    println!("Esperando datos de Locust para reenviar al servicio Go...");

    // Inicio del servidor
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

// ======================================================================
// CONTROLADORES
// ======================================================================

// Health check endpoint para k8s
async fn health_check() -> StatusCode {
    StatusCode::OK
}

// Handler que recibe datos de Locust y los reenvía al servicio Go
async fn recibir_y_reenviar_clima(Json(payload): Json<DatosClima>) -> impl IntoResponse {
    // Mostramos los datos recibidos de Locust
    println!("\n=== Datos recibidos de Locust ===");
    println!("Lugar: {}", payload.name);
    println!("Temperatura: {}°C", payload.temperatura);
    println!("Humedad: {}%", payload.humedad);
    println!("Clima: {}", payload.clima);
    
    // Obtener URL del servicio Go desde variable de entorno
    let go_service_url = std::env::var("GO_SERVICE_URL")
        .unwrap_or_else(|_| "http://localhost:8080".to_string());
    
    // Crear cliente HTTP y reenviar los datos al servicio Go
    let client = reqwest::Client::new();
    
    println!("\n=== Reenviando datos al servicio Go ===");
    println!("URL destino: {}/clima", go_service_url);
    
    match client
        .post(&format!("{}/clima", go_service_url))
        .json(&payload)
        .send()
        .await
    {
        Ok(response) => {
            if response.status().is_success() {
                println!("Datos reenviados exitosamente al servicio Go");
            } else {
                println!("Error al reenviar datos: {}", response.status());
            }
        }
        Err(e) => {
            println!("Error de conexión con servicio Go: {}", e);
        }
    }
}