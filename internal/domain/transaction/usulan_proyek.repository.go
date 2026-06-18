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
	BulkInsert(data []UsulanProyek) error
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
        (id, tahun_anggaran, bidang_id, nama_proyek, lokasi, volume, satuan, nilai_rab, status_sifat, status_tahapan, sumber_dana_id, created_by, created_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `
	_, err := r.db.Write.Exec(query, data.ID, data.TahunAnggaran, data.BidangID, data.NamaProyek, data.Lokasi, data.Volume, data.Satuan, data.NilaiRAB, data.StatusSifat, data.StatusTahapan, data.SumberDanaID, data.CreatedBy, data.CreatedAt)
	return err
}

func (r *usulanProyekRepository) ResolveAll() ([]UsulanProyek, error) {
	var data []UsulanProyek
	query := `
        SELECT 
            u.id, u.tahun_anggaran, u.bidang_id, b.nama_bidang as bidang_name, 
            u.nama_proyek, u.lokasi, 
            COALESCE(u.volume, 0) as volume, 
            COALESCE(u.satuan, '') as satuan, 
            COALESCE(u.nilai_rab, 0) as nilai_rab, 
            u.status_sifat, u.status_tahapan, u.sumber_dana_id, s.nama_sumber as sumber_dana_name, 
            u.created_at, u.updated_at,
            COALESCE(p.nilai_preferensi_v, 0) as nilai_preferensi_v,
            EXISTS (SELECT 1 FROM penilaian_usulan pu WHERE pu.usulan_id = u.id) as sudah_dinilai
        FROM usulan_proyek u
        LEFT JOIN sumber_dana s ON u.sumber_dana_id = s.id
        LEFT JOIN bidang_pembangunan b ON u.bidang_id = b.id
        LEFT JOIN arsip_perankingan p ON u.id = p.usulan_id
        WHERE u.is_deleted = false 
        ORDER BY nilai_preferensi_v DESC, u.created_at DESC
    `
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *usulanProyekRepository) ResolveByID(id uuid.UUID) (UsulanProyek, error) {
	var data UsulanProyek
	query := `
        SELECT 
            u.id, u.tahun_anggaran, u.bidang_id, b.nama_bidang as bidang_name, 
            u.nama_proyek, u.lokasi, 
            COALESCE(u.volume, 0) as volume, 
            COALESCE(u.satuan, '') as satuan, 
            COALESCE(u.nilai_rab, 0) as nilai_rab, 
            u.status_sifat, u.status_tahapan, u.sumber_dana_id, s.nama_sumber as sumber_dana_name, 
            u.created_at, u.updated_at,
            COALESCE(p.nilai_preferensi_v, 0) as nilai_preferensi_v,
            EXISTS (SELECT 1 FROM penilaian_usulan pu WHERE pu.usulan_id = u.id) as sudah_dinilai
        FROM usulan_proyek u
        LEFT JOIN sumber_dana s ON u.sumber_dana_id = s.id
        LEFT JOIN bidang_pembangunan b ON u.bidang_id = b.id
        LEFT JOIN arsip_perankingan p ON u.id = p.usulan_id
        WHERE u.id = $1 AND u.is_deleted = false
    `
	err := r.db.Read.Get(&data, query, id)
	return data, err
}

func (r *usulanProyekRepository) Update(data UsulanProyek) error {
	query := `
        UPDATE usulan_proyek 
        SET tahun_anggaran = $1, 
            bidang_id = $2,
            nama_proyek = $3, 
            lokasi = $4, 
            volume = $5, 
            satuan = $6, 
            nilai_rab = $7, 
            status_sifat = $8, 
            sumber_dana_id = $9, 
            status_tahapan = $10, 
            updated_by = $11, 
            updated_at = $12 
        WHERE id = $13 AND is_deleted = false
    `
	_, err := r.db.Write.Exec(query, 
		data.TahunAnggaran,
		data.BidangID,
		data.NamaProyek,
		data.Lokasi,
		data.Volume,
		data.Satuan,
		data.NilaiRAB,
		data.StatusSifat,
		data.SumberDanaID,
		data.StatusTahapan,
		data.UpdatedBy,
		data.UpdatedAt,
		data.ID,
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

// =================================================================================
// BULK INSERT DENGAN AUTO-MAPPING MASTER DATA (SQL SUBQUERY)
// =================================================================================
func (r *usulanProyekRepository) BulkInsert(data []UsulanProyek) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := r.db.Write.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Query ini sangat cerdas. Ia menggunakan teks mentah dari Excel ($6 dan $7) 
	// untuk mencari UUID di tabel Master secara otomatis (Case Insensitive).
	query := `
		INSERT INTO usulan_proyek 
		(id, tahun_anggaran, nama_proyek, lokasi, nilai_rab, status_tahapan, status_sifat, volume, satuan, is_deleted, bidang_id, sumber_dana_id) 
		VALUES ($1, $2, $3, $4, $5, 'RKP', 'Reguler', 0, '', false,
			(CASE WHEN $6 = '' THEN NULL ELSE (SELECT id FROM bidang_pembangunan WHERE nama_bidang ILIKE '%' || $6 || '%' OR $6 ILIKE '%' || nama_bidang || '%' LIMIT 1) END),
			(CASE WHEN $7 = '' THEN NULL ELSE (SELECT id FROM sumber_dana WHERE nama_sumber ILIKE '%' || $7 || '%' OR $7 ILIKE '%' || nama_sumber || '%' LIMIT 1) END)
		)
	`
	for _, val := range data {
		// val.StatusTahapan = Teks Bidang Mentah (contoh: "Bidang Penyelenggaraan Pemerintahan Desa")
		// val.StatusSifat = Teks Sumber Dana Mentah (contoh: "DD/ADD")
		_, err := tx.Exec(query, val.ID, val.TahunAnggaran, val.NamaProyek, val.Lokasi, val.NilaiRAB, val.StatusTahapan, val.StatusSifat)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}