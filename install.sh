#!/usr/bin/env bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

info "Проверка зависимостей..."

if ! command -v docker &>/dev/null; then
    error "Docker не установлен. Установите Docker: https://docs.docker.com/get-docker/"
fi

if ! docker compose version &>/dev/null; then
    error "Docker Compose не доступен. Обновите Docker до актуальной версии."
fi

info "Генерация JWT_SECRET..."
JWT_SECRET=$(openssl rand -hex 32 2>/dev/null || python3 -c "import secrets; print(secrets.token_hex(32))" 2>/dev/null || echo "dev_secret_change_me_$(date +%s)")

ENV_FILE="url-shortener/.env"
if [ ! -f "$ENV_FILE" ]; then
    info "Создание $ENV_FILE..."
    cat > "$ENV_FILE" <<EOF
JWT_SECRET=$JWT_SECRET
FRONTEND_URL=http://localhost
API_BASE_URL=/api
VITE_API_BASE_URL=/api
VITE_SHORT_LINK_DOMAIN=http://localhost
EOF
else
    warn "$ENV_FILE уже существует, пропускаем"
fi

BACKEND_ENV="url-shortener/backend/.env"
if [ ! -f "$BACKEND_ENV" ]; then
    info "Создание $BACKEND_ENV..."
    cat > "$BACKEND_ENV" <<EOF
DATABASE_URL=postgres://appuser:apppassword@postgres:5432/urlshortener?sslmode=disable
JWT_SECRET=$JWT_SECRET
JWT_EXPIRATION=3600
SERVER_PORT=:8080
FRONTEND_URL=http://localhost
EOF
else
    warn "$BACKEND_ENV уже существует, пропускаем"
fi

FRONTEND_ENV="url-shortener/frontend/.env"
if [ ! -f "$FRONTEND_ENV" ]; then
    info "Создание $FRONTEND_ENV..."
    cat > "$FRONTEND_ENV" <<EOF
VITE_API_BASE_URL=http://localhost:8080
EOF
else
    warn "$FRONTEND_ENV уже существует, пропускаем"
fi

info "Запуск Docker Compose..."
cd url-shortener
docker compose up -d --build

echo ""
info "========================================"
info " Linqs запущен!"
info " Frontend: http://localhost"
info " Backend:  http://localhost:8080"
info "========================================"
echo ""
warn "Nginx сконфигурирован для localhost."
warn "Для продакшена замените server_name в url-shortener/frontend/nginx.conf"
echo ""
