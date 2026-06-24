package auth

import (
	"bytes"
	"sipropeda-backend/infras"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type RoleRepository interface {
	Create(data Role) error
	ResolveAll() ([]Role, error)
	ResolvePaging(req model.StandardRequest) (pagination.Response, error) // Endpoint Pagination
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

func (r *roleRepository) ResolvePaging(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var filterBuff bytes.Buffer

	filterBuff.WriteString(" WHERE is_deleted = false")

	if req.Keyword != "" {
		filterBuff.WriteString(" AND (name ILIKE ? OR description ILIKE ?) ")
		keyword := "%" + req.Keyword + "%"
		searchParams = append(searchParams, keyword, keyword)
	}

	selectDto := `SELECT id, name, description, created_at, updated_at FROM roles`

	// 1. Hitung Total Data
	countQuery := r.db.Read.Rebind("SELECT count(*) FROM (" + selectDto + filterBuff.String() + ")x")
	var totalData int
	err = r.db.Read.QueryRowx(countQuery, searchParams...).Scan(&totalData)
	if err != nil {
		return
	}

	if totalData < 1 {
		data.Items = []Role{}
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return data, nil
	}

	// 2. Sorting Dinamis
	orderByColumn, ok := ColumnMapRole[req.SortBy].(string)
	if !ok {
		orderByColumn = "created_at"
	}
	sortType := req.SortType
	if sortType == "" {
		sortType = "desc"
	}
	filterBuff.WriteString(" ORDER BY " + orderByColumn + " " + sortType)

	// 3. Limit Offset
	offset := (req.PageNumber - 1) * req.PageSize
	filterBuff.WriteString(" LIMIT ? OFFSET ? ")
	searchParams = append(searchParams, req.PageSize, offset)

	searchQuery := r.db.Read.Rebind(selectDto + filterBuff.String())
	rows, err := r.db.Read.Queryx(searchQuery, searchParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	var items []Role
	for rows.Next() {
		var item Role
		err = rows.StructScan(&item)
		if err != nil {
			return
		}
		items = append(items, item)
	}

	data.Items = items
	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return data, nil
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