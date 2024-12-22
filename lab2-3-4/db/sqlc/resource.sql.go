// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: resource.sql

package db

import (
	"context"
)

const createResource = `-- name: CreateResource :one
INSERT INTO resources (name, org_id)
VALUES (?, ?)
RETURNING id, name, org_id, created_at
`

type CreateResourceParams struct {
	Name  string `json:"name"`
	OrgID int64  `json:"org_id"`
}

func (q *Queries) CreateResource(ctx context.Context, arg CreateResourceParams) (Resource, error) {
	row := q.db.QueryRowContext(ctx, createResource, arg.Name, arg.OrgID)
	var i Resource
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OrgID,
		&i.CreatedAt,
	)
	return i, err
}

const getResource = `-- name: GetResource :one
SELECT id, name, org_id, created_at
FROM resources
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetResource(ctx context.Context, id int64) (Resource, error) {
	row := q.db.QueryRowContext(ctx, getResource, id)
	var i Resource
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OrgID,
		&i.CreatedAt,
	)
	return i, err
}

const listOrganizationResources = `-- name: ListOrganizationResources :many
SELECT id, name, org_id, created_at
FROM resources
WHERE org_id = ?
ORDER BY created_at
`

func (q *Queries) ListOrganizationResources(ctx context.Context, orgID int64) ([]Resource, error) {
	rows, err := q.db.QueryContext(ctx, listOrganizationResources, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Resource{}
	for rows.Next() {
		var i Resource
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OrgID,
			&i.CreatedAt,
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