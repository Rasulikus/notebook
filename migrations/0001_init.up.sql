-- Таблица users
CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    email        TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name         TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    deleted_at    TIMESTAMPTZ
);

-- Таблица notes
CREATE TABLE notes (
    id        BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id   BIGINT REFERENCES users(id),
    title     TEXT NOT NULL,
    text      TEXT

);

-- Таблица tags
CREATE TABLE tags (
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    name    TEXT NOT NULL,
    UNIQUE (user_id, name)
);

-- Связка многие-ко-многим
CREATE TABLE notes_tags (
    note_id BIGINT REFERENCES notes(id) ON DELETE CASCADE,
    tag_id  BIGINT REFERENCES tags(id),
    PRIMARY KEY (note_id, tag_id)
);

-- Таблица sessions
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash BYTEA UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
)
