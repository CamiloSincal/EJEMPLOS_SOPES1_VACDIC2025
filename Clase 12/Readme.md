# Kafka
Es una plataforma de streaming de eventos distribuida y de código abierto que permite publicar, almacenar y procesar flujos de datos en tiempo real de manera escalable y confiable, actuando como un "bus" o "cola" de mensajes de alto rendimiento para aplicaciones impulsadas por eventos, como análisis de datos, monitoreo de actividad o sistemas de comercio electrónico
## Strimzi kafka
Strimzi es un proyecto de código abierto que facilita la ejecución y gestión de clústeres de Apache Kafka en Kubernetes, proporcionando operadores nativos que simplifican enormemente su despliegue, configuración y mantenimiento, permitiendo gestionar Kafka de forma "nativa" a Kubernetes usando recursos personalizados como ```KafkaTopic``` y ```KafkaUser```
# Ejemplos
Para esta clase se probará kafka de 2 maneras. De manera local y un despliegue en k8s.

## Local con Docker
Para el ejemplo local utiilizaremos el docker-compose con:
```bash
docker compose up
```

Con el docker compose configurado tendremos que instalar y actualizar los paquetes de go:
```bash
go mod init
# luego
go mod tidy
```

Ejecutamos el ```main.go``` y el ```consumer-kafka.go``` con:
```bash
go run main
go run consumer-kafka
```

## Desplegado en k8s
Primero construimos nuestra imagen de main-k8s y el consumer con:
```bash
docker build -t <Link-Zot>:5000/main-k8s -f Dockerfile
docker build -t <Link-Zot>:5000/main-consumer-k8s -f Dockerfile.consumer
```

Luego hacemos el push de las imagenes a nuestro repositorio de imagenes privado con:
```bash
docker push <Link-Zot>:5000/main-k8s
docker push <Link-Zot>:5000/main-consumer-k8sconsumer
```

Finalmente desplegamos nuestra aplicación con el archivo ```clima-app.yml``` con:
```bash
kubectl apply -f clima-app.yml
```

Para revisar los logs y que nuestra app esté funcionado correctamente:
```bash
kubectl logs -f <pod-que-publica-en-kafka>
kubectl logs -f <pod-que-consume-kafka>
```