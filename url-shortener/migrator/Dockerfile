FROM python:3.13-slim-bookworm

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    MIGRATIONS_DIR=/app/migrations

WORKDIR /app

COPY migrator/requirements.txt ./requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

COPY migrator/migrate.py ./migrate.py
COPY backend/migrations/ ./migrations/

USER 10001:10001

CMD ["python", "migrate.py"]
