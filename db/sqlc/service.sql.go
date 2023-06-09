// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: service.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createService = `-- name: CreateService :one
INSERT INTO services (
    name, description
) VALUES (
    $1, $2
) RETURNING id, name, description, created_at, updated_at, deleted_at
`

type CreateServiceParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) CreateService(ctx context.Context, arg CreateServiceParams) (Service, error) {
	row := q.db.QueryRowContext(ctx, createService, arg.Name, arg.Description)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteAllServices = `-- name: DeleteAllServices :exec
DELETE FROM services WHERE 1 = 1
`

func (q *Queries) DeleteAllServices(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllServices)
	return err
}

const getServiceByID = `-- name: GetServiceByID :one
SELECT id, name, description, created_at, updated_at, deleted_at FROM services
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetServiceByID(ctx context.Context, id uuid.UUID) (Service, error) {
	row := q.db.QueryRowContext(ctx, getServiceByID, id)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getServicesCount = `-- name: GetServicesCount :one
SELECT COUNT(id) AS total
FROM services
`

func (q *Queries) GetServicesCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getServicesCount)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const listServices = `-- name: ListServices :many
SELECT id, name, description, created_at, updated_at, deleted_at
FROM services
WHERE
    name LIKE CASE WHEN $3::text IS NULL THEN '' ELSE concat('%', $3, '%')::text END
ORDER BY (
    CASE WHEN $4::text = 'name' THEN name END,
    CASE WHEN $4::text = 'description' THEN description END
)
LIMIT $1
OFFSET $2
`

type ListServicesParams struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

func (q *Queries) ListServices(ctx context.Context, arg ListServicesParams) ([]Service, error) {
	rows, err := q.db.QueryContext(ctx, listServices,
		arg.Limit,
		arg.Offset,
		arg.Search,
		arg.OrderBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Service
	for rows.Next() {
		var i Service
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
