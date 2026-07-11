CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name          text        NOT NULL,
    username      text        NOT NULL,
    email         text        NOT NULL,
    password_hash text        NOT NULL,
    bio           text        NOT NULL DEFAULT '',
    avatar_id     text        NOT NULL DEFAULT 'a1',
    banner_id     text        NOT NULL DEFAULT 'b1',
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

-- Case-insensitive uniqueness for email and username.
CREATE UNIQUE INDEX idx_users_email ON users (lower(email));
CREATE UNIQUE INDEX idx_users_username ON users (lower(username));
