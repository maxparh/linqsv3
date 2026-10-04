# PostgreSQL migrations

The one-shot `migrations` Compose service runs Python with Psycopg 3. PostgreSQL
must be healthy before it starts; backend startup requires its successful exit.
SQL files from `backend/migrations` are included in the image at build time.
The PostgreSQL init-script mount is no longer used.

From `url-shortener`, start the project with:

```sh
docker compose up -d --build
docker compose logs migrations
```

The runner validates filenames and the applied migration history before applying
pending files in numeric order. `public.schema_migrations` records each version,
filename, SHA-256 checksum and application timestamp. Repeated runs skip applied
files. Checksums normalize CRLF/LF and a UTF-8 BOM, but other edits or renames to
applied files cause an error. Missing applied files and newly inserted versions
before already applied migrations also cause an error.

Each SQL file and its history record commit in one transaction. A failure rolls
back that file and stops execution; earlier successful files stay committed.
A PostgreSQL advisory lock serializes concurrent runners. Lock waits, including
DDL locks, time out after 60 seconds. Connection timeout is 10 seconds.

Existing databases initialized by the previous Compose configuration can replay
`001` without dropping tables or data: its table and index creation is conditional.
With no history table, the runner executes all files and records their results.
This supports the repository's previous initial schema, not arbitrary manual
schema changes. It does not compare the live schema against SQL definitions or
remove old redundant indexes.

Add new files as `005_description.sql`, `006_description.sql`, etc. Never edit
files already recorded in the history. Files must contain transactional SQL,
without explicit BEGIN/COMMIT/ROLLBACK, psql commands or CREATE INDEX CONCURRENTLY.
Do not manually insert history records to bypass a failed migration.

For an update with an already running backend, stop it before applying migrations:

```sh
docker compose stop frontend backend
docker compose build migrations
docker compose run --rm migrations
# Continue only if the migration command succeeds:
docker compose up -d --build
```

Compose startup dependencies do not stop backend instances that are already running.

For local development, install `requirements.txt`, set `DATABASE_URL` and
`MIGRATIONS_DIR`, then run `python migrate.py`. Tests use the standard library:

```sh
python -m unittest discover -s tests -v
```

Integration tests run only when `TEST_DATABASE_URL` is set. They create and drop
isolated, randomly named schemas in that test database; use a dedicated test DB.
