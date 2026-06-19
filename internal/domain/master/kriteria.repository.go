package master

import (
	"bytes"
	"sipropeda-backend/infras"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type KriteriaRepository interface {
	Create(data Kriteria) error
	ResolveAll() ([]Kriteria, error) // Tetap ada untuk drop-down dan kalkulasi bobot
	ResolvePaging(req model.StandardRequest) (pagination.Response, error) // Endpoint Pagination
	ResolveByID(id uuid.UUID) (Kriteria, error)
	Update(data Kriteria) error
	Delete(data Kriteria) error
}

type kriteriaRepository struct {
	db *infras.PostgresqlConn
}

func ProvideKriteriaRepository(db *infras.PostgresqlConn) KriteriaRepository {
	return &kriteriaRepository{db: db}
}

func (r *kriteriaRepository) Create(data Kriteria) error {
	query := `
		INSERT INTO m_kriteria (id, kode, nama, jenis, bobot, is_active, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Write.Exec(query, data.ID, data.Kode, data.Nama, data.Jenis, data.Bobot, data.IsActive, data.CreatedBy, data.CreatedAt)
	return err
}

func (r *kriteriaRepository) ResolveAll() ([]Kriteria, error) {
	var data []Kriteria
	query := `SELECT * FROM m_kriteria WHERE is_deleted = false ORDER BY kode ASC`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *kriteriaRepository) ResolvePaging(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var filterBuff bytes.Buffer

	filterBuff.WriteString(" WHERE is_deleted = false")

	if req.Keyword != "" {
		filterBuff.WriteString(" AND (nama ILIKE ? OR kode ILIKE ?) ")
		searchParams = append(searchParams, "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	selectDto := `SELECT * FROM m_kriteria`

	// 1. Hitung Total Data
	countQuery := r.db.Read.Rebind("SELECT count(*) FROM (" + selectDto + filterBuff.String() + ")x")
	var totalData int
	err = r.db.Read.QueryRowx(countQuery, searchParams...).Scan(&totalData)
	if err != nil {
		return
	}

	if totalData < 1 {
		data.Items = []Kriteria{}
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return data, nil
	}

	// 2. Sorting Dinamis
	orderByColumn, ok := ColumnMapKriteria[req.SortBy].(string)
	if !ok {
		orderByColumn = "kode" // Default sort by Kode C1, C2
	}
	
	sortType := req.SortType
	if sortType == "" {
		sortType = "asc"
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

	var items []Kriteria
	for rows.Next() {
		var item Kriteria
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

func (r *kriteriaRepository) ResolveByID(id uuid.UUID) (Kriteria, error) {
	var data Kriteria
	query := `SELECT * FROM m_kriteria WHERE id = $1 AND is_deleted = false`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *kriteriaRepository) Update(data Kriteria) error {
	query := `
		UPDATE m_kriteria 
		SET kode = $1, nama = $2, jenis = $3, bobot = $4, updated_by = $5, updated_at = $6
		WHERE id = $7 AND is_deleted = false
	`
	_, err := r.db.Write.Exec(query, data.Kode, data.Nama, data.Jenis, data.Bobot, data.UpdatedBy, data.UpdatedAt, data.ID)
	return err
}

func (r *kriteriaRepository) Delete(data Kriteria) error {
	query := `
		UPDATE m_kriteria 
		SET is_deleted = $1, deleted_at = $2, updated_by = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedBy, data.UpdatedAt, data.ID)
	return err
}