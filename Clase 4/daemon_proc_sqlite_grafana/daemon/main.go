package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
)



func main() {
	// Inicializar aleatoriedad
	rand.Seed(time.Now().UnixNano())

	// Obtener directorio del ejecutable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Error obteniendo ruta del ejecutable:", err)
	}
	exeDir := filepath.Dir(exePath)
	
	// Crear ruta absoluta para la DB
	dbPath := filepath.Join(exeDir, "containers.db")
	log.Println("Usando base de datos en:", dbPath)

	// IMPORTANTE: Eliminar si existe como directorio
	if info, err := os.Stat(dbPath); err == nil && info.IsDir() {
		log.Println("containers.db es un directorio, eliminando...")
		os.RemoveAll(dbPath)
	}

	// Asegurar que el archivo se cree 
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Creando archivo de base de datos...")
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal("Error creando archivo DB:", err)
		}
		file.Close()
		log.Println("Archivo containers.db creado")
	}

	// Conexión a SQLite
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Error abriendo base de datos:", err)
	}
	defer db.Close()

	// Crear tabla si no existe
	createTable := `
	CREATE TABLE IF NOT EXISTS registros (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		total_ram REAL,
		ram_libre REAL,
		total_procesos REAL,
		created_at INTEGER
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Error creando tabla:", err)
	}

	// ====================================================  1. EJECUCIÓN DEL DOCKER COMPOSE DE GRAFANA ====================================================
	ejecutarDockerCompose()

	// ==================================================== 2. CREACIÓN DEL CRONJOB A PARTIR DEL .SH ====================================================

	// ==================================================== 3. CARGA DE MODULOS DE KERNEL ====================================================


	// ==================================================== 4. INICIO DEL LOOP PRINCIPAL ====================================================

	// Canal para capturar señales (detener daemon)
	// 'sigs' es un canal que se utiliza para recibir señales del sistema operativo.
	// En este caso, se capturan señales SIGINT (Ctrl+C) y SIGTERM (terminación del proceso).
	// Estas señales permiten detener el daemon de manera controlada.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Se define el ticker de 20 segundos
	ticker := time.NewTicker(20 * time.Second)
	// 'defer ticker.Stop()' asegura que el temporizador se detenga correctamente cuando el programa termine.
	defer ticker.Stop()

	// Este bucle infinito permite que el daemon esté en ejecución constante,
	// esperando eventos del ticker o señales del sistema operativo.
loop:
	// Se lee el archivo de proc cada 20 segundos a través de un for select
	for {
		select {
		case <-ticker.C:
			// Este caso se ejecuta cada vez que el ticker genera un evento (cada 20 segundos).
			lectura_sysinfo, err := leerProc("/proc/sysinfo")
			if err != nil {
				fmt.Println(err)
				return
			}

			// Extraer los valores directamente del map
			totalram := int64(lectura_sysinfo["Totalram"].(float64))
			freeram := int64(lectura_sysinfo["Freeram"].(float64))
			procs := int(lectura_sysinfo["Procs"].(float64))

			// Insertar una sola fila con los valores del sistema
			_, err = db.Exec("INSERT INTO registros (total_ram, ram_libre, total_procesos, created_at) VALUES (?, ?, ?, ?)",
				totalram, freeram, procs, time.Now().UnixMilli())
			if err != nil {
				// Si ocurre un error al insertar, se registra en los logs.
				log.Println("Error insertando:", err)
			}

		case <-sigs:
			// Este caso se ejecuta cuando se recibe una señal de interrupción o terminación.
			// Detiene el daemon y sale del bucle.
			log.Println("Daemon detenido.")
			break loop
		}
	}
}

// ======================================================= FUNCIONES COMPLEMENTARIAS PARA EL MAIN ==============================
func ejecutarDockerCompose() {
	composeFilePath := "../grafana/docker-compose.yaml"

	cmd := exec.Command("docker", "compose", "-f", composeFilePath, "up", "-d")
	output, err := cmd.CombinedOutput()
    	if err != nil {
    		log.Fatalf("Error al ejecutar el docker compose: %v\nOutput: %s", err, output)
    	}

    	fmt.Printf("Docker Compose Up exitoso:\n%s\n", output)
}

func leerProc(ruta string) (map[string]interface{}, error) {
	// Se abre el archivo
	file, err := os.Open(ruta)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo %s: %v", ruta, err)
	}
	defer file.Close()

	// Se lee todo el contenido del archivo
	contenido, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo: %v", err)
	}

	// Se parsea el JSON a un map genérico
	var datos map[string]interface{}
	err = json.Unmarshal(contenido, &datos)
	if err != nil {
		return nil, fmt.Errorf("error al parsear JSON: %v", err)
	}

	return datos, nil
}