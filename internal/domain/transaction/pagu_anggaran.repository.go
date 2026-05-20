package transaction

import (
	"sipropeda-backend/infras"

	"github.com/gofrs/uuid"
)

type PaguAnggaranRepository interface {
	Create(data PaguAnggaran) error
	ResolveAll() ([]PaguAnggaran, error)
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

func (r *paguAnggaranRepository) ResolveAll() ([]PaguAnggaran, error) {
	var data []PaguAnggaran
	query := `
		SELECT p.id, p.tahun, p.sumber_dana_id, s.nama_sumber as sumber_dana_name, p.pagu_estimasi, p.pagu_definitif, p.created_at 
		FROM pagu_anggaran p
		LEFT JOIN sumber_dana s ON p.sumber_dana_id = s.id
		WHERE p.is_deleted = false 
		ORDER BY p.tahun DESC, p.created_at DESC
	`
	err := r.db.Read.Select(&data, query)
	return data, err
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