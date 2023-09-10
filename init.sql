CREATE DATABASE jwt_test;

\c jwt_test;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE jwt (
    id SERIAL PRIMARY KEY,
    expiration DATE NOT NULL,
    token TEXT NOT null,
    user_id INTEGER REFERENCES users(id)
);
