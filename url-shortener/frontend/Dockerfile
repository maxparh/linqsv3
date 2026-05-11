# Build stage
FROM node:24-alpine AS frontend-builder

WORKDIR /app

# Кэшируем node_modules
COPY package*.json ./
RUN npm ci

# Копируем весь исходный код (нужен для Tailwind JIT)
COPY . .

# Передаем переменную окружения для API URL
ARG VITE_API_BASE_URL=/
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

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