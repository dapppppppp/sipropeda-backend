package transaction

import (
	"bytes"
	"sipropeda-backend/infras"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type PaguAnggaranRepository interface {
	Create(data PaguAnggaran) error
	ResolveAll(req model.StandardRequest) (pagination.Response, error)
	ResolveByID(id uuid.UUID) (PaguAnggaran, error)
	Update(data PaguAnggaran) error
	Delete(data PaguAnggaran) error
}

type paguAnggaranRepository struct {
	db *infras.PostgresqlConn
}

func ProvidePaguAnggaranRepository(db *infras.PostgresqlConn) PaguAnggaranRepository {
	return &paguAnggaranRepository{db: db}
}

func (r *paguAnggaranRepository) Create(data PaguAnggaran) error {
	query := `
		INSERT INTO pagu_anggaran (id, tahun, sumber_dana_id, pagu_estimasi, pagu_definitif, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Write.Exec(query, data.ID, data.Tahun, data.SumberDanaID, data.PaguEstimasi, data.PaguDefinitif, data.CreatedBy, data.CreatedAt)
	return err
}

func (r *paguAnggaranRepository) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var filterBuff bytes.Buffer

	filterBuff.WriteString(" WHERE coalesce(p.is_deleted, false) = false")

	// Filter Pencarian
	if req.Keyword != "" {
		filterBuff.WriteString(" AND ")
		filterBuff.WriteString(" concat(p.tahun::text, s.nama_sumber) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.Tahun != "" {
		filterBuff.WriteString(" AND p.tahun = ? ")
		searchParams = append(searchParams, req.Tahun)
	}

	selectDto := `
		SELECT p.id, p.tahun, p.sumber_dana_id, s.nama_sumber as sumber_dana_name, p.pagu_estimasi, p.pagu_definitif, p.created_at 
		FROM pagu_anggaran p
		LEFT JOIN sumber_dana s ON p.sumber_dana_id = s.id
	`

	// 1. Hitung Total
	countQuery := r.db.Read.Rebind("SELECT count(*) FROM (" + selectDto + filterBuff.String() + ")x")
	var totalData int
	err = r.db.Read.QueryRowx(countQuery, searchParams...).Scan(&totalData)
	if err != nil {
		return
	}

	if totalData < 1 {
		data.Items = []PaguAnggaran{}
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return data, nil
	}

	// 2. Sorting
	orderByColumn, ok := ColumnMapPaguAnggaran[req.SortBy].(string)
	if !ok {
		orderByColumn = "p.created_at"
	}
	filterBuff.WriteString(" ORDER BY " + orderByColumn + " " + req.SortType)

	// 3. Limit Offset Pagination
	offset := (req.PageNumber - 1) * req.PageSize
	filterBuff.WriteString(" LIMIT ? OFFSET ? ")
	searchParams = append(searchParams, req.PageSize, offset)

	searchQuery := r.db.Read.Rebind(selectDto + filterBuff.String())
	rows, err := r.db.Read.Queryx(searchQuery, searchParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	var items []PaguAnggaran
	for rows.Next() {
		var item PaguAnggaran
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

func (r *paguAnggaranRepository) ResolveByID(id uuid.UUID) (PaguAnggaran, error) {
	var data PaguAnggaran
	query := `
		SELECT p.id, p.tahun, p.sumber_dana_id, s.nama_sumber as sumber_dana_name, p.pagu_estimasi, p.pagu_definitif, p.created_at 
		FROM pagu_anggaran p
		LEFT JOIN sumber_dana s ON p.sumber_dana_id = s.id
		WHERE p.id = $1 AND p.is_deleted = false
	`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *paguAnggaranRepository) Update(data PaguAnggaran) error {
	query := `
		UPDATE pagu_anggaran 
		SET tahun = $1, sumber_dana_id = $2, pagu_estimasi = $3, pagu_definitif = $4, updated_by = $5, updated_at = $6 
		WHERE id = $7 AND is_deleted = false
	`
	_, err := r.db.Write.Exec(query, data.Tahun, data.SumberDanaID, data.PaguEstimasi, data.PaguDefinitif, data.UpdatedBy, data.UpdatedAt, data.ID)
	return err
}

func (r *paguAnggaranRepository) Delete(data PaguAnggaran) error {
	query := `
		UPDATE pagu_anggaran 
		SET is_deleted = $1, deleted_at = $2, updated_by = $3, updated_at = $4 
		WHERE id = $5
	`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedBy, data.UpdatedAt, data.ID)
	return err
}