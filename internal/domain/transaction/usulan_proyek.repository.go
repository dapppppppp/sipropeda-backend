package transaction

import (
	"sipropeda-backend/infras"

	"github.com/gofrs/uuid"
)

type UsulanProyekRepository interface {
	Create(data UsulanProyek) error
	ResolveAll() ([]UsulanProyek, error)
	ResolveByID(id uuid.UUID) (UsulanProyek, error)
	Update(data UsulanProyek) error
	Delete(data UsulanProyek) error
}

type usulanProyekRepository struct {
	db *infras.PostgresqlConn
}

func ProvideUsulanProyekRepository(db *infras.PostgresqlConn) UsulanProyekRepository {
	return &usulanProyekRepository{db: db}
}

func (r *usulanProyekRepository) Create(data UsulanProyek) error {
	query := `
		INSERT INTO usulan_proyek 
		(id, tahun_anggaran, nama_proyek, lokasi, volume, satuan, nilai_rab, status_sifat, status_tahapan, sumber_dana_id, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Write.Exec(query, data.ID, data.TahunAnggaran, data.NamaProyek, data.Lokasi, data.Volume, data.Satuan, data.NilaiRAB, data.StatusSifat, data.StatusTahapan, data.SumberDanaID, data.CreatedBy, data.CreatedAt)
	return err
}

func (r *usulanProyekRepository) ResolveAll() ([]UsulanProyek, error) {
	var data []UsulanProyek
	query := `
		SELECT u.id, u.tahun_anggaran, u.nama_proyek, u.lokasi, u.volume, u.satuan, u.nilai_rab, u.status_sifat, u.status_tahapan, u.sumber_dana_id, s.nama_sumber as sumber_dana_name, u.created_at, u.updated_at 
		FROM usulan_proyek u
		LEFT JOIN sumber_dana s ON u.sumber_dana_id = s.id
		WHERE u.is_deleted = false 
		ORDER BY u.created_at DESC
	`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *usulanProyekRepository) ResolveByID(id uuid.UUID) (UsulanProyek, error) {
	var data UsulanProyek
	query := `
		SELECT u.id, u.tahun_anggaran, u.nama_proyek, u.lokasi, u.volume, u.satuan, u.nilai_rab, u.status_sifat, u.status_tahapan, u.sumber_dana_id, s.nama_sumber as sumber_dana_name, u.created_at, u.updated_at 
		FROM usulan_proyek u
		LEFT JOIN sumber_dana s ON u.sumber_dana_id = s.id
		WHERE u.id = $1 AND u.is_deleted = false
	`
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *usulanProyekRepository) Update(data UsulanProyek) error {
	query := `
		UPDATE usulan_proyek 
		SET tahun_anggaran = $1, 
		    nama_proyek = $2, 
		    lokasi = $3, 
		    volume = $4, 
		    satuan = $5, 
		    nilai_rab = $6, 
		    status_sifat = $7, 
		    sumber_dana_id = $8, 
		    status_tahapan = $9, 
		    updated_by = $10, 
		    updated_at = $11 
		WHERE id = $12 AND is_deleted = false
	`
	_, err := r.db.Write.Exec(query, 
		data.TahunAnggaran, // $1
		data.NamaProyek,    // $2
		data.Lokasi,        // $3
		data.Volume,        // $4
		data.Satuan,        // $5
		data.NilaiRAB,      // $6
		data.StatusSifat,   // $7
		data.SumberDanaID,  // $8
		data.StatusTahapan, // $9 (Ini tambahan krusialnya!)
		data.UpdatedBy,     // $10
		data.UpdatedAt,     // $11
		data.ID,            // $12
	)
	return err
}

func (r *usulanProyekRepository) Delete(data UsulanProyek) error {
	query := `
		UPDATE usulan_proyek 
		SET is_deleted = $1, deleted_at = $2, updated_by = $3, updated_at = $4 
		WHERE id = $5
	`
	_, err := r.db.Write.Exec(query, data.IsDeleted, data.DeletedAt, data.UpdatedBy, data.UpdatedAt, data.ID)
	return err
}