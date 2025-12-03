#!/bin/bash

# Script para capturar las teclas y escribir en un archivo cada 5 segundos

# Creamos el archivo donde vamos a guardar las teclas
archivo="teclas.log"

# Limpiamos el archivo
> $archivo

# Mostramos un mensaje

echo "Presiona las teclas que quieres capturar. Presiona Ctrl+C para salir"

# Creamos un bucle infinito 
while true; do

    # capturar la entrada
    read -n 1 -s tecla # -n 1 captura solo un caracter, -s no muestra la tecla, tecla es la variable donde se guarda la tecla

    # escribimos y guardamos 
    # El comando >> conocido como redirección de salida, se utiliza para redirigir la salida de un comando a un archivo
    # Existen dos tipos de redirección de salida: >> y >
    # >> se utiliza para redirigir la salida de un comando a un archivo, pero si el archivo ya existe, el comando añadirá la salida al final del archivo
    # > se utiliza para redirigir la salida de un comando a un archivo, pero si el archivo ya existe, el comando sobrescribirá el archivo
    
    echo -n "$tecla" >> $archivo # -n no añade un salto de línea al final 

done