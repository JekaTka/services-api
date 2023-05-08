-- name: CreateServiceVersion :one
INSERT INTO service_versions (
    changelog, version, service_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetVersionsByServiceID :many
SELECT *
FROM service_versions
WHERE service_id = $1;

-- name: DeleteAllServiceVersions :exec
DELETE FROM service_versions WHERE 1 = 1;

-- name: GetServiceVersionsCount :one
SELECT COUNT(id) AS total
FROM service_versions
WHERE service_id = $1;
