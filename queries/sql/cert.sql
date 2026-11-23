-- name: ListCSRs :many
SELECT * FROM csrs;

-- name: ListPendingCSRs :many
SELECT * FROM csrs WHERE status = 'pending';

-- name: ListApprovedCSRs :many
SELECT * FROM csrs WHERE status = 'approved';

-- name: ListRejectedCSRs :many
SELECT * FROM csrs WHERE status = 'rejected';

---- name: AddCSR :one
--INSERT INTO csrs () VALUES ();
