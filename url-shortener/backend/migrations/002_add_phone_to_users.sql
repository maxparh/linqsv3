-- Migration 002: Add a unique phone number to users.

ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20) UNIQUE;
