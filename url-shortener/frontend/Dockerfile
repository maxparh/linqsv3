# Build stage
FROM node:24-alpine AS frontend-builder

WORKDIR /app

# Кэшируем node_modules
COPY package*.json ./
RUN npm ci

# Копируем весь исходный код (нужен для Tailwind JIT)
COPY . .

# 🔥 Передаём переменные окружения для Vite
ARG VITE_API_BASE_URL=/
ARG VITE_SHORT_LINK_DOMAIN=http://localhost

ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_SHORT_LINK_DOMAIN=${VITE_SHORT_LINK_DOMAIN}

# Собираем проект
RUN npm run build

# Production stage (Nginx)
FROM nginx:alpine

# Копируем конфиг Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем собранный фронтенд
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]