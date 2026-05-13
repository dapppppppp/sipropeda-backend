package auth

import (
	"sipropeda-backend/infras"
	"github.com/gofrs/uuid"
)

type RoleRepository interface {
	Create(data Role) error
	ResolveAll() ([]Role, error)
	ResolveByID(id uuid.UUID) (Role, error)
	Update(data Role) error
	Delete(data Role) error
}

type roleRepository struct {
	db *infras.PostgresqlConn
}

func ProvideRoleRepository(db *infras.PostgresqlConn) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(data Role) error {
	query := `INSERT INTO roles (id, name, description, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Write.Exec(query, data.ID, data.Name, data.Description, data.CreatedAt)
	return err
}

func (r *roleRepository) ResolveAll() ([]Role, error) {
	var data []Role
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE is_deleted = false ORDER BY name ASC`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *roleRepository) ResolveByID(id uuid.UUID) (Role, error) {
	var data Role
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1 AND is_deleted = false`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *roleRepository) Update(data Role) error {
	query := `UPDATE roles SET name = $1, description = $2, updated_at = $3 WHERE id = $4 AND is_deleted = false`
	_, err := r.db.Write.Exec(query, data.Name, data.Description, data.UpdatedAt, data.ID)
	return err
}

func (r *roleRepository) Delete(data Role) error {
	query := `UPDATE roles SET is_deleted = $1, deleted_at = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedAt, data.ID)
	return err
}