import hashlib
import logging
import os
import re
from dataclasses import dataclass
from pathlib import Path

import psycopg


LOG = logging.getLogger("migrations")
LOCK_ID = 764913208
FILE_PATTERN = re.compile(r"([0-9]+)_[a-z0-9_]+\.sql")


class MigrationError(Exception):
    pass


@dataclass(frozen=True)
class Migration:
    version: int
    name: str
    checksum: str
    sql: str


def load_migrations(directory: Path) -> list[Migration]:
    if not directory.is_dir():
        raise MigrationError(f"Migration directory does not exist: {directory}")

    migrations = []
    versions = set()
    for path in sorted(directory.glob("*.sql")):
        match = FILE_PATTERN.fullmatch(path.name)
        if not match:
            raise MigrationError(f"Invalid migration filename: {path.name}")
        version = int(match.group(1))
        if not 0 < version < 2**63 or version in versions:
            raise MigrationError(f"Invalid or duplicate migration version: {path.name}")
        # Normalize line endings and a possible UTF-8 BOM across host platforms.
        sql = path.read_text(encoding="utf-8-sig")
        if not sql.strip():
            raise MigrationError(f"Empty migration: {path.name}")
        checksum = hashlib.sha256(sql.encode("utf-8")).hexdigest()
        migrations.append(Migration(version, path.name, checksum, sql))
        versions.add(version)

    if not migrations:
        raise MigrationError("No SQL migrations found")
    return sorted(migrations, key=lambda migration: migration.version)


def pending_migrations(
    migrations: list[Migration], applied: list[tuple[int, str, str]]
) -> list[Migration]:
    by_version = {migration.version: migration for migration in migrations}
    for version, name, checksum in applied:
        migration = by_version.get(version)
        if migration is None:
            raise MigrationError(f"Applied migration is missing from the image: {name}")
        if (name, checksum) != (migration.name, migration.checksum):
            raise MigrationError(f"Applied migration has changed: {name}")

    if [row[0] for row in applied] != [m.version for m in migrations[:len(applied)]]:
        raise MigrationError("Migration history is out of order; add new versions at the end")
    return migrations[len(applied):]


def run_migrations(connection, migrations: list[Migration]) -> int:
    # The session lock covers history validation and all per-file transactions.
    connection.execute("SET lock_timeout = '60s'")
    LOG.info("Waiting for the migration lock")
    connection.execute("SELECT pg_advisory_lock(%s)", (LOCK_ID,))
    try:
        connection.execute("""
            CREATE TABLE IF NOT EXISTS schema_migrations (
                version BIGINT PRIMARY KEY,
                name TEXT NOT NULL UNIQUE,
                checksum TEXT NOT NULL,
                applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
            )
        """)
        applied = connection.execute(
            "SELECT version, name, checksum FROM schema_migrations ORDER BY version"
        ).fetchall()
        pending = pending_migrations(migrations, applied)
        LOG.info("Validated %s applied migrations; %s pending", len(applied), len(pending))

        for migration in pending:
            LOG.info("Applying %s", migration.name)
            with connection.transaction():
                connection.execute(migration.sql, prepare=False)
                connection.execute(
                    "INSERT INTO schema_migrations (version, name, checksum) VALUES (%s, %s, %s)",
                    (migration.version, migration.name, migration.checksum),
                )
            LOG.info("Applied %s", migration.name)
        return len(pending)
    finally:
        if not connection.closed:
            connection.execute("SELECT pg_advisory_unlock(%s)", (LOCK_ID,))


def main() -> int:
    logging.basicConfig(level=logging.INFO, format="%(levelname)s %(message)s")
    try:
        database_url = os.environ.get("DATABASE_URL")
        if not database_url:
            raise MigrationError("DATABASE_URL is required")
        directory = Path(os.environ.get("MIGRATIONS_DIR", "/app/migrations"))
        migrations = load_migrations(directory)
        with psycopg.connect(
            database_url,
            autocommit=True,
            connect_timeout=10,
            application_name="linqs-migrations",
            options="-c search_path=public",
        ) as connection:
            count = run_migrations(connection, migrations)
        LOG.info("Migration check completed successfully: %s applied", count)
        return 0
    except psycopg.OperationalError as error:
        LOG.error("PostgreSQL connection or operation failed (SQLSTATE: %s)", error.sqlstate)
    except psycopg.Error as error:
        LOG.error("PostgreSQL error [%s]: %s", error.sqlstate, error.diag.message_primary)
    except (MigrationError, OSError, UnicodeError) as error:
        LOG.error("%s", error)
    return 1


if __name__ == "__main__":
    raise SystemExit(main())
