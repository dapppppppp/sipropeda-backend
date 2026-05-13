package transaction

import (
	"github.com/gofrs/uuid"
)

// PenilaianUsulan merepresentasikan tabel penilaian_usulan di database
type PenilaianUsulan struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UsulanID     uuid.UUID `db:"usulan_id" json:"usulanId"`
	KriteriaID   uuid.UUID `db:"kriteria_id" json:"kriteriaId"`
	NilaiInput   float64   `db:"nilai_input" json:"nilaiInput"`
	// Hasil JOIN agar Frontend tahu ini nilai untuk kriteria apa
	KriteriaName *string   `db:"kriteria_name" json:"kriteriaName,omitempty"`
}

// RequestDetailPenilaian adalah format untuk 1 baris nilai
type RequestDetailPenilaian struct {
	KriteriaID uuid.UUID `json:"kriteriaId" validate:"required" example:"uuid-kriteria"`
	NilaiInput float64   `json:"nilaiInput" validate:"required" example:"85.5"`
}

// RequestBulkPenilaian adalah format JSON untuk menyimpan banyak nilai sekaligus untuk 1 usulan
type RequestBulkPenilaian struct {
	UsulanID uuid.UUID                `json:"usulanId" validate:"required" example:"uuid-usulan"`
	Data     []RequestDetailPenilaian `json:"data" validate:"required"`
}