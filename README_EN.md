[Русский](README.md) | [English](README_EN.md)

---

# Linqs

**Linqs** is an open-source URL shortener.

## Tech Stack

- **Backend:** Go (pgx, golang-jwt, godotenv)
- **Frontend:** Vue 3 + TypeScript + Vite + Pinia + Vue Router + ECharts
- **Database:** PostgreSQL + Redis
- **Infrastructure:** Docker Compose (Nginx, multi-stage builds)

## Features

### Users & Authentication
- Registration via email or phone
- Login/logout with JWT tokens
- Protected routes

### URL Shortening
- Create short links with auto-generated codes
- Support for custom short codes
- Redirect with click tracking
- Link management (list, delete)

### Analytics
- **Overview:** total clicks, unique clicks, bounce rate, average time on site
- **Clicks over time:** daily traffic chart
- **Geography:** distribution by country
- **Devices:** desktop vs mobile, browser breakdown (Chrome, Firefox, Safari, Edge, Opera)
- **Bounce rate:** per-link bounce percentage
- Session tracking with IP hashing (SHA-256)

### Additional
- Admin user panel
- Profile settings
- Cookie Consent Banner
- Responsive design (Tailwind CSS)
- Fully containerized setup

## Deployment

### Important

The nginx configuration is tailored for the `linqs.ru` domain. When deploying on your own server, **first of all** update `server_name` and related settings in `url-shortener/frontend/nginx.conf` to match your domain.

### Quick Start (local)

```bash
chmod +x install.sh
./install.sh
```

The script automatically checks dependencies, generates a JWT secret, configures the environment, and starts the project via Docker Compose.

> **Note:** `install.sh` has not been tested, manual setup via `docker compose up -d` is recommended.

### Self-Hosted (manual)

```bash
cd url-shortener
docker-compose up -d
```

> **Note:** self-hosted deployment includes **basic** URL shortening and link management functionality only.

### Extended Features

The full feature set, including advanced analytics, priority support, and additional capabilities, is available exclusively at [linqs.ru](https://linqs.ru).

## License

This project is licensed under the **MIT License**. You are free to modify, distribute, and use the code for both personal and commercial purposes, provided that the original copyright notice is included.

See [LICENSE](LICENSE) for details.
