-- name: GetPolicies :many
SELECT * FROM policies;

-- name: GetPoliciesByUsername :many
SELECT p.* FROM policies p INNER JOIN user_policies up ON up.policy_id = p.id INNER JOIN users u ON up.user_id = u.id WHERE u.username = ?;

-- name: GetPolicyByName :one
SELECT * FROM policies WHERE name = ?;

-- name: GetPolicyByID :one
SELECT * FROM policies WHERE id = ?;

-- name: AddPolicy :one
INSERT INTO policies (name, description, auto_sign, subject_cn_regex, san_dns_regex, max_requested_days, created_at) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *;
