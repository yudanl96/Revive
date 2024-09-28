-- name: CreateAccount :execresult
INSERT INTO accounts(
    id, username
) VALUES (uuid(), ?);

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id;
