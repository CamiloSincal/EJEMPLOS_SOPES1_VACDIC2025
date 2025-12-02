FROM python:3.9-slim

# Establecer directorio de trabajo
WORKDIR /app

# Copiar requirements primero (para aprovechar cache de Docker)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copiar el resto de la aplicación
COPY . .

# Variables de entorno
ENV PYTHONUNBUFFERED=1
ENV PORT=8000

# Exponer puerto
EXPOSE 8000

# Comando para ejecutar
CMD ["python", "app.py"]