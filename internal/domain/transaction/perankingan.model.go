package transaction

import (
	"time"

	"github.com/gofrs/uuid"
)

// ArsipPerankingan merepresentasikan tabel arsip_perankingan di database
type ArsipPerankingan struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	UsulanID         uuid.UUID  `db:"usulan_id" json:"usulanId"`
	UsulanName       *string    `db:"usulan_name" json:"usulanName,omitempty"` // Dari JOIN usulan_proyek
	NilaiPreferensiV float64    `db:"nilai_preferensi_v" json:"nilaiPreferensiV"`
	Ranking          int        `db:"ranking" json:"ranking"`
	TahapVersi       string     `db:"tahap_versi" json:"tahapVersi"`
	DetailKalkulasi  *string    `db:"detail_kalkulasi" json:"detailKalkulasi"` // Disimpan sebagai JSON string
	CreatedAt        *time.Time `db:"created_at" json:"createdAt"`
}

// RequestHitungTopsis adalah input dari Frontend untuk memicu kalkulasi
type RequestHitungTopsis struct {
	TahunAnggaran int    `json:"tahunAnggaran" validate:"required" example:"2025"`
	TahapVersi    string `json:"tahapVersi" validate:"required" example:"RKP"` // 'RKP' atau 'RAPBDes'
}

// --- Struct Pembantu untuk Algoritma (Tidak langsung ke DB) ---

type KriteriaTopsis struct {
	ID    uuid.UUID `db:"id"`
	Kode  *string   `db:"kode"`
	Bobot float64   `db:"bobot"`
	Jenis string    `db:"jenis"` // benefit / cost
}

type MatriksPenilaian struct {
	UsulanID   uuid.UUID `db:"usulan_id"`
	KriteriaID uuid.UUID `db:"kriteria_id"`
	NilaiInput float64   `db:"nilai_input"`
}