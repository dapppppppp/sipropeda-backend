package master

import (
	"bytes"
	"sipropeda-backend/infras"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type SumberDanaRepository interface {
	Create(data SumberDana) error
	ResolveAll() ([]SumberDana, error) // Tetap utuh tanpa limit untuk Dropdown halaman lain
	ResolvePaging(req model.StandardRequest) (pagination.Response, error) // Fungsi baru khusus tabel master
	ResolveByID(id uuid.UUID) (SumberDana, error)
	Update(data SumberDana) error
	Delete(data SumberDana) error
}

type sumberDanaRepository struct {
	db *infras.PostgresqlConn
}

func ProvideSumberDanaRepository(db *infras.PostgresqlConn) SumberDanaRepository {
	return &sumberDanaRepository{db: db}
}

func (r *sumberDanaRepository) Create(data SumberDana) error {
	query := `INSERT INTO sumber_dana (id, nama_sumber, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.Write.Exec(query, data.ID, data.NamaSumber, data.CreatedAt)
	return err
}

func (r *sumberDanaRepository) ResolveAll() ([]SumberDana, error) {
	var data []SumberDana
	query := `SELECT id, nama_sumber, created_at, updated_at FROM sumber_dana WHERE is_deleted = false ORDER BY nama_sumber ASC`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *sumberDanaRepository) ResolvePaging(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var filterBuff bytes.Buffer

	filterBuff.WriteString(" WHERE coalesce(is_deleted, false) = false")

	if req.Keyword != "" {
		filterBuff.WriteString(" AND nama_sumber ILIKE ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	selectDto := `SELECT id, nama_sumber, created_at, updated_at FROM sumber_dana`

	// 1. Hitung Total Data
	countQuery := r.db.Read.Rebind("SELECT count(*) FROM (" + selectDto + filterBuff.String() + ")x")
	var totalData int
	err = r.db.Read.QueryRowx(countQuery, searchParams...).Scan(&totalData)
	if err != nil {
		return
	}

	if totalData < 1 {
		data.Items = []SumberDana{}
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return data, nil
	}

	// 2. Sorting Dinamis
	orderByColumn, ok := ColumnMapSumberDana[req.SortBy].(string)
	if !ok {
		orderByColumn = "created_at"
	}
	filterBuff.WriteString(" ORDER BY " + orderByColumn + " " + req.SortType)

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

	var items []SumberDana
	for rows.Next() {
		var item SumberDana
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

func (r *sumberDanaRepository) ResolveByID(id uuid.UUID) (SumberDana, error) {
	var data SumberDana
	query := `SELECT id, nama_sumber, created_at, updated_at FROM sumber_dana WHERE id = $1 AND is_deleted = false`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *sumberDanaRepository) Update(data SumberDana) error {
	query := `UPDATE sumber_dana SET nama_sumber = $1, updated_at = $2 WHERE id = $3 AND is_deleted = false`
	_, err := r.db.Write.Exec(query, data.NamaSumber, data.UpdatedAt, data.ID)
	return err
}

func (r *sumberDanaRepository) Delete(data SumberDana) error {
	query := `UPDATE sumber_dana SET is_deleted = $1, deleted_at = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedAt, data.ID)
	return err
}