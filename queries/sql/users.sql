-- name: CreateUser :one
INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?) RETURNING *;
 
-- name: GetUserByUsernameAndPassword :one
SELECT * FROM users WHERE username = ? AND password_hash = ?;

-- name: ChangePassword :execrows
UPDATE users SET password_hash = ? WHERE username = ?;
