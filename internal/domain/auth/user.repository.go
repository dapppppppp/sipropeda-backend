package auth

import (
	"bytes"
	"sipropeda-backend/infras"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type UserRepository interface {
	GetByUsername(username string) (User, error)
	Create(data User) error
	ResolveAll(req model.StandardRequest, roleId string) (pagination.Response, error)
	ResolveByID(id uuid.UUID) (User, error)
	Update(data User) error
	Delete(data User) error
	UpdatePassword(id uuid.UUID, hashedPassword string) error
}

type userRepository struct {
	db *infras.PostgresqlConn
}

func ProvideUserRepository(db *infras.PostgresqlConn) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(username string) (User, error) {
	var user User
	// Query menggunakan u.nama
	query := `
		SELECT u.id, u.nama, u.email, u.username, u.password, u.role_id, r.name as role_name 
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 AND u.is_deleted = false
	`
	err := r.db.Read.Get(&user, query, username)
	return user, err
}

func (r *userRepository) Create(data User) error {
	// Insert ke kolom nama
	query := `INSERT INTO users (id, nama, email, username, password, role_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Write.Exec(query, data.ID, data.Name, data.Email, data.Username, data.Password, data.RoleID, data.CreatedAt)
	return err
}

func (r *userRepository) ResolveAll(req model.StandardRequest, roleIdFilter string) (response pagination.Response, err error) {
	var searchParams []interface{}
	var queryBuff bytes.Buffer

	// Select u.nama
	querySelect := `SELECT u.id, u.nama, u.email, u.username, u.role_id, r.name as role_name, u.created_at, u.updated_at 
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id `
		
	queryBuff.WriteString(" WHERE coalesce(u.is_deleted, false) = ? ")
	searchParams = append(searchParams, false)

	if req.Keyword != "" {
		// Pencarian menggunakan u.nama
		queryBuff.WriteString(" AND (u.nama ILIKE ? OR u.username ILIKE ? OR u.email ILIKE ? OR r.name ILIKE ?) ")
		keyword := "%" + req.Keyword + "%"
		searchParams = append(searchParams, keyword, keyword, keyword, keyword)
	}

	if roleIdFilter != "" {
		queryBuff.WriteString(" AND u.role_id = ? ")
		searchParams = append(searchParams, roleIdFilter)
	}

	queryCount := r.db.Read.Rebind("SELECT count(*) FROM users u LEFT JOIN roles r ON u.role_id = r.id " + queryBuff.String())
	var totalData int
	err = r.db.Read.QueryRow(queryCount, searchParams...).Scan(&totalData)
	if err != nil {
		return
	}

	if totalData == 0 {
		response.Items = make([]interface{}, 0)
		response.Meta = pagination.CreateMeta(0, req.PageSize, req.PageNumber)
		return
	}

	queryBuff.WriteString(" ORDER BY u.created_at DESC ")
	offset := (req.PageNumber - 1) * req.PageSize
	queryBuff.WriteString(" LIMIT ? OFFSET ? ")
	searchParams = append(searchParams, req.PageSize, offset)

	finalQuery := r.db.Read.Rebind(querySelect + queryBuff.String())
	rows, err := r.db.Read.Queryx(finalQuery, searchParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	var items []interface{}
	for rows.Next() {
		var user User
		if err = rows.StructScan(&user); err != nil {
			return
		}
		items = append(items, user)
	}

	response.Items = items
	response.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *userRepository) ResolveByID(id uuid.UUID) (User, error) {
	var data User
	// Select u.nama
	query := `
		SELECT u.id, u.nama, u.email, u.username, u.password, u.role_id, r.name as role_name, u.created_at, u.updated_at 
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1 AND u.is_deleted = false
	`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *userRepository) Update(data User) error {
	// Update kolom nama
	query := `
		UPDATE users 
		SET nama = $1, email = $2, username = $3, password = $4, role_id = $5, updated_at = $6
		WHERE id = $7 AND is_deleted = false
	`
	_, err := r.db.Write.Exec(query, data.Name, data.Email, data.Username, data.Password, data.RoleID, data.UpdatedAt, data.ID)
	return err
}

func (r *userRepository) Delete(data User) error {
	query := `UPDATE users SET is_deleted = $1, deleted_at = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedAt, data.ID)
	return err
}

func (r *userRepository) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Write.Exec(query, hashedPassword, id)
	return err
}