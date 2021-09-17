CREATE TABLE IF NOT EXISTS sessions(
    id BIGSERIAL PRIMARY KEY,
    refresh_token VARCHAR(99),
    expires_at TIMESTAMP,
    user_id BIGINT NOT NULL,
    CONSTRAINT oto_session_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create data entries in fresh table for all users
INSERT INTO sessions(user_id)
SELECT u.id FROM users u;