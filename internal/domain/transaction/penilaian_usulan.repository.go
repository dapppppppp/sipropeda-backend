package transaction

import (
	"sipropeda-backend/infras"

	"github.com/gofrs/uuid"
)

type PenilaianUsulanRepository interface {
	SaveBulk(usulanID uuid.UUID, data []PenilaianUsulan) error
	GetByUsulanID(usulanID uuid.UUID) ([]PenilaianUsulan, error)
}

type penilaianUsulanRepository struct {
	db *infras.PostgresqlConn
}

func ProvidePenilaianUsulanRepository(db *infras.PostgresqlConn) PenilaianUsulanRepository {
	return &penilaianUsulanRepository{db: db}
}

func (r *penilaianUsulanRepository) SaveBulk(usulanID uuid.UUID, data []PenilaianUsulan) error {
	// 1. Hapus semua data lama untuk usulan ini (Hard Delete, karena ada ON DELETE CASCADE)
	deleteQuery := `DELETE FROM penilaian_usulan WHERE usulan_id = $1`
	_, err := r.db.Write.Exec(deleteQuery, usulanID)
	if err != nil {
		return err
	}

	// 2. Insert data baru satu per satu
	insertQuery := `INSERT INTO penilaian_usulan (usulan_id, kriteria_id, nilai_input) VALUES ($1, $2, $3)`
	for _, val := range data {
		// Abaikan jika usulan_id kosong untuk mencegah error
		if val.UsulanID != uuid.Nil && val.KriteriaID != uuid.Nil {
			_, err = r.db.Write.Exec(insertQuery, val.UsulanID, val.KriteriaID, val.NilaiInput)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *penilaianUsulanRepository) GetByUsulanID(usulanID uuid.UUID) ([]PenilaianUsulan, error) {
	var data []PenilaianUsulan
	// Join dengan tabel m_kriteria untuk mengambil nama kriteria
	query := `
		SELECT p.id, p.usulan_id, p.kriteria_id, k.nama as kriteria_name, p.nilai_input
		FROM penilaian_usulan p
		LEFT JOIN m_kriteria k ON p.kriteria_id = k.id
		WHERE p.usulan_id = $1
		ORDER BY k.kode ASC
	`
	err := r.db.Read.Select(&data, query, usulanID)
	return data, err
}