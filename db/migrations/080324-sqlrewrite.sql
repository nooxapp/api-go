-- FUCK UUID (breaks auth) CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id)
);

CREATE TYPE request_status AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE friend_requests (
    id SERIAL PRIMARY KEY,
    requester_id INTEGER NOT NULL REFERENCES users(id),
    requestee_id INTEGER NOT NULL REFERENCES users(id),
    status request_status DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE friends (
    user_id INTEGER NOT NULL REFERENCES users(id),
    friend_id INTEGER NOT NULL REFERENCES users(id),
    PRIMARY KEY (user_id, friend_id)
);