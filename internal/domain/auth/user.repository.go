package auth

import (
	"sipropeda-backend/infras"

	"github.com/gofrs/uuid"
)

type UserRepository interface {
	GetByUsername(username string) (User, error)
	Create(data User) error
	ResolveAll() ([]User, error)
	ResolveByID(id uuid.UUID) (User, error)
	Update(data User) error
	Delete(data User) error
}

type userRepository struct {
	db *infras.PostgresqlConn
}

func ProvideUserRepository(db *infras.PostgresqlConn) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(username string) (User, error) {
	var user User
	query := `
		SELECT u.id, u.username, u.password, u.role_id, r.name as role_name 
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 AND u.is_deleted = false
	`
	err := r.db.Read.Get(&user, query, username)
	return user, err
}

func (r *userRepository) Create(data User) error {
	query := `INSERT INTO users (id, username, password, role_id, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Write.Exec(query, data.ID, data.Username, data.Password, data.RoleID, data.CreatedAt)
	return err
}

func (r *userRepository) ResolveAll() ([]User, error) {
	var data []User
	query := `
		SELECT u.id, u.username, u.role_id, r.name as role_name, u.created_at, u.updated_at 
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.is_deleted = false ORDER BY u.created_at DESC
	`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *userRepository) ResolveByID(id uuid.UUID) (User, error) {
	var data User
	query := `
		SELECT u.id, u.username, u.password, u.role_id, r.name as role_name, u.created_at, u.updated_at 
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1 AND u.is_deleted = false
	`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *userRepository) Update(data User) error {
	query := `
		UPDATE users 
		SET username = $1, password = $2, role_id = $3, updated_at = $4
		WHERE id = $5 AND is_deleted = false
	`
	_, err := r.db.Write.Exec(query, data.Username, data.Password, data.RoleID, data.UpdatedAt, data.ID)
	return err
}

func (r *userRepository) Delete(data User) error {
	query := `UPDATE users SET is_deleted = $1, deleted_at = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedAt, data.ID)
	return err
}