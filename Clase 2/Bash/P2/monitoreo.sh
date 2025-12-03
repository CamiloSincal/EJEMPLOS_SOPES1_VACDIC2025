#!/bin/bash

# Script para monitorear el uso de memoria y enviar una advertencia

# Obtener el uso de memoria
# Comandos a utilizar
# free -m : Muestra el uso de memoria en MB
# grep Mem : Filtra las lineas que contienen la palabra Mem
# awk '{print $4/$2 * 100.0}' : Imprime la tercera columna de un archivo

mem_libre=$(free -m | grep Mem | awk '{print $4/$2 * 100.0}')

# Comprobar si la memoria libre es inferior al 20%
# (( )) : Evalua una expresión aritmética
# |  : Pipe, redirige la salida de un comando a la entrada de otro
# bc -l : Calculadora de precisión arbitraria
if (( $(echo "$mem_libre < 20" | bc -l) )); then 
    echo "Advertencia: La memoria libre es de $mem_libre%"
else
    echo "Memoria libre normal"
fi