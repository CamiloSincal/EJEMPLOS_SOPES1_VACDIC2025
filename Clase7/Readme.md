# Crear una imagen de Rust
Primero tendremos que hacer el build de nuestra aplicación de Rust:
```bash
cargo run
```

Con la aplicación complicada colocamos nuestro Dockerfile en el mismo directorio o nivel que nuestro Cargo.toml

En este directorio colocaremos nuestro Dockerfile:
```Dockerfile
# ============================================
# ETAPA 1: BUILD
# ============================================
FROM rust:1.75-slim as builder

WORKDIR /app

# Copiamos solo los archivos de dependencias primero
# Esto permite cachear las dependencias si no cambian
COPY Cargo.toml Cargo.lock ./

# Creamos un src dummy para compilar solo las dependencias
RUN mkdir src && \
    echo "fn main() {}" > src/main.rs && \
    cargo build --release && \
    rm -rf src

# Ahora copiamos el código real
COPY . .

# Forzamos recompilación del código (no las dependencias)
RUN touch src/main.rs && \
    cargo build --release

# ============================================
# ETAPA 2: RUNTIME (Imagen final)
# ============================================
FROM debian:bookworm-slim

# Instalamos solo las librerías runtime necesarias
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    libssl3 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copiamos SOLO el binario compilado desde la etapa de build
COPY --from=builder /app/target/release/nombre-aplicacion /app/nombre-aplicacion

# Creamos un usuario no-root por seguridad
RUN useradd -m -u 1001 appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

# Mismo que el de Cargo.toml
CMD ["./nombre-aplicacion"] 
```
Finalmente podemos construir la imagen y ejecutar el contenedor con:
```bash
# Construir la imagen
docker build -t api-rust .

# Ejecutar el contenedor
docker run -d -p 3000:3000 api-rust

# Ver logs del contenedor
docker logs -f id-mi-api-rust
```


# Locust
Locust es una herramienta de código abierto escrita en Python para realizar pruebas de carga y rendimiento

## Cómo empezar - Instalación
La instalación es sencilla y requiere solo unos pasos:
1. **Requisito:** Tener Python instalado.
2. **Instalación:** Ejecuta en tu terminal o consola:
```bash
# Instalar Locust
pip install locust

#Verificar instalación
locust --version
```
## Prueba simple
El núcleo de Locust es un script de Python. Este es un ejemplo básico:
```py
from locust import HttpUser, task,between

class myUser(HttpUser):
    wait_time = between(1,3)

    @task
    def home_page(self):
        self.client.get("https://www.google.com")
```

Explicación rápida:
- ```HttpUser```: Clase que representa a un usuario virtual que hace peticiones HTTP.
- ```@task```: Decorador que marca el método que el usuario ejecutará repetidamente.
- ```self.client```: Es el cliente HTTP utilizado para hacer las solicitudes.
  
### Ejecutar la prueba simple
Tienes dos formas principales de ejecutar tus pruebas:
1. Con Interfaz Web (Recomendado para empezar):
- Se ejecuta en la terminal: ```locust```
- Se puede verificar abriendo el navegador en: ```http://localhost:8089```

Es necesario configurar:
- Número de usuarios totales a simular.
- Tasa de generación (usuarios por segundo).
- Host (la URL de tu aplicación a probar).
1. Por Línea de Comando (Headless):
- Útil para automatización o CI/CD.
```bash
locust --headless --users 100 --spawn-rate 10 -H http://tu-app.com
```
*Nota: No nombrar los archivos como locust.py para evitar conflictos con el comando para ejecutarlo*

## Prueba de ejemplo para proyecto
```py
from locust import HttpUser, TaskSet, task, between
import random
import json

class MyTasks(TaskSet):
    
    @task(1)
    def engineering(self):
        # Listado aleatorio de nombres
        names = ["guatemala", "mexico", "panama", "inglaterra", "francia", "italia", "españa", "argentina", "chile", "colombia"]

        climas = ["soleado", "nublado", "lluvioso"]
    
        # Datos de estudiantes
        weather_data = {
            "name": random.choice(names),  # Random name
            "temperatura": random.randint(18, 28),  # Temperatura random entre 18 y 28
            "humedad": random.randint(40, 80),  # Humedad random entre 40 y 80
            "clima": random.choice(climas)  # Clima aleatorio
        }
        
        # Envio de JSON hacia /engineering route como POST
        headers = {'Content-Type': 'application/json'}
        self.client.post("/clima", json=weather_data, headers=headers)

    

class WebsiteUser(HttpUser):
    tasks = [MyTasks]
    wait_time = between(1, 5)  # Tiempo de espera entre tareas entre 1 y 5 segundos
```