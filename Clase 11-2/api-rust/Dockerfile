# ============================================
# ETAPA 1: BUILD
# ============================================
FROM rust:1.83-slim AS builder

WORKDIR /app

# Instalar dependencias necesarias para compilar reqwest y openssl
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copiamos solo los archivos de dependencias primero
COPY Cargo.toml Cargo.lock ./

# Creamos un src dummy para compilar solo las dependencias
RUN mkdir -p src && \
    echo "fn main() {}" > src/main.rs && \
    cargo build --release 2>&1 || echo "Expected to fail" && \
    rm -rf src

# Ahora copiamos el código real
COPY . .

# Forzamos recompilación del código (no las dependencias)
RUN touch src/main.rs && \
    cargo build --release

# ============================================
# ETAPA 2: RUNTIME (Imagen final)
# ============================================
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    libssl3 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/target/release/api-rust /app/api-rust

RUN useradd -m -u 1001 appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 3000

CMD ["./api-rust"]