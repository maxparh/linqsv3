import os
import sys
import tempfile
import unittest
import uuid
from pathlib import Path

import psycopg
from psycopg import sql

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from migrate import LOCK_ID, MigrationError, load_migrations, pending_migrations, run_migrations


PROJECT_MIGRATIONS = Path(__file__).resolve().parents[2] / "backend" / "migrations"


class MigrationValidationTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.path = Path(self.directory.name)

    def write(self, name, content="SELECT 1;"):
        path = self.path / name
        path.write_bytes(content.encode("utf-8"))
        return path

    def test_numeric_order_and_pending_files(self):
        self.write("10_later.sql")
        self.write("2_first.sql")
        migrations = load_migrations(self.path)
        self.assertEqual([m.version for m in migrations], [2, 10])
        first = migrations[0]
        applied = [(first.version, first.name, first.checksum)]
        self.assertEqual(pending_migrations(migrations, applied), migrations[1:])

    def test_duplicate_versions_rejected(self):
        self.write("001_first.sql")
        self.write("1_duplicate.sql")
        with self.assertRaisesRegex(MigrationError, "duplicate"):
            load_migrations(self.path)

    def test_invalid_or_empty_input_rejected(self):
        with self.assertRaisesRegex(MigrationError, "No SQL"):
            load_migrations(self.path)
        for name, content in [("init.sql", "SELECT 1;"), ("000_zero.sql", "SELECT 1;"),
                              ("001_empty.sql", " \n")]:
            with self.subTest(name=name):
                path = self.write(name, content)
                with self.assertRaises(MigrationError):
                    load_migrations(self.path)
                path.unlink()

    def test_line_endings_and_bom_do_not_change_checksum(self):
        self.write("001_first.sql", "SELECT 1;\nSELECT 2;\n")
        original = load_migrations(self.path)[0]
        self.write("001_first.sql", "\ufeffSELECT 1;\r\nSELECT 2;\r\n")
        self.assertEqual(original.checksum, load_migrations(self.path)[0].checksum)

    def test_changed_or_renamed_applied_file_rejected(self):
        path = self.write("001_first.sql")
        original = load_migrations(self.path)[0]
        applied = [(original.version, original.name, original.checksum)]
        self.write(path.name, "SELECT 2;")
        with self.assertRaisesRegex(MigrationError, "changed"):
            pending_migrations(load_migrations(self.path), applied)
        self.write(path.name)
        path.rename(self.path / "001_renamed.sql")
        with self.assertRaisesRegex(MigrationError, "changed"):
            pending_migrations(load_migrations(self.path), applied)

    def test_missing_applied_file_rejected(self):
        self.write("001_first.sql")
        migrations = load_migrations(self.path)
        with self.assertRaisesRegex(MigrationError, "missing"):
            pending_migrations(migrations, [(2, "002_missing.sql", "checksum")])

    def test_out_of_order_history_rejected(self):
        self.write("001_first.sql")
        self.write("003_later.sql")
        migrations = load_migrations(self.path)
        later = migrations[1]
        with self.assertRaisesRegex(MigrationError, "out of order"):
            pending_migrations(migrations, [(later.version, later.name, later.checksum)])


@unittest.skipUnless(os.environ.get("TEST_DATABASE_URL"), "TEST_DATABASE_URL is not set")
class PostgreSQLMigrationTests(unittest.TestCase):
    def setUp(self):
        self.connection = psycopg.connect(os.environ["TEST_DATABASE_URL"], autocommit=True)
        self.addCleanup(self.connection.close)
        self.schema = "migration_test_" + uuid.uuid4().hex
        self.connection.execute(sql.SQL("CREATE SCHEMA {}").format(sql.Identifier(self.schema)))
        self.addCleanup(self.drop_schema)
        self.connection.execute(sql.SQL("SET search_path TO {}").format(sql.Identifier(self.schema)))
        self.migrations = load_migrations(PROJECT_MIGRATIONS)

    def drop_schema(self):
        self.connection.execute(sql.SQL("DROP SCHEMA {} CASCADE").format(sql.Identifier(self.schema)))

    def test_fresh_database_and_repeat(self):
        self.assertEqual(run_migrations(self.connection, self.migrations), len(self.migrations))
        history = self.connection.execute("SELECT * FROM schema_migrations ORDER BY version").fetchall()
        self.assertEqual(run_migrations(self.connection, self.migrations), 0)
        self.assertEqual(history, self.connection.execute(
            "SELECT * FROM schema_migrations ORDER BY version"
        ).fetchall())
        self.connection.execute(
            "INSERT INTO users (email, password_hash, phone) VALUES ('test@example.com', 'hash', '+123')"
        )
        self.connection.execute("INSERT INTO analytics_sessions (session_id) VALUES ('test')")
        self.connection.execute("INSERT INTO click_stats (ip_address) VALUES (%s)", ("a" * 64,))

    def test_legacy_initial_schema_preserves_data(self):
        self.connection.execute(self.migrations[0].sql, prepare=False)
        self.connection.execute("CREATE INDEX idx_users_email ON users(email)")
        self.connection.execute("CREATE INDEX idx_links_short_code ON links(short_code)")
        self.connection.execute(
            "INSERT INTO users (email, password_hash) VALUES ('legacy@example.com', 'hash')"
        )
        self.assertEqual(run_migrations(self.connection, self.migrations), len(self.migrations))
        self.assertEqual(self.connection.execute("SELECT email, phone FROM users").fetchall(),
                         [("legacy@example.com", None)])

    def test_failed_file_rolls_back_and_can_be_retried(self):
        run_migrations(self.connection, self.migrations)
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / f"{self.migrations[-1].version + 1:03d}_failure.sql"
            path.write_text("CREATE TABLE rollback_probe (id INT); SELECT 1 / 0;", encoding="utf-8")
            with self.assertRaises(psycopg.errors.DivisionByZero):
                run_migrations(self.connection, self.migrations + load_migrations(Path(directory)))
            self.assertIsNone(self.connection.execute("SELECT to_regclass('rollback_probe')").fetchone()[0])
            self.assertEqual(self.connection.execute("SELECT count(*) FROM schema_migrations").fetchone()[0], len(self.migrations))
            with psycopg.connect(os.environ["TEST_DATABASE_URL"], autocommit=True) as other:
                self.assertTrue(other.execute("SELECT pg_try_advisory_lock(%s)", (LOCK_ID,)).fetchone()[0])
                other.execute("SELECT pg_advisory_unlock(%s)", (LOCK_ID,))
            path.write_text("CREATE TABLE rollback_probe (id INT);", encoding="utf-8")
            self.assertEqual(run_migrations(self.connection, self.migrations + load_migrations(Path(directory))), 1)

    def test_changed_history_stops_before_applying_pending_files(self):
        run_migrations(self.connection, self.migrations[:1])
        self.connection.execute("UPDATE schema_migrations SET checksum = 'changed'")
        with self.assertRaisesRegex(MigrationError, "changed"):
            run_migrations(self.connection, self.migrations)
        self.assertIsNone(self.connection.execute("SELECT to_regclass('analytics_sessions')").fetchone()[0])


if __name__ == "__main__":
    unittest.main()
