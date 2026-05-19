package transaction

import (
	"sipropeda-backend/infras"

	"github.com/gofrs/uuid"
)

type PerankinganRepository interface {
	GetKriteriaAktif() ([]KriteriaTopsis, error)
	GetMatriksPenilaian(tahun int, tahap string) ([]MatriksPenilaian, error)
	SaveHasilPerankingan(data []ArsipPerankingan) error
	GetArsip(tahun int, tahap string) ([]ArsipPerankingan, error)
}

type perankinganRepository struct {
	db *infras.PostgresqlConn
}

func ProvidePerankinganRepository(db *infras.PostgresqlConn) PerankinganRepository {
	return &perankinganRepository{db: db}
}

func (r *perankinganRepository) GetKriteriaAktif() ([]KriteriaTopsis, error) {
	var data []KriteriaTopsis
	query := `SELECT id, kode, bobot, jenis FROM m_kriteria WHERE is_deleted = false ORDER BY kode ASC`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *perankinganRepository) GetMatriksPenilaian(tahun int, tahap string) ([]MatriksPenilaian, error) {
	var data []MatriksPenilaian
	// Tarik hanya usulan yang sesuai tahun, tahap, dan belum dihapus
	query := `
		SELECT p.usulan_id, p.kriteria_id, p.nilai_input
		FROM penilaian_usulan p
		JOIN usulan_proyek u ON p.usulan_id = u.id
		WHERE u.tahun_anggaran = $1 AND u.status_tahapan::text = $2::text AND u.is_deleted = false
	`
	err := r.db.Read.Select(&data, query, tahun, tahap)
	return data, err
}

func (r *perankinganRepository) SaveHasilPerankingan(data []ArsipPerankingan) error {
	if len(data) == 0 {
		return nil
	}

	// Hapus arsip lama untuk tahap ini agar tidak duplikat saat dihitung ulang
	deleteQuery := `
		DELETE FROM arsip_perankingan 
		WHERE usulan_id IN (
			SELECT id FROM usulan_proyek WHERE status_tahapan::text = $1::text
		) AND tahap_versi::text = $1::text
	`
	_, err := r.db.Write.Exec(deleteQuery, data[0].TahapVersi)
	if err != nil {
		return err
	}

	// Insert data perankingan baru
	insertQuery := `
		INSERT INTO arsip_perankingan (id, usulan_id, nilai_preferensi_v, ranking, tahap_versi, detail_kalkulasi)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, val := range data {
		newID, _ := uuid.NewV4()
		_, err = r.db.Write.Exec(insertQuery, newID, val.UsulanID, val.NilaiPreferensiV, val.Ranking, val.TahapVersi, val.DetailKalkulasi)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *perankinganRepository) GetArsip(tahun int, tahap string) ([]ArsipPerankingan, error) {
	var data []ArsipPerankingan
	
	// Tambahkan u.is_deleted = false agar usulan yang sudah dihapus (soft delete) tidak muncul lagi
	query := `
		SELECT a.id, a.usulan_id, u.nama_proyek as usulan_name, a.nilai_preferensi_v, a.ranking, a.tahap_versi, a.created_at
		FROM arsip_perankingan a
		JOIN usulan_proyek u ON a.usulan_id = u.id
		WHERE u.tahun_anggaran = $1 
		  AND a.tahap_versi::text = $2::text 
		  AND u.status_tahapan::text = $2::text 
		  AND u.is_deleted = false
		ORDER BY a.ranking ASC
	`
	err := r.db.Read.Select(&data, query, tahun, tahap)
	return data, err
}