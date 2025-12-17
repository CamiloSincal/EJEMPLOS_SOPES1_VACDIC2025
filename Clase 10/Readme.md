# Ingress para Api Rust
Para este ejemplo se considerará que ya se tiene creada la VM en GCP y que se tiene instalado Zot + la imagen de la API de Rust cargada en el registro de imágenes privado.

## Deployment y service de API-RUST
Para facilitar la descarga de imagenes de Zot en la VM, se recomienda usar Ngrok. Por este motivo, lo primero es que en la vm se ejecuten estos comandos:
```yml
# Para instalar ngrok:
curl -sSL https://ngrok-agent.s3.amazonaws.com/ngrok.asc \
  | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null \
  && echo "deb https://ngrok-agent.s3.amazonaws.com bookworm main" \
  | sudo tee /etc/apt/sources.list.d/ngrok.list \
  && sudo apt update \
  && sudo apt install ngrok

# Para agregar el token:
ngrok config add-authtoken 2w3ewknemOBcjLqYLPfD1kp0MXs_5QCzSL96KdfS1cMLbUhc6

# En la VM con Zot:
ngrok http 5000 --host-header=rewrite
```

Ngrok dará una URL como: https://xxxx-xx-xx-xx-xx.ngrok-free.app

Para poder realizar el deployment de la api de rust utilizaremos esta estructura, tomar en cuenta que en image no se coloca el https o http de ngrok, únicamente lo que está después del enlace:
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-rust-deploy
  labels:
    app: api-rust
spec:
  replicas: 1
  selector:
    matchLabels:
      app: api-rust
  template:
    metadata:
      labels:
        app: api-rust
    spec:
      containers:
      - name: api-rust
        # Usa la URL de ngrok que se obtuvo
        image: xxxx-xx-xx-xx-xx.ngrok-free.app/api_rust:v1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
```

Luego usaremos lo siguiente para levantar el servicio del deployment:
```yml
apiVersion: v1
kind: Service
metadata:
  name: api-rust-service
  labels:
    app: api-rust
spec:
  selector:
    app: api-rust
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

**Nota: Es posible ejecutar primero el archivo del deployment y luego el del service o implementar un multistage como en el archivo adjunto para solo ejecutar un archivo y no dos de manera separada**

Para este caso en específico, podemos ver los logs en tiempo real de nuestra api de rust con:
```bash
kubectl logs -f <nombre-pod-api-rust>
```

## Ingress
Con el deployment de k8s realizado, lo primero es instalar el controller de Nginx:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
```

Luego levantamos nuestro ingress con:
```bash
kubectl apply -f <nombre-archivo-ingress>
```

Luego obtenemos la IP del ingress con:

```bash
kubectl get svc -n ingress-nginx ingress-nginx-controller
# Ver EXTERNAL-IP
# Acceder: http://<EXTERNAL-IP>/
```
# Logs en Tiempo Real

## 1. Ver logs de tu API Rust (más importante)
```powershell
kubectl logs -f deployment/api-rust-deploy
```
**Se debería ver:** Cada petición que llega con los datos de clima

## 2. Ver logs del Ingress Controller
```powershell
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller -f --tail=100
```
**Se debería ver:** Cada request HTTP que pasa por el Ingress

## 3. Monitorear ambos a la vez
```powershell
# Terminal 1
kubectl logs -f deployment/api-rust-deploy

# Terminal 2
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller -f --tail=100
```

## 4. Ver estado de pods en tiempo real
```powershell
kubectl get pods --watch
```

## Interpretación Rápida

###  TODO FUNCIONA
- **Logs Rust:** Se ven "Datos de clima recibidos: Lugar: guatemala..."
- **Logs Ingress:** Se ven códigos 200

### Ingress no recibe nada
- **Logs Rust:** Vacío (no llegan peticiones)
- **Logs Ingress:** Vacío o errores 404/502
- **Problema:** IP incorrecta en Locust o Ingress mal configurado

### Ingress recibe pero Rust no
- **Logs Ingress:** Ves peticiones pero con códigos 502/503
- **Logs Rust:** Vacío
- **Problema:** Service no conecta con los pods (puerto incorrecto)

### Pods reiniciando
- **Logs Pods:** `CrashLoopBackOff`
- **Problema:** Aplicación Rust falla al iniciar