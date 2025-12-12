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

De manera técnica no decimos que hacemos deploymentes de contenedores directamente en k8s, hacemos deployments de pods y usamos cosas sobre el pod para controlarlo.

### Controlador
Sirve para crear o actualizar pods y otros objetos de k8s, casi siempre será necesario crear un controlador, en k8s podemos encontrar los siguientes tipos de controladores:

- ### Deployments, ReplicaSet
    Son los que utilizarán casi siempre, controlan a los pods en niveles muy bajos.

### Service o Servicio
Es el final o punto final dado hacia un pod, al usar un deployment para un conjunto de pods, es en ese momento que se configura un servicio.

### Namespace
Un Namespace es un "filtro" para la línea de comandos, es para ver solo lo que nosotros queramos.

### ConfigMap y Secret
ConfigMap: Almacena configuraciones no sensibles.
Secret: Almacena datos sensibles (como contraseñas) de forma segura.

---

## Kubectl
Es la herramienta de línea de comandos de k8s, que también se conoce como kube control

Se utiliza para interactuar y administrar el clúster

### Algunos comandos:
- ```kubectl version```: Para ver la versión instalada de Kubectl
- ```kubectl run```: Crear un pod
- ```kubectl cluster-info```: Para ver información del cluster
- ```kubectl get nodes```: Para ver qué nodos están disponibles en el clúster

Ejemplo simple:
```bash
kubectl run nginx --image nginx #Crea un Pod con nombre nginx
kubectl create deployment nginx --image nginx #Creamos un deployment con nombre de ngnix
```

---

## PODS
Nuestro objetivo es desplegar nuestra aplicación en forma de contenedores, sin embargo, k8s no implementa directamente sobre los nodos, lo hace en pods (la unidad mínima)

Lo usual es que por cada pod exista un contenedor. Si decrece el número de usuarios, solo se elimina un pod, no se suelen añadir contenedores adicionales a un Pod existente para el escalamiento.

Sin embargo, k8s no limita tener un contenedor por pod, podemos tener una aplicación entera en un pod (no es lo normal) y escalar todo con la creación de pods.

### Creación de pods y relación con deployments
Empecemos creando un pod:
```bash
kubectl run nginx --image nginx
```

Para eliminar el pod:
```bash
kubectl delete pod <nombre-del-pod>
```
Esto solo crea el pod, no el deployment, si no creamos un Deployment en Kubernetes, los Pods se crean individualmente (o como un ReplicaSet no administrado), pero se pierden las funcionalidades clave de autocuración, escalado automático y actualizaciones declarativas

El siguiente comando realiza un deployment de un contenedor de docker creando un pod de k8s:
```bash
kubectl create deployment nginx --image nginx #Creamos un deployment con nombre de ngnix
```
Para entender mejor, primero se crea un pod automáticamente y hace un deployment de una instancia de la imagen del docker de nginx. Tener en cuenta que para obtener la imagen, se necesita especificar el nombre usando ```--image```

Para ver el pod creado:
```bash
kubectl get pods # Información general del pod
kubectl get pods -o wide # Agregar ip y el nodo donde está disponible en la salida
```

Para ver el deployment:
```bash
kubectl get deployments
```

Para eliminar el pod y deployment:
```bash
kubectl delete deployment <nombre-del-deployment>
```

Cuando creamos un pod a través de un deployment, la eliminación del deployment automáticamente elimina el pod.

### PODs y YAML
### Logs
Existe el comando ```kubectl logs``` que nos mostrará todos los recursos relacionados a un pod (también se puede especificar otros objetos).

Ejemplo de Pingpong:
```bash
kubectl run nginx --image pingpong ping 1.1.1.1
kubectl logs pod/pingpong
```

### Pod con Yaml
Kubernetes usa yaml como "entradas" para crear sus diferentes objetos (pods, replicasets, deployments, etc). Para esto se definen 4 campos obligatorios:
- ```apiVersion```: La versión de api de kubernetes a usar para crear objetos, esta varía según el objeto, para pods es v1.
- ```kind```:Indica el tipo de objeto a crear
- ```metadata```:Se indican cosas como names y labels
- ```spec```:Se indican especifaciones específicas a k8s, como los contenedores para el pod.

Cuando tengamos un archivo Yaml para algún objeto usaremos el siguiente comando para crearlo:
```bash
kubectl create -f <nombre-del-yml> # Se considera que el comando se ejecuta en el directorio del .yml
```

---

## Replication Controller
Recordemos que los *controladores* son el cerebro de k8s, ya que se encargan de monitorear los objetos para que actuen como se espera

El replication controller se encarga y asegura que se ejecuten el número específico de pods que se le hayan indicado en todo momento, también sirve para balancear y escalar la aplicación.

Este controlar nos ayuda a ejecutar varias instancias de un mismo pod en un mismo clúster de k8s. En caso de que solo tengamos 1 pod, este controlador crea un nuevo pod en caso de que el anterior falle.

La estructura de un yml de este controlador cuenta con los mismos niveles, con la única diferencia que spec tiene un nuevo campo:
- ```template```: Acá se indica el pod al cuál se le aplicará el controlador

Además también se agrega ```replicas```, la cual sirve para indicar cuantas copias del pod se necesitan en todo momento para nuestra aplicación

Ubicamos nuestra carpeta con el controlador y para "ejecutar" el archivo usamos:
```bash
kubectl apply -f rctest.yml
```

En caso de que queramos ver nuestros controladores de réplicas, usamos:
```bash
kubectl get replicationcontroller 
```

Para ver los pods creados por el replication controller:
```bash
kubectl get pods
```

Si borramos un pod, el replication controller creará uno nuevo, esto se puede probar con:
```bash
kubectl delete pod <nombre-del-bod>
```

Si volvemos a usar un get de los pods, entonces veremos que aún existe la cantidad de pods que indicamos, pero ya no estará el que eliminamos, será uno nuevo

Para borrar un replication controller:
```bash
kubectl delete replicationcontroller <nombre-del-controller>
```

### Replication controller y replica set
- **Replication Controller (RC)**: Utiliza únicamente selectores basados en igualdad (equality-based selectors). Esto significa que solo puede seleccionar Pods que tengan una etiqueta con un valor exacto y específico (ej: app=nginx).
- **Replica Set (RS)**: Introduce selectores basados en conjuntos (set-based selectors), que son mucho más potentes y expresivos. Permite filtrar etiquetas utilizando operadores como In, NotIn, Exists o DoesNotExist (ej: app in (nginx, apache)).

---
## Deployments
Es un objeto de kubernetes en 1 nivel superior a los ReplicaSet.
La función principal de un deployment es actualizar y hacer un upgrade a instancias subyacentes con actualizaciones continuas (rolling updates), también permite deshacer, pausar y reanudar cambios según la necesidad.
Al crearautomáticamente se agrega un ReplicaSet que a la vez crean los PODs de nuestra aplicación, algunas características de los deployments son:

### Rolling update
Con k8s vimos que podemos tener varias instancias de una aplicación, si nosotros tratamos de actualizar nuestra aplicación, aplicar los cambios a todas las instancias de golpe podría generar problemas a los usuarios. Rolling Update permite aplicar los cambios a una instancia tras otra

### Rollback 
Es simplemente deshacer cambios recientes y volver atrás.

### ¿Cómo funciona un Deployment?
¿Cómo funciona?
- **Etiquetado de Pods**: Cuando se define un Deployment, también se describe un template (plantilla) para los Pods que creará. En esa plantilla, se usan labels para asignar etiquetas (ej: app: mi-aplicacion, version: v1) a los Pods que se generen.
- **Selección por matchLabels**: La sección matchLabels en el Deployment especifica las etiquetas que debe buscar.
- **Control y Mantenimiento**: El controlador del Deployment (y su ReplicaSet) monitorea constantemente el clúster. Si encuentra Pods que coincidan con esas etiquetas matchLabels, los gestiona (los mantiene vivos, los actualiza o los elimina). Si un Pod con esas etiquetas desaparece, el Deployment creará uno nuevo para mantener el replicas deseado. 

### Comandos para crear un deployment con un yml
Con nuestro yml de deploy creado, usamos el siguiente comando para hacer deploy:
```bash
kubectl apply -f <nombre-del-deploy>
```

Para verificar todo lo que el deploy crea usamos:
```bash
kubectl get all
```

Es posible cambiar la cantidad de replicas del replica set definido en el deploy con:
```bash
kubectl deploy/<nombre-del-deploy> --replicas <num-replicas-deseadas>
```

Para eliminar el deploy usamos:
```bash
kubectl delete deployment <nombre-del-deploy>
```

### Actualización de Deployments
Hay 2 formas de actualizar la versión de la aplicación, contenedores, labels, etc.
#### Recreate
Esta estrategia consiste en dar de baja todas las instancias de la aplicación par luego crear las nuevas, esto si bien funciona generará inactividad de la aplicación

#### Rolling Update
En esta estrategia se da de baja y se crean las instancias de manera consecutiva, una instancia tras otra

#### Rollback
Es el proceso de revertir una actualización si algo sale mal durante el proceso. Para esto, el deployment eliminará los pods nuevos uno por uno y traerá los pods de la versión anterior.

Para poner a prueba estos conceptos tendremos que tener la modificar el archivo ```ej3-deploy```, específicamente indicaremos la versión de nginx en imagen, en este caso usaremos la ```1.17.3```

Con la imagen definida usamos:
```bash
kubectl apply -f <nombre-del-deploy>
```

Recibimos el siguiente comando para revisar el deploy:
```bash
kubectl apply describe <nombre-del-deploy>
```

Ahora usemos el siguiente comando para ver el estado del rollout:
```bash
kubectl rollout status deploy/<nombre-del-deploy>
```

Ahora revisar el historial del deploy:
```bash
kubectl rollout history deploy/<nombre-del-deploy>
```

No se verá un registro como tal, para que en un cambio se mire el registro o record se usa:
```bash
kubectl apply -f <nombre-del-deploy> --record
```

Para actualizar el deployment sin usar el yml, se usa:
```bash
kubectl set image deployment/<nombre-del-deploy> nginx-container=nginx:1.17.10
```

Revisemos el deployment con:
```bash
kubectl describe deployment
```

Si tuvieramos un error en nuestra nueva versión y necesitamos hacer un rollback, usamos:
```bash
kubectl rollout undo deployment/<nombre-del-deploy>
```

Para revisar las versiones:
```bash
kubectl rollout history deploy/<nombre-del-deploy>
```

## Servicios
Los servicios permiten comunicar diferentes componentes dentro y fuera de la aplicación.

Para la conexión a pods, se requiere de un servicio, en pocas palabras, un service o servicio es una dirección estable para un pod o grupo de pods.

En la práctica buscamos exponer servicios para crear estos recursos que apuntan a pods de backends. Esto normalmente se hace con ```kubectl expose``` o con archivos yaml.

Hay distintos tipos de servicios:
- **ClusterIP**: Funciona en cualquier configuración, expone el servicio en una *IP interna* del clúster. Es decir, solo se encuentra disponible dentro del clúster (entre nodos y pods)
- **NodePort**: Funciona en cualquier configuración, está diseñado para ser accesible desde fuera de nuestro clúster, disponible para todos los nodos y cualquiera que se quiera conectar a el. (Se tendrán puertos altos)
- **LoadBalancer**: Es usar un servicio externo de un tercero, puede ser un proxy o firewall (AWS,Azure,GCP). Acá se expone el servicio de manera externa con un proveedor cloud.
- **ExternalName**: Proporciona un alias inerno para un nombre de DNS externo. Los clientes solicitan el DNS y las solicitudes se redireccionan a un nombre externo.

### Ejemplos
1. #### NodePort
   
    El primer ejemplo es de NodePort, para este tendremos que tener ejecutado el pod del yaml ```ej1-pod.yml```, si no lo está, entonces ejecutamos con:
    ```bash
    kubectl apply -f <nombre-del-yml>
    ```
    Para aplicar el servicio usamos:
    ```bash
    kubectl apply -f <nombre-del-servicio>
    ```

    Para revisar nuestros servicios usamos:
    ```bash
    kubectl get services
    ```

2. #### ClusterIP

    El segundo ejemplo es de ClusterIP, para este caso los comandos son los mismos que del ejemplo anterior, y la única diferencia notable a simple vista será el tipo de Servicio Creado. Para esto tendremos que tener ejecutado el pod del yaml ```ej1-pod.yml```, si no lo está, entonces ejecutamos con:
    ```bash
    kubectl apply -f <nombre-del-yml>
    ```
    Para aplicar el servicio usamos:
    ```bash
    kubectl apply -f <nombre-del-servicio.yml>
    ```

    Para revisar nuestros servicios usamos:
    ```bash
    kubectl get services