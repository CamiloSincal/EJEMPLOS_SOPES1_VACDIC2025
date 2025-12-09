// ======================================================================
// SECCIÓN 1: IMPORTACIÓN DE BIBLIOTECAS
// ======================================================================

// Importamos componentes del framework web Axum
use axum::{
    // `get` y `post` son funciones para definir rutas HTTP (GET y POST)
    routing::{get, post},
    // `Json` permite enviar y recibir datos en formato JSON
    Json,
    // `Router` es el enrutador principal que maneja las rutas de la API
    Router,
    // `StatusCode` contiene códigos HTTP como 200 OK, 404 NOT FOUND, 201 CREATED, etc.
    http::StatusCode,
    // `IntoResponse` permite que diferentes tipos de datos se conviertan en respuestas HTTP
    response::IntoResponse,
};
// Importamos Serde para serialización/deserialización (convertir objetos Rust a JSON y viceversa)
use serde::{Deserialize, Serialize};
// Importamos herramientas para redes del estándar de Rust (direcciones IP y puertos)
use std::net::SocketAddr;

// ======================================================================
// SECCIÓN 2: DEFINICIÓN DE ESTRUCTURAS DE DATOS
// ======================================================================

// Definimos una estructura (struct) llamada 'Usuario' que representa un usuario en nuestro sistema.
// #[derive(...)] son "atributos" que añaden funcionalidad automáticamente:
// - `Serialize`: Permite convertir esta struct a JSON cuando se envía en una respuesta.
// - `Deserialize`: Permite crear esta struct desde JSON cuando se recibe en una petición (ej., en un POST).
// - `Clone`: Permite crear copias explícitas del objeto con el método `.clone()`.
#[derive(Serialize, Deserialize, Clone)]
struct Usuario {
    id: u32,        // Entero sin signo de 32 bits (solo números positivos). Identificador único.
    nombre: String, // Cadena de texto (String) para el nombre. String es dinámica y se almacena en el heap.
    email: String,  // Cadena de texto (String) para el email.
}

// Definimos el estado compartido (state) de la aplicación.
// En una aplicación real, aquí irían conexiones a bases de datos, clientes de caché, etc.
// Para este ejemplo, simulamos una "base de datos" en memoria con un simple vector.
// `Clone` es necesario para poder pasar una copia del estado a cada manejador (handler) de ruta.
#[derive(Clone)]
struct AppState {
    // `Vec` es un vector/array dinámico y redimensionable (como un ArrayList en Java o un Array en JavaScript).
    // `Vec<Usuario>` significa "un vector que contiene elementos del tipo `Usuario`".
    usuarios: Vec<Usuario>,
}

// ======================================================================
// SECCIÓN 3: FUNCIÓN PRINCIPAL - PUNTO DE ENTRADA DEL PROGRAMA
// ======================================================================

// El atributo `#[tokio::main]` transforma la función `main` en una función asíncrona
// y configura un "runtime" (sistema de ejecución) para manejar operaciones concurrentes.
// Es similar a `async` en JavaScript pero con un sistema más explícito y de alto rendimiento.
#[tokio::main]
async fn main() {
    // ==================================================================
    // 3.1: INICIALIZACIÓN DEL ESTADO (BASE DE DATOS EN MEMORIA)
    // ==================================================================

    // Creamos el estado inicial de la aplicación con algunos usuarios de ejemplo.
    // En un proyecto real, estos datos se cargarían desde una base de datos persistente.
    let state = AppState {
        // `vec!` es un macro conveniente de Rust para crear un vector con valores iniciales.
        usuarios: vec![
            Usuario {
                id: 1,
                // `.to_string()` convierte un "string literal" (tipo `&str`) a un `String` (cadena en el heap).
                nombre: "Ana".to_string(),
                email: "ana@email.com".to_string(),
            },
            Usuario {
                id: 2,
                nombre: "Luis".to_string(),
                email: "luis@email.com".to_string(),
            },
        ],
    };

    // ==================================================================
    // 3.2: CONFIGURACIÓN DEL ENRUTADOR (DEFINICIÓN DE RUTAS)
    // ==================================================================

    // Creamos el enrutador (Router) principal. Es el corazón de nuestra API.
    // Asigna direcciones URL (rutas) y métodos HTTP (GET, POST) a funciones manejadoras (handlers).
    let app = Router::new()
        // Ruta: `GET /` -> Llama a la función `raiz()` cuando alguien accede a la página principal.
        .route("/", get(raiz))
        // Ruta: `GET /usuarios` -> Llama a `obtener_usuarios()` para listar todos los usuarios.
        .route("/usuarios", get(obtener_usuarios))
        // Ruta: `POST /usuarios` -> Llama a `crear_usuario()` para agregar un nuevo usuario.
        .route("/usuarios", post(crear_usuario))
        // Ruta: `GET /saludar/:nombre` -> `:nombre` es un parámetro de ruta variable.
        // Ejemplo: `/saludar/Maria` extraerá "Maria" y llamará a `saludar("Maria")`.
        .route("/saludar/:nombre", get(saludar))
        // Compartimos el estado (`state`) con todas las rutas definidas arriba.
        // Esto permite que cada función manejadora acceda a los datos de la aplicación.
        .with_state(state);

    // ==================================================================
    // 3.3: CONFIGURACIÓN DE LA DIRECCIÓN DEL SERVIDOR
    // ==================================================================

    // Creamos una `SocketAddr`, que es una combinación de dirección IP y puerto.
    // `[127, 0, 0, 1]` es la dirección IP de localhost (tu propia máquina).
    // `3000` es el número de puerto donde escuchará el servidor.
    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    // Mostramos un mensaje en la terminal para saber dónde se está ejecutando el servidor.
    println!("Servidor corriendo en http://{}", addr);

    // ==================================================================
    // 3.4: INICIO DEL SERVIDOR WEB (CÓDIGO CORREGIDO)
    // ==================================================================

    // **CORRECCIÓN APLICADA:**
    // En versiones recientes de Axum, la forma de iniciar el servidor cambió.
    // Ya no se usa `axum::Server::bind`. En su lugar:
    // 1. Primero creamos un "oyente" (listener) TCP en la dirección especificada.
    //    `tokio::net::TcpListener::bind(&addr).await` intenta vincularse al puerto.
    //    `.unwrap()` detiene el programa si falla (ej., si el puerto 3000 ya está en uso).
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();

    // 2. Luego, pasamos ese "oyente" y nuestra aplicación (`app`) a `axum::serve`.
    //    `axum::serve(listener, app)` inicia el servidor HTTP y comienza a atender peticiones.
    //    `.await.unwrap()` hace que el programa espere aquí a que el servidor se ejecute.
    axum::serve(listener, app).await.unwrap();
}

// ======================================================================
// SECCIÓN 4: CONTROLADORES (HANDLERS) - FUNCIONES QUE MANEJAN LAS RUTAS
// ======================================================================
// Todas estas funciones son `async` (asíncronas). Esto significa que pueden "esperar"
// (`await`) a que se completen operaciones (como leer de una base de datos) sin
// bloquear el hilo de ejecución, permitiendo atender muchas peticiones a la vez.

// ----------------------------------------------------------------------
// 4.1: CONTROLADOR PARA LA RUTA RAÍZ (GET /)
// ----------------------------------------------------------------------
// Esta función se llama cuando un cliente hace una petición `GET` a la raíz (`/`).
// Retorna: Una referencia a un string estático (`&'static str`).
// `'static` significa que el string vive durante toda la ejecución del programa.
async fn raiz() -> &'static str {
    // Este mensaje se envía directamente como cuerpo de la respuesta HTTP.
    "¡Hola desde la API de Rust!"
}

// ----------------------------------------------------------------------
// 4.2: CONTROLADOR PARA OBTENER USUARIOS (GET /usuarios)
// ----------------------------------------------------------------------
// Parámetro: `state` - Recibe una referencia al estado compartido de la aplicación.
// Retorna: Un `Json` que contiene un `Vec<Usuario>` (lista de usuarios).
// Axum se encarga automáticamente de convertir el `Vec` en JSON para la respuesta.
async fn obtener_usuarios(state: axum::extract::State<AppState>) -> Json<Vec<Usuario>> {
    // `state.usuarios.clone()` crea una copia completa del vector de usuarios.
    // Esto es necesario porque `Json(...)` tomará posesión (ownership) de los datos
    // que le pasemos, y no podemos "mover" los datos originales fuera del estado.
    // Envolvemos la copia en `Json(...)` para indicar que la respuesta debe ser JSON.
    Json(state.usuarios.clone())
}

// ----------------------------------------------------------------------
// 4.3: CONTROLADOR PARA CREAR USUARIO (POST /usuarios)
// ----------------------------------------------------------------------
// Parámetro 1: `state` - Una referencia al estado compartido de la aplicación.
// Parámetro 2: `Json(payload)` - Axum extrae el cuerpo JSON de la petición POST
//             y trata de convertirlo (deserializarlo) en un valor del tipo `Usuario`.
//             Si el JSON no tiene la forma correcta, Axum devuelve un error automáticamente.
// Retorna: `impl IntoResponse` - Esto significa "cualquier tipo que se pueda convertir
//          en una respuesta HTTP". Da flexibilidad para devolver diferentes cosas.
async fn crear_usuario(
    state: axum::extract::State<AppState>,
    Json(payload): Json<Usuario>,
) -> impl IntoResponse {
    // Creamos una copia mutable del estado. `mut` significa que podemos modificarlo.
    let mut new_state = state.clone();
    // Agregamos el nuevo usuario (`payload`) al final del vector.
    new_state.usuarios.push(payload);

    // Devolvemos una tupla con dos elementos:
    // 1. `StatusCode::CREATED` - El código de estado HTTP 201 (Creado), indicando éxito en la creación.
    // 2. `Json(new_state.usuarios.clone())` - El cuerpo de la respuesta.
    //    Debemos clonar `new_state.usuarios` por la misma razón que en `obtener_usuarios`:
    //    no podemos mover el vector original fuera de `new_state`. `.clone()` resuelve esto
    //    creando una nueva copia de los datos para la respuesta.
    (StatusCode::CREATED, Json(new_state.usuarios.clone()))
}

// ----------------------------------------------------------------------
// 4.4: CONTROLADOR CON PARÁMETRO DINÁMICO (GET /saludar/:nombre)
// ----------------------------------------------------------------------
// Parámetro: `axum::extract::Path(nombre)` - Axum extrae el segmento `:nombre` de la URL
//             y lo intenta convertir a un `String`.
// Retorna: Un `String` que será el cuerpo de la respuesta HTTP.
async fn saludar(axum::extract::Path(nombre): axum::extract::Path<String>) -> String {
    // `format!` es un macro similar a `println!`, pero en lugar de imprimir,
    // construye y devuelve un nuevo `String` con el valor interpolado.
    format!("¡Hola, {}! 👋", nombre)
}