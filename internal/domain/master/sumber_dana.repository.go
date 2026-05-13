package master

import (
	"sipropeda-backend/infras"

	"github.com/gofrs/uuid"
)

type SumberDanaRepository interface {
	Create(data SumberDana) error
	ResolveAll() ([]SumberDana, error)
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
	query := `SELECT id, nama_sumber, created_at, updated_at FROM sumber_dana WHERE is_deleted = false ORDER BY created_at DESC`
	err := r.db.Read.Select(&data, query)
	return data, err
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