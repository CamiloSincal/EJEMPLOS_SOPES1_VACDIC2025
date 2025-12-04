# Creación de módulos en Linux
## Instalación/configuración de C y python en Linux
Para instalar C en Linux, se debe instalar el compilador de C, gcc. Para instalar gcc se debe ejecutar el siguiente comando:
```bash
sudo apt install gcc

# ver versión
gcc --version
```

---
Además, dado que se usará *Makefile* para la compilación, si no tienes instalado el paquete make puedes instalarlo con el siguiente comando:
```bash
sudo apt-get install make
```

También será necesario instalar los essentials de desarrollo en Ubuntu, los essentials de desarrollo incluyen herramientas y bibliotecas necesarias para compilar programas en C, para ello se debe ejecutar el siguiente comando:
```bash
sudo apt-get install build-essential
```
## Instalación de Python
Python usualmente viene instalado en la mayoría de las distribuciones de Linux, sin embargo, se puede verificar si está instalado ejecutando el siguiente comando:
```
python --version
```
En caso de que no esté instalado, se puede instalar Python en Ubuntu o Fedora ejecutando el siguiente comando:
```bash
sudo apt-get/dnf install python3
```
## Ejemplos
### 1. Módulo de hola mundo
```c
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
```Makefile
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
Ejecutar:
```bash
make # se compila y genera los archivos para instalar
sudo insmod <file>.ko # instalar el modulo en kernel

sudo dmesg | tail -n 20 # para ver los logs del kernel

sudo rmmod <name> # desinstalar modulo
```

### 2. Módulo para imprimir las métricas del SO en un archivo en /proc
Antes de crear y compilar este tipo de módulos es necesario preparar el entorno, para esto debemos:
- Asegurarnos de tener instalado un compilador de kernel, como gcc.
- Tener acceso a los encabezados del kernel (kernel headers). En Ubuntu y sistemas basados en el mismo, puedes instalarlos con:
```bash
sudo apt install linux-headers-$(uname -r)
```

Si durante la compilación con ```make``` se genera un error, es posible que que se tengan que generar los archivos de configuración, para eso usamos los siguiente comandos:
```bash
# Se navega al directorio de los headers
cd /usr/src/linux-headers-$(uname -r)

# Se copia la configuración actual del kernel
sudo cp /boot/config-$(uname -r) .config

# Senera los archivos faltantes
sudo make oldconfig
sudo make prepare
sudo make modules_prepare
```

En último caso es posible que se necesite una reinstalación:
```bash
# Se remueven headers actuales
sudo apt remove linux-headers-6.14.0-36-generic

# Se limpia la configuración
sudo rm -rf /usr/src/linux-headers-6.14.0-36-generic

# Se actualiza e instalar de nuevo
sudo apt update
sudo apt install --reinstall linux-headers-6.14.0-36-generic
```

#### Archivo C:
```c
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
```Makefile
obj-m += sysinfo.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```

Ejecutar:
```bash
make # se compila y genera los archivos para instalar
sudo insmod <file>.ko # instalar el modulo en kernel

sudo dmesg | tail -n 20 # para ver los logs del kernel


cat /proc/sysinfo # imprime lo escrito en el archivo 

sudo rmmod <name> # desinstalar modulo
```

### 3. Modulo para listar procesos y analizar recursos(memoria)
#### Archivo C:
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/string.h> 
#include <linux/init.h>
#include <linux/proc_fs.h> 
#include <linux/seq_file.h> 
#include <linux/mm.h> 
#include <linux/sched.h> 
#include <linux/timer.h> 
#include <linux/jiffies.h> 
#include <linux/uaccess.h>
#include <linux/tty.h>
#include <linux/sched/signal.h>
#include <linux/fs.h>        
#include <linux/slab.h>      
#include <linux/sched/mm.h>
#include <linux/binfmts.h>
#include <linux/timekeeping.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tu Nombre");
MODULE_DESCRIPTION("Modulo para leer informacion de memoria y CPU en JSON");
MODULE_VERSION("1.0");

#define PROC_NAME "sysinfo"
#define MAX_CMDLINE_LENGTH 256

// Función para obtener la línea de comandos de un proceso
static char *get_process_cmdline(struct task_struct *task)
{
    struct mm_struct *mm;
    char *cmdline;
    unsigned long arg_start = 0, arg_end = 0;
    int len = 0, i;

    cmdline = kmalloc(MAX_CMDLINE_LENGTH, GFP_KERNEL);
    if (!cmdline)
        return NULL;

    mm = get_task_mm(task);
    if (!mm) {
        kfree(cmdline);
        return NULL;
    }

    /* Kernel >= 6.8 usa mmap_lock */
    down_read(&mm->mmap_lock);
    arg_start = mm->arg_start;
    arg_end = mm->arg_end;
    up_read(&mm->mmap_lock);

    if (arg_end > arg_start)
        len = arg_end - arg_start;
    else
        len = 0;

    if (len > MAX_CMDLINE_LENGTH - 1)
        len = MAX_CMDLINE_LENGTH - 1;

    if (len > 0) {
        if (access_process_vm(task, arg_start, cmdline, len, 0) != len) {
            mmput(mm);
            kfree(cmdline);
            return NULL;
        }
    } else {
        cmdline[0] = '\0';
    }

    cmdline[len] = '\0';

    for (i = 0; i < len; i++)
        if (cmdline[i] == '\0')
            cmdline[i] = ' ';

    while (len > 0 && cmdline[len - 1] == ' ')
        cmdline[--len] = '\0';

    mmput(mm);
    return cmdline;
}


// Función para mostrar la información en formato JSON
static int sysinfo_show(struct seq_file *m, void *v) {
    struct sysinfo si;
    struct task_struct *task;
    unsigned long total_jiffies;
    int first_process = 1;
    int process_count = 0;

    // Obtenemos la información de memoria
    si_meminfo(&si);
    total_jiffies = jiffies;

    seq_printf(m, "{\n");
    seq_printf(m, "  \"Totalram\": %lu,\n", si.totalram << (PAGE_SHIFT - 10));
    seq_printf(m, "  \"Freeram\": %lu,\n", si.freeram << (PAGE_SHIFT - 10));
    
    // Contar todos los procesos
    for_each_process(task) {
        process_count++;
    }
    
    seq_printf(m, "  \"Procs\": %d,\n", process_count);
    seq_printf(m, "  \"Processes\": [\n");

    // Iterar sobre todos los procesos del sistema
    for_each_process(task) {
        unsigned long vsz = 0;
        unsigned long rss = 0;
        unsigned long totalram = si.totalram << (PAGE_SHIFT - 10);
        unsigned long mem_usage = 0;
        unsigned long cpu_usage = 0;
        char *cmdline = NULL;

        // Obtenemos los valores de VSZ y RSS
        if (task->mm) {
            vsz = task->mm->total_vm << (PAGE_SHIFT - 10);
            rss = get_mm_rss(task->mm) << (PAGE_SHIFT - 10);
            
            // Calcular porcentaje de memoria (multiplicamos por 1000 para tener 1 decimal)
            if (totalram > 0)
                mem_usage = (rss * 1000) / totalram;
        }

        // Calcular uso de CPU
        unsigned long total_time = task->utime + task->stime;
        if (total_jiffies > 0) {
            cpu_usage = (total_time * 10000) / total_jiffies;
            // Ajustar por número de CPUs
            cpu_usage = cpu_usage / num_online_cpus();
        }

        // Obtener línea de comandos
        cmdline = get_process_cmdline(task);

        // Imprimir separador entre procesos
        if (!first_process) {
            seq_printf(m, ",\n");
        } else {
            first_process = 0;
        }

        // Imprimir información del proceso
        seq_printf(m, "    {\n");
        seq_printf(m, "      \"PID\": %d,\n", task->pid);
        seq_printf(m, "      \"Name\": \"%s\",\n", task->comm);
        seq_printf(m, "      \"Cmdline\": \"%s\",\n", cmdline ? cmdline : "N/A");
        seq_printf(m, "      \"vsz\": %lu,\n", vsz);
        seq_printf(m, "      \"rss\": %lu,\n", rss);
        seq_printf(m, "      \"Memory_Usage\": %lu.%lu,\n", mem_usage / 10, mem_usage % 10);
        seq_printf(m, "      \"CPU_Usage\": %lu.%02lu\n", cpu_usage / 100, cpu_usage % 100);
        seq_printf(m, "    }");

        // Liberar memoria de cmdline
        if (cmdline) {
            kfree(cmdline);
        }
    }

    seq_printf(m, "\n  ]\n}\n");
    return 0;
}

// Función que se ejecuta al abrir el archivo /proc
static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

// Estructura que contiene las operaciones del archivo /proc
static const struct proc_ops sysinfo_ops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

// Función de inicialización del módulo
static int __init sysinfo_init(void) {
    proc_create(PROC_NAME, 0444, NULL, &sysinfo_ops);
    printk(KERN_INFO "sysinfo_json modulo cargado\n");
    return 0;
}

// Función de limpieza del módulo
static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "sysinfo_json modulo desinstalado\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);
```

#### Makefile:
```Makefile
obj-m += sysinfo.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```

Ejecutar:
```bash
make # se compila y genera los archivos para instalar
sudo insmod <file>.ko # instalar el modulo en kernel

sudo dmesg | tail -n 20 # para ver los logs del kernel


cat /proc/sysinfo # imprime lo escrito en el archivo 

sudo rmmod <name> # desinstalar modulo
```

