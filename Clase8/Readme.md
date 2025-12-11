# gRPC en Golang
## 1. Instalar Protocol Buffers (protoc)
### Linux
Se abre una terminal y actualiza los paquetes del sistema:
```
sudo apt update && sudo apt upgrade -y
```
Se instalan las herramientas necesarias para compilar protoc:
```
sudo apt install -y build-essential libtool pkg-config protobuf-compiler
```
Se puede verificar la instalación ejecutando:
```
protoc --version
```
Se debería ver la versión instalada de ```protoc```.

### Windows
Descargamos los binarios de protoc en:
```
https://github.com/protocolbuffers/protobuf/releases
```
*Nota: se debe descargar el archivo con terminación win según la arqutiectura de cada computadora*

- Con el zip descargado se extraen los archivos.
- Ubicamos el folder en donde queramos.
- Abrimos la carpeta de binario y copiamos el path del protoc (aplicación)
- Abrimos las variables de entorno de windows
- En la sección de variables del sistema buscamos y seleccionamos path.
- Agregamos el path de protoc.

## 2. Instalar plugins de Go en gRPC
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
## 3. Generar código Go desde ```.proto```
El archivo .proto es el corazón de cualquier proyecto gRPC. Es un archivo de definición de contrato que describe los servicios, métodos y mensajes que usarán tanto el cliente como el servidor.

- Se coloca el archivo ```.proto``` (por ejemplo, tweet.proto) en el directorio del proyecto.

- Generamos el código necesario ejecutando:
```
protoc --go_out=. --go-grpc_out=. tweet.proto
```
- Esto generará dos archivos en /proto:

  - ```tweet.pb.go```
  - ```tweet_grpc.pb.go```

### 4. Configurar el código de Go
Copiamos el codigo de server.go y client.go
```
go mod init grpc
# luego
go mod tidy
```

### 5. Se ejecuta el proyecto
Iniciamos el servidor en una terminal:
```
go run server.go
```
En otra terminal, ejecutamos el cliente:
```
go run client.go
```

---
# Kubernetes
Kubernetes es un sistema de orquestación de contenedores que automatiza el despliegue, la gestión y la escalabilidad de aplicaciones en contenedores. A continuación se presenta una descripción general de Kubernetes, junto con algunos comandos y ejemplos útiles.

---
## Conceptos Principales
### Cluster
Un cluster es un conjunto de máquinas (nodos) que ejecutan aplicaciones en contenedores, coordinados por un nodo maestro (Control Plane).

### Node
Cada node es una máquina (física o virtual) que forma parte del clúster y ejecuta los Pods.

### Pod
El Pod es la unidad mínima de trabajo en Kubernetes. Puede contener uno o más contenedores y comparte recursos como almacenamiento y red.

### Deployment
Un Deployment es una definición que especifica cómo deben ejecutarse los Pods, gestionando réplicas, actualizaciones y fallos.

### Service
Un Service proporciona una dirección IP estable y un puerto para acceder a los Pods, incluso si cambian de IP.

### Namespace
Un Namespace permite segmentar lógicamente los recursos dentro del clúster para aislamiento y organización.

### ConfigMap y Secret
ConfigMap: Almacena configuraciones no sensibles.
Secret: Almacena datos sensibles (como contraseñas) de forma segura.

--- 
## ¿Cómo configuramos kubernetes (k8s) en GCP?
