# Comandos para usar Zot como registro de contenedores

1. **Iniciar el registro Zot en una VM con DOCKER**

    Ejecuta el siguiente comando para iniciar un registro Zot en segundo plano, exponiendo el puerto 5000:
    ```sh
    docker run -d -p 5000:5000 --name zot ghcr.io/project-zot/zot-linux-amd64:latest
    ```
    Esto descargará la imagen de Zot y la ejecutará como un contenedor llamado `zot`.

    
## En la Computadora donde se tenga el desarrollo de la imagen:
2. **Edita la configuración de Docker para que pueda subir la imagen a la VM con DOCKER**

   a. **Editar la configuración de Docker**:
    ```bash
    sudo nano /etc/docker/daemon.json
    ```

   b. **Agregar o modificar el contenido del archivo** (si está vacío, agregar lo siguiente):
    ```json
    {
        "insecure-registries": ["<IP_VM_DOCKER>:5000"]
    }
    ```

   c. **Reiniciar Docker** para aplicar los cambios:
    ```bash
    sudo systemctl restart docker
    ```

3. **Crear y etiquetar la imagen para el registro privado**
   
   Para subir la imagen a docker, primero es necesario crear la imagen a través del comando:
   ```bash
    docker build -t api-go .
    ```
    Recordando siempre estar en la carpeta del ```Dockerfile```

    Usaremos de ejemplo la imagen llamada api-go.
    Cambia la etiqueta de la imagen para que apunte a tu registro privado (reemplaza `<IP_VM_DOCKER>` por la IP de tu VM):
    ```sh
    docker tag api-go:v1 <IP_VM_DOCKER>:5000/api-go:v1
    ```

1. **Subir la imagen al registro Zot**

    Sube la imagen etiquetada a tu registro Zot:
    ```sh
    docker push <IP_VM_DOCKER>:5000/api-go:v1
    ```

2. **Verificar las imágenes disponibles en el registro**

    Consulta el catálogo de imágenes almacenadas en el registro Zot:
    ```sh
    curl http://<IP_VM_DOCKER>:5000/v2/_catalog
    ```
   ***También puedes pegar la URL en tu navegador para verificar que funciona***

## En las maquinas virtuales con unicamente containerd y ctr
6. **Descargar la imagen desde el registro Zot**

    Descarga la imagen desde tu registro privado para comprobar que está disponible:
    ```sh
    sudo ctr images pull --plain-http <IP_VM_DOCKER>:5000/api-go:v1
    ```
    Lista las imagenes
    ```sh
    sudo ctr images ls

    # si el comando anterior no funciona prueba con este
    sudo ctr images list
    ```

    Listo ya puedes conectarte con tu Registro de Contenedores privados de ZOT en tu maquina virtual, ahora pongamoslo a prueba con containerd.

# Containerd
**containerd** es un runtime de contenedores de nivel industrial que maneja el ciclo de vida completo de los contenedores en un sistema host. Es el componente principal que Docker utiliza internamente, pero también puede usarse de forma independiente.

Características principales de containerd:

- **Runtime estándar:** Implementa las especificaciones OCI Runtime y Image
- **Gestión de imágenes:** Descarga, almacena y gestiona imágenes de contenedores
- **Ciclo de vida:** Crea, ejecuta, detiene y elimina contenedores
- **Snapshots:** Maneja sistemas de archivos en capas para los contenedores
- **Networking:** Proporciona capacidades básicas de red para contenedores

**ctr** es la herramienta de línea de comandos que viene con containerd, similar a como docker es el cliente para Docker Engine.

## Comandos básicos
### 1.Para verificar si containerd está corriendo
```bash
sudo systemctl status containerd
```
```bash
ctr --version
```

### 2. Para listar las imagenes disponibles
```bash
sudo ctr images ls

# si el comando anterior no funciona prueba con 
sudo ctr images list
```

### 3. Para descargar una imagen desde un registry (como Docker Hub)
```bash
sudo ctr images pull docker.io/library/hello-world:latest
```
Otra opción es con un registro privado como Zot en una maquina virtual
```bash
sudo ctr images pull --plain-http <IP_VM1_DOCKER>:5000/api-go:v1
```

### 4. Para levantar el contenedor a partir de la imagen
```bash
sudo ctr run -t --rm docker.io/library/hello-world:latest my-hello
```
De este comando es necesario saber que:
- ```-t```: modo interactivo
- ```--rm```: elimina el contenedor al salir
- ```my-hello```: ID local para el contenedor, no puede conteneder dos puntos

### 4.1. Para ejecutar el contenedor en segundo plano con la red de host
Para ejecutar un contenedor en segundo plano (modo detached) y que utilice la red del host, usa el siguiente comando. Debes reemplazar ```<IP_VM_DOCKER>``` por la IP correspondiente de tu máquina virtual o servidor donde está el registro de imágenes:
```bash
sudo ctr run -d --net-host <IP_VM_DOCKER>:5000/api-go:v1 my-api-go
```

De este comando es necesario saber que:
- ```-d```: ejecuta el contenedor en segundo plano (detached).
- ```--tty```: asigna una terminal al contenedor.
- ```--plain-http```: utiliza la red del host, permitiendo que el contenedor acceda directamente a los puertos y servicios del host.
- ```<IP_VM_DOCKER>:5000/api-go:v1```: imagen a ejecutar, obtenida desde el registro privado.
- ```my-api-go```: nombre local para el contenedor.

### 5. Para listar los contenedores activos
```bash
sudo ctr containers ls
```

### 6. Para crear y ejecutar un contenedor en pasos separados
```bash
sudo ctr images pull docker.io/library/alpine:latest
sudo ctr containers create docker.io/library/alpine:latest my-alpine
sudo ctr tasks start -d my-alpine
```

### 7.  Para subir una imagen a un registry (como Zot o Docker Hub)
```bash
sudo ctr images tag docker.io/library/hello-world:latest localhost:5000/hello-world:latest
sudo ctr images push <IP_VM1_DOCKER>:5000/hello-world:latest
```

### 8. Para eliminar una imagen
```bash
sudo ctr images rm docker.io/library/hello-world:latest
```

### 9. Para eliminar un contenedor
Primero detenemos la tarea del contenedor con:
```bash
sudo ctr tasks list #para obtener el nombre de la tarea
sudo ctr task kill <nombre-de-tarea> #para detener la tarea
sudo ctr task kill --signal SIGKILL <nombre-de-tarea> #para obligar a detener la tarea
```

Con la tarea detenida ya podemos eliminar el contenedor:
```bash
sudo ctr containers delete my-hello
```

### 10. Para acceder a un shell en un contenedor
Si el contenedor está corriendo y tiene /bin/sh:

```bash
sudo ctr tasks exec -t --exec-id myexecid my-alpine /bin/sh
```

### 11. Para eliminar todo
```bash
sudo ctr tasks list -q | xargs -I {} sudo ctr task kill {}
sudo ctr containers list -q | xargs -I {} sudo ctr containers delete {}
sudo ctr images list -q | xargs -I {} sudo ctr images remove {}
```
