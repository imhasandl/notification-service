-- +goose Up
CREATE TABLE users (
   id UUID NOT NULL PRIMARY KEY,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP NOT NULL,
   email TEXT NOT NULL UNIQUE,
   password TEXT NOT NULL,
   username TEXT NOT NULL,
   is_premium BOOLEAN NOT NULL DEFAULT FALSE,
   verification_code INT NOT NULL,
   is_verified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE refresh_tokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    expiry_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE refresh_tokens;
DROP INDEX idx_users_username;
DROP TABLE users;