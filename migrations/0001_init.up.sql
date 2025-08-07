BEGIN;

-- Таблица users
CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    email        TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    username     TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ 
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

ALTER TABLE asdf ADD COLUMN asd TEXT;

COMMIT;