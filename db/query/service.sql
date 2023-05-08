-- name: CreateService :one
INSERT INTO services (
    name, description
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListServices :many
SELECT *
FROM services
WHERE
    name LIKE CASE WHEN @search::text IS NULL THEN '' ELSE concat('%', @search, '%')::text END
ORDER BY (
    CASE WHEN @order_by::text = 'name' THEN name END,
    CASE WHEN @order_by::text = 'description' THEN description END
)
LIMIT $1
OFFSET $2;

-- name: GetServiceByID :one
SELECT * FROM services
WHERE id = $1 LIMIT 1;

-- name: DeleteAllServices :exec
DELETE FROM services WHERE 1 = 1;

-- name: GetServicesCount :one
SELECT COUNT(id) AS total
FROM services;
