// main.rs
use axum::{
    routing::post,
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

// Respuesta que enviamos de vuelta al cliente
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
    // Configuración del enrutador con una sola ruta POST /clima
    let app = Router::new()
        .route("/clima", post(recibir_clima));

    // Configuración de la dirección del servidor
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    println!("Servidor de clima corriendo en http://{}", addr);
    println!("Endpoint disponible: POST /clima");

    // Inicio del servidor
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

// ======================================================================
// CONTROLADOR
// ======================================================================

// Handler que recibe los datos de clima desde Locust
async fn recibir_clima(Json(payload): Json<DatosClima>) -> impl IntoResponse {
    // Mostramos los datos recibidos en la consola del servidor
    println!("\nDatos de clima recibidos:");
    println!("Lugar: {}", payload.name);
    println!("Temperatura: {}°C", payload.temperatura);
    println!("Humedad: {}%", payload.humedad);
    println!("Clima: {}", payload.clima);
    
    // Creamos la respuesta
    let respuesta = Respuesta {
        mensaje: format!("Datos de {} recibidos correctamente", payload.name),
        datos_recibidos: payload,
    };
    
    // Retornamos código 200 OK con los datos en JSON
    (StatusCode::OK, Json(respuesta))
}