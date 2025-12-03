#!/bin/bash
# Bucle while
contador=1
while [ $contador -le 5 ]; do # Mientras contador sea menor o igual a 5
  echo "Iteración $contador"
  contador=$((contador + 1)) # Incrementar el contador0
done