# Creación de módulos en Linux
## Instalación/configuración de C y python en Linux
Para instalar C en Linux, se debe instalar el compilador de C, gcc. Para instalar gcc se debe ejecutar el siguiente comando:
```
sudo apt install gcc

# ver versión
gcc --version
```

---
Además, dado que se usará *Makefile* para la compilación, si no tienes instalado el paquete make puedes instalarlo con el siguiente comando:
```
sudo apt-get install make
```

También será necesario instalar los essentials de desarrollo en Ubuntu, los essentials de desarrollo incluyen herramientas y bibliotecas necesarias para compilar programas en C, para ello se debe ejecutar el siguiente comando:
```
sudo apt-get install build-essential
```
## Instalación de Python
Python usualmente viene instalado en la mayoría de las distribuciones de Linux, sin embargo, se puede verificar si está instalado ejecutando el siguiente comando:
```
python --version
```
En caso de que no esté instalado, se puede instalar Python en Ubuntu o Fedora ejecutando el siguiente comando:
```
sudo apt-get/dnf install python3
```
## Ejemplos
### 1. Módulo de hola mundo
```
#include <linux/init.h> // Este archivo contiene las macros __init y __exit
/*  
    Que son los macro?
    Los macros son una forma de definir constantes en C.
    En este caso, __init y __exit son macros que se utilizan para indicarle al kernel que funciones 
    se deben llamar al cargar y descargar el modulo.

*/
#include <linux/module.h> // Este archivo contiene las funciones necesarias para la creacion de un modulo
#include <linux/kernel.h> // Este archivo contiene las funciones necesarias para la impresion de mensajes en el kernel

/*  
    El modulo debe tener una licencia, una descripcion y un autor.
*/
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("A simple Hello, World Module");
MODULE_AUTHOR("SGG");

static int __init hello_init(void) {
    printk(KERN_INFO "Hello, World!\n");
    return 0;

}

static void  __exit hello_exit(void) {
    printk(KERN_INFO "Goodbye, World!\n");
}

/* 
    Se debe indicarle al kernel que funciones se deben llamar al cargar y descargar el modulo.
*/
module_init(hello_init);
module_exit(hello_exit);
```

#### Makefile:
```
obj-m += basic.o # obj-m es una variable que contiene el nombre del modulo a compilar

all: # all es una regla que se ejecuta por defecto si no se especifica una regla
# Se compila el modulo
# Paso a paso:
# 1. Se ejecuta el comando make en el directorio /lib/modules/$(shell uname -r)/build
# 2. Se ejecuta el comando make en el directorio actual (PWD) con la regla modules
# 3. Se compila el modulo basic.c

	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules 

clean: # clean es una regla que se ejecuta para limpiar los archivos generados por la compilacion
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```

### 2. Módulo para imprimir las métricas del SO en un archivo en /proc
#### Archivo C:
```
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h> // trae las funciones para crear archivos en /proc
#include <linux/seq_file.h> // trae las funciones para escribir en archivos en /proc
#include <linux/mm.h> // trae las funciones para manejar la memoria
#include <linux/sched.h> // trae las funciones para manejar los procesos
#include <linux/timer.h> // trae las funciones para manejar los timers
#include <linux/jiffies.h> // trae las funciones para manejar los jiffies, que son los ticks del sistema


MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tu Nombre");
MODULE_DESCRIPTION("Modulo para leer informacion de memoria y CPU");
MODULE_VERSION("1.0");

#define PROC_NAME "sysinfo" // nombre del archivo en /proc

/* 
    Esta función se encarga de obtener la información de la memoria
    - si_meminfo: recibe un puntero a una estructura sysinfo, la cual se llena con la información de la memoria
*/
static int sysinfo_show(struct seq_file *m, void *v) {
    struct sysinfo si; // estructura que contiene la informacion de la memoria

    si_meminfo(&si); // obtiene la informacion de la memoria

    /*  
        El seq_printf se encarga de escribir en el archivo en /proc
        - m: es el archivo en /pro
    */

    seq_printf(m, "Total RAM: %lu KB\n", si.totalram * 4);
    seq_printf(m, "Free RAM: %lu KB\n", si.freeram * 4);
    seq_printf(m, "Shared RAM: %lu KB\n", si.sharedram * 4);
    seq_printf(m, "Buffer RAM: %lu KB\n", si.bufferram * 4);
    seq_printf(m, "Total Swap: %lu KB\n", si.totalswap * 4);
    seq_printf(m, "Free Swap: %lu KB\n", si.freeswap * 4);

    seq_printf(m, "Number of processes: %d\n", num_online_cpus());

    return 0;
};

/* 
    Esta función se ejecuta cuando se abre el archivo en /proc
    - single_open: se encarga de abrir el archivo y ejecutar la función sysinfo_show
*/
static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

/* 
    Esta estructura contiene las operaciones a realizar cuando se accede al archivo en /proc
    - proc_open: se ejecuta cuando se abre el archivo
    - proc_read: se ejecuta cuando se lee el archivo
*/

static const struct proc_ops sysinfo_ops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
};


/* 
    Esta macro se encarga de hacer dos cosas:
    1. Ejecutar la función proc_create, la cual recibe el nombre del archivo a guardar en /proc, permisos,
        y la estructura con las operaciones a realizar

    2. Imprimir un mensaje en el log del kernel
*/
static int __init sysinfo_init(void) {
    proc_create(PROC_NAME, 0, NULL, &sysinfo_ops);
    printk(KERN_INFO "sysinfo module loaded\n");
    return 0;
}

/* 
    Esta macro se encarga de hacer dos cosas:
    1. Ejecutar la función remove_proc_entry, la cual recibe el nombre del archivo a eliminar de /proc
*/
static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "sysinfo module unloaded\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);
```

#### Makefile:
```
obj-m += sysinfo.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```

Para ambos ejemplos se deben compliar con:
```
make # se compila y genera los archivos para instalar
sudo insmod <file>.ko # instalar el modulo en kernel

sudo dmesg | tail -n 20 # para ver los logs del kernel


cat /proc/sysinfo # imprime lo escrito en el archivo 

sudo rmmod <name> # desinstalar modulo
```