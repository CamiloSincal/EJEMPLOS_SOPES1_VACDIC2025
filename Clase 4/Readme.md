# Ejemplo de Grafana+SQLite
Este ejemplo hace lo siguiente:
- Crea una base de datos SQLite en la misma carpeta del programa.
- Define una tabla containers para guardar la información.
- Cada 20 segundos, genera hasta 4 contenedores con datos aleatorios (nombre, CPU, memoria, estado) y los guarda (Esto es solo para simular lo que ustedes deben hacer con /proc/continfo).
- Se ejecuta como daemon en segundo plano.

## ¿Qué es SQLite?
Es una biblioteca de C que implementa un motor de *base de datos relacional ligero y autónomo*, que se integra directamente en la aplicación y almacena todos los datos en un único archivo.  A diferencia de otros sistemas de bases de datos que requieren un servidor separado, SQLite funciona como una biblioteca dentro del proceso de la aplicación.

## Compilación
```bash
go mod init daemon-test #  Crea un nuevo proyecto Go llamado "daemon-test"
go get github.com/mattn/go-sqlite3 # Descarga e instala una librería para usar SQLite
go build -o daemon main.go # Compila el archivo main.go y crea un ejecutable llamado daemon
```

## Ejecución manual
```bash
./daemon &
```
(el ```&``` lo manda a segundo plano)

Ver logs en tiempo real:
```bash
tail -f nohup.out
```
Detenerlo:
```bash
pkill daemon
```
## Ejecución manual
Crea el archivo de servicio:

```bash
sudo nano /etc/systemd/system/grafana-db-daemon.service
```

Contenido:
```
[Unit]
Description=Daemon en Go - (descripcion corta)
After=network.target

[Service]
ExecStart=/ruta/a/tu/daemon
WorkingDirectory=/ruta/a/tu/carpeta
Restart=always

[Install]
WantedBy=multi-user.target
```
*Nota: El **WorkingDirectory** es el  directorio de trabajo desde el cual se ejecutará el programa*

Guardar y habilitar:
```
sudo systemctl daemon-reload
sudo systemctl enable grafana-db-daemon
sudo systemctl start grafana-db-daemon
```

Ver estado:
```
sudo systemctl status grafana-db-daemon
```

Parar y quitar el daemon
```
sudo systemctl stop grafana-db-daemon
sudo systemctl disable grafana-db-daemon
sudo rm /etc/systemd/system/grafana-db-daemon.service
```

## Grafana
Para este ejemplo, se crea una carpeta con el nombre **grafana** donde se tendran 2 archivos:
- Dockerfile con SQLite plugin
- Docker-compose

### Dockerfile
Grafana por defecto no soporta SQLite como datasource, pero hay un plugin de la comunidad llamado ```frser-sqlite-datasource```.

### Docker-compose
```
version: "3.9"

services:
  grafana:
    build: .
    container_name: grafana-sqlite
    ports:
      - "3000:3000"
	user: "0" 
    volumes:
      - ./grafana-data:/var/lib/grafana   # Datos de grafana
      - ../containers.db:/db/containers.db  # Montamos la DB creada por GO
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=frser-sqlite-datasource
```

Con los archivos configurados procedemos a levantar grafana:
```
docker-compose up -d
```
Grafana estará en: http://localhost:3000 Usuario: admin,Password: admin

Luego configuramos SQLite en el source de grafana:
- En Grafana, ve a Connections → Data Sources → Add new data source.
- Busca SQLite (plugin frser-sqlite-datasource).
  Si el plugin no se instalo es posible ingresar al CLI de grafana con:
  ```bash
  docker exec -it <container_name_or_id> grafana-cli <command>
  ```
  Para luego instalar con el comando de la página de plugins de grafana.
- Configura el path al archivo dentro del contenedor:
```
/db/containers.db
```

Finalmente ya solo es necesario crear el dashboard con los datos del source:
```sql
SELECT 
  name,
  cpu,
  memory,
  status,
  created_at / 1000 as time  -- ← DIVIDIR entre 1000 para que Grafana lo tome en segundos
FROM containers
ORDER BY memory DESC
LIMIT 10
```

# Ejemplo de grafana con lectura de proc
Este ejemplo hace lo siguiente:
- Crea una base de datos SQLite en la misma carpeta del programa.
- Define una tabla de procesos para guardar la información.
- Se carga el módulo de kernel
- Se hace la lectura de proc
- Se ejecuta como daemon en segundo plano.