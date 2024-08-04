CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id)
);

CREATE TYPE request_status AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE friend_requests (
    status request_status DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE friends (
)
