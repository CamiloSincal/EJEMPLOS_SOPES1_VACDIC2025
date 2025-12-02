# Instalación y configuración de Docker
---
## Instalación de Docker Linux
Para instalar Docker en Linux, se recomienda seguir la documentación oficial de Docker, ya que la instalación puede variar dependiendo de la distribución de Linux que se esté utilizando.

### Instalación de Docker en Ubuntu (oficial)
1. Desinstala versiones antiguas (si existen):

```
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
```

2. Instala dependencias y agrega el repositorio oficial:
```
# Agregar llave GPG oficial de Docker:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Agregar el repositorio a los recursos Apt:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
```

3.  Instala Docker Engine y complementos:
```
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```
4. Verificación de la instalación
```
sudo docker run hello-world
```
Si ves un mensaje de bienvenida, Docker está correctamente instalado.

### Explicación del Docker Engine
#### Arquitectura del Docker Engine
El Docker Engine es una aplicación cliente-servidor con estos componentes:

Un servidor que es un tipo de demonio que se ejecuta en la máquina host.
Una API REST que especifica interfaces que los programas pueden usar para hablar con el demonio y darle instrucciones.
#### ¿Qué es un Deamon (demonio)?

Un demonio es un programa que se ejecuta en segundo plano, sin interacción directa con el usuario. Los demonios se utilizan para realizar tareas de mantenimiento y administración del sistema, como la gestión de servicios, la programación de tareas y la monitorización del sistema.

En el caso de Docker, el demonio es el servidor de Docker que se ejecuta en la máquina host y se encarga de gestionar los contenedores y las imágenes de Docker.