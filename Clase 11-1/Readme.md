# Kubernetes p2

#### Modelo de Network en Kubernetes
En kubernetes todo puede ser accesible a todo, si se quiere seguridad, se deben adicionar políticas de red, esto se hace mediante implementación de red de terceros (Flannel, Romana, Cilium, etc)

## Namespaces
Es un filtro en la línea de comando de kubectl
Al ejecutar   ```kubectl get namespaces``` podemos ver solo las cosas que nos interesan.

Estos permiten separar separar por etiquetas distintos objetos y llevar un mejor control sobre pods, deployments, etc. o incluso para llevar un control sobre lo que debe trabar cada persona en un equipo grande.

Para crear un namespace:
 ```bash
kubectl create namespace <nombre-namespace-crear>
# O usando la abreviatura
kubectl create ns <nombre-namespace-crear>
```

Es posible crear los objetos a través de la línea de comandos
 ```bash
kubectl run mi-pod-nginx --image=nginx --restart=Never --namespace=<nombre-namespace-creado>
# O la versión corta:
kubectl run mi-pod-nginx --image=nginx --restart=Never -n <nombre-namespace-creado>
```

También es posible indicar el namespace en un archivo yml:
```yml
apiVersion: v1
kind: Namespace
metadata:
  name: mi-nuevo-namespace
  labels:
    app: ejemplo-ns
```

Para ver todos los recursos de todos los namespaces usamos:
 ```bash
kubectl get pods --all-namespaces
```

Para filtrar un namespace en específico podemos usar:
 ```bash
kubectl get <objeto> --namespace=<namespace_a_buscar>

kubectl get <objeto> -n <namespace_a_buscar>
```

## ConfigMaps
Es un recurso en k8s que se encarga de almacenar contenido de archivos o key/values en el API de K8s dentro de etcd.

Estos pueden servir para:
- Contener uno o más archivos de configuración
- Contener parámetros de configuración individuales

Generalmente sirve para almacenar información NO CONFIDENCIAL en el formato de llave valor. Existen 3 formas de aplicar una configuración específica a un contenedor:
- Archivos de configuración (Configmaps)
- Arugmentos por terminal
- Variables de entorno

Este archivo debe crearse ANTES de la creación del pod, ya que se le hace referencia en la sección de spec.

Para crearlo usaremos:
 ```bash
 # map-name = nombre del configmap
 # data-source = origen de datos (directorio, archivo, llave-valor)
kubectl create configmap <map-name> <data-source>
```
Si hace referencia a un directorio o archivo, se usa ```--from-file```, si se usa un llave-valor, se usa ```--from-literal```

### Creación de un Config a partir de un directorio
Para este ejemplo nos ubicaremos en la carpeta con nuestros archivos.
```bash
kubectl create configmap <nombre-configmap> --from-file=config
```

Podemos revisar el archivo configurado en k8s con:
 ```bash
kubectl get configmap -o wide
```

Podemos revisar el contenido usamos:
 ```bash
kubectl get configmap <nombre-configmap> -o yaml
```

Para eliminar el configmap se usa:
 ```bash
kubectl delete configmap <nombre-configmap> 
```

### Creación de un Config a partir de un solo archivo
De maner similar al ejemplo anterior, nos colocaremos en el directorio del archivo de redis-config para luego ejecutar:
 ```bash
#Para este ejemplo se usará como nombre "example-redis-config"
kubectl create configmap  <nombre-configmap> --from-file=redis-config
```
Podemos revisar el archivo configurado en k8s con:
 ```bash
kubectl get configmaps
```

Podemos revisar las características del config usamos:
 ```bash
kubectl describe configmap <nombre-configmap>
```


Para cargar o introducir el configmap en un pod en este ejemplo usaremos volúmenes, sin embargo, es necesario aclarar que es posible hacerlo de esta manera o con variables de entorno.

Para hacerlo con volúmenes usarmos un configmap.yml:
```yml
apiVersion: v1
kind: Pod
metadata:
  name: redis
spec:
  containers:
  - name: redis
    image: redis:5.0.4
  # Especificamos donde se montará el volumen
    volumeMounts:
    - mountPath: /redis-master
      name: config #nombre del volúmen

  volumes:
    - name: config
      configMap:
        name: example-redis-config
        items:
        - key: redis-config # nombre para el volúmen
          path: redis.conf # ubicación donde la información se presentará en el pod
```

Con el archivo creado ejecutaremos:
```bash
kubectl create -f <nombre-del-yml-del-configmap>
```

Para verificar que el volúmen está correctamente configurado usamos:
```bash
# El directorio del volúmen del archivo anterior es "/redis-master/redis.conf"
kubectl exec redis -- cat <directorio-del-volúmen>
```

Para eliminar el configmap y e pod:
```bash
kubectl delete configmap <nombre-configmap>
kubectl delete pod redis
```

### Creación de un Config a partir de valores literales
Este tipo de configmaps se crean completamente en la línea de comandos, para esto usaremos:
```bash
kubectl create configmap <nombre-configmap> --from-literal=special.how=very
# La clave será special.how
# Y el valor será very
```

Con el configmap configurado, ahora los insertaremos en un pod. Para este ejemplo se usarán variables de entorno, y el archivo yml para hacer esto es el siguiente:

```yml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: test-container
    image: k8s.gcr.io/busybox
    command: ["/bin/sh","-c","env"] # Se muestran las variables de entorno al ejecutar el pod
    env: # Acá se definen las variables de entorno
      - name: SPECIAL_LEVEL_KEY
        valueFrom:
          configMapKeyRef:
            name: special-config # nombre del config
            key: special.how # nombre de la clave a introducir
  restartPolicy: Never
  ```

Con el archivo configurado procedemos a aplicar el archivo del pod:
```bash
kubectl apply -f <nombre-del-archivo-del-pod-con-config-por-env>
```

Para verificar su funcionamiento revisaremos los logs con:
```bash
kubectl logs <nombre-pod> | grep SPECIAL
```

Eliminamos todo con:
```bash
kubectl delete pod <nombre-pod>
kubectl delete configmap <nombre-config>
```

## Secrets
Cuando necesitamos manejar información sensible como usuarios, contraseñas, tokens, etc., es necesario utilizar secrets para contener esta información.

Los secrets buscan evitar exponer este tipo de información al hacer deployments en los pods.

Estas son almacenadas en el etcd en k8s y podrán ser utilizados por diferentes objetos.

Al igual que los configmaps se rigen por claves valor, pero con la diferencia de que los Secrets tienen un tamaño máximo de 1MB, y de manera similar es posible insertarlos a través de volúmenes o variables de entorno.

### Secrets a través de kubectl
Para crear un secret a partir de archivos y con kubectl usamos el siguiente comando:
```bash
# para este ejemplo se usará el tipo genérico 
# Para el tipo de información usaremos generic mi-secret --from-file=./secrets/username.txt --from-file=./secrets/password.txt
kubectl create secret <tipo-secret> <nombre-secret> <tipo-de-información>
```

Para ver que se creó adecuadamente usamos:
```bash
kubectl get secrets
```

Para ver información detallada:
```bash
kubectl describe secrets <nombre-secret>
```

Para eliminar el secret:
```bash
kubectl delete secret <nombre-secret>
```

### Secrets a través de archivos de manifiesto
Los archivos de manifesto también son conocidos como "Manifest Files" y es la forma manual de crear secrets, se ceran los secretos a partir de un archivo tipo YAML y luego se creará el objeto.

Para este procedimiento se debe encriptar la información usando base64, teniendo los archivos con base64 crearemos nuestro yml con:
```yml
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MTNmMjM2OGFzZg==
  ```
Con el archivo creado ejecutamos:
```bash
kubectl apply -f <nombre-del-archivo-secret>
```

Para verificar la creación usamos:
```bash
kubectl get secrets
```

Para ver información detallada:
```bash
kubectl describe secrets <nombre-del-secret>
```

Para eliminar el secret:
```bash
kubectl delete secret <nombre-secret>
```

### Introducir Secrets en un pod
Para introducir secrets dentro un pod primero tendremos que tener creados dichos valores. Para este ejmplo usaremos el secret del archivo mysecret.yml.

Existen 2 formas (al igual que en configmaps) de introducir estos valores, por volúmenes o por variables de entorno.

#### Volúmenes
Para este primer método usaremos el siguiente archivo yml:
```yml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: mypod
    image: redis
  # Especificamos donde se montará el volumen
    volumeMounts:
    - name: foo
      mountPath: "/etc/foo"
      readOnly: true

  volumes:
    - name: foo
      secret:
        secretName: mysecret
```

Con el archivo creado ejecutaremos el siguinte comando:
```bash
kubectl apply -f <nombre-del-yml>
```

Verificamos el pod:
```bash
kubectl get pods
```

Para verificar el contenido del volúmen:
```bash
kubectl exec <nombre-pod> -- ls /etc/foo

kubectl exec <nombre-pod> -- cat /etc/foo
```

Eliminamos el pod con:
```bash
kubectl delete <nombre-pod>
```

#### Variables de entorno
Para introducir nuestras variables de entorno dentro de un pod con variables de entorno usaremos un archivo yml con esta estructura:
```yml
apiVersion: v1
kind: Pod
metadata:
  name: secret-env-pod
spec:
  containers:
  - name: mycontainer
    image: redis
    env:
      - name: SECRET_USERNAME
        valueFrom:
          secretKeyRef:
            name: mysecret
            key: username
      - name: SECRET_PASSWORD
        valueFrom:
          secretKeyRef:
            name: mysecret
            key: password
restartPolicy: Never
```

Con el archivo yml creado ejecutamos:
```bash
kubectl create -f <nombre-del-yml>
```

Verificamos el pod:
```bash
kubectl get pods
```

Para verificar las variables de entorno:
```bash
kubectl exec <nombre-pod> env | grep SECRET
```

Eliminamos el pod con:
```bash
kubectl delete <nombre-pod>
```

Eliminamos el secret:
```bash
kubectl delete secret <nombre-secret>
```