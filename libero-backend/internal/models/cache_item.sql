-- SQL migration for cache_items table (PostgreSQL)
CREATE TABLE IF NOT EXISTS cache_items (
    id SERIAL PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    value BYTEA NOT NULL,
    e_tag TEXT,
    last_modified TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
