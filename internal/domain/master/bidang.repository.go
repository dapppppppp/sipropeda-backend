package master

import (
	"sipropeda-backend/infras"

	"github.com/gofrs/uuid"
)

type BidangPembangunanRepository interface {
	Create(data BidangPembangunan) error
	ResolveAll() ([]BidangPembangunan, error)
	ResolveByID(id uuid.UUID) (BidangPembangunan, error)
	Update(data BidangPembangunan) error
	Delete(data BidangPembangunan) error
}

type bidangPembangunanRepository struct {
	db *infras.PostgresqlConn
}

func ProvideBidangPembangunanRepository(db *infras.PostgresqlConn) BidangPembangunanRepository {
	return &bidangPembangunanRepository{db: db}
}

func (r *bidangPembangunanRepository) Create(data BidangPembangunan) error {
	query := `
		INSERT INTO bidang_pembangunan (id, nama_bidang, created_by, created_at) 
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Write.Exec(query, data.ID, data.NamaBidang, data.CreatedBy, data.CreatedAt)
	return err
}

func (r *bidangPembangunanRepository) ResolveAll() ([]BidangPembangunan, error) {
	var data []BidangPembangunan
	query := `
		SELECT id, nama_bidang, created_at, updated_at 
		FROM bidang_pembangunan 
		WHERE is_deleted = false 
		ORDER BY nama_bidang ASC
	`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *bidangPembangunanRepository) ResolveByID(id uuid.UUID) (BidangPembangunan, error) {
	var data BidangPembangunan
	query := `
		SELECT id, nama_bidang, created_at, updated_at 
		FROM bidang_pembangunan 
		WHERE id = $1 AND is_deleted = false
	`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *bidangPembangunanRepository) Update(data BidangPembangunan) error {
	query := `
		UPDATE bidang_pembangunan 
		SET nama_bidang = $1, updated_by = $2, updated_at = $3 
		WHERE id = $4 AND is_deleted = false
	`
	_, err := r.db.Write.Exec(query, data.NamaBidang, data.UpdatedBy, data.UpdatedAt, data.ID)
	return err
}

func (r *bidangPembangunanRepository) Delete(data BidangPembangunan) error {
	query := `
		UPDATE bidang_pembangunan 
		SET is_deleted = $1, deleted_at = $2, updated_by = $3, updated_at = $4 
		WHERE id = $5
	`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedBy, data.UpdatedAt, data.ID)
	return err
}