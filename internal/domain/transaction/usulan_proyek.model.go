package transaction

import (
	"time"

	"github.com/gofrs/uuid"
)

// UsulanProyek merepresentasikan tabel usulan_proyek di database
type UsulanProyek struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	TahunAnggaran  int        `db:"tahun_anggaran" json:"tahunAnggaran"`
	NamaProyek     string     `db:"nama_proyek" json:"namaProyek"`
	Lokasi         string     `db:"lokasi" json:"lokasi"`
	Volume         float64    `db:"volume" json:"volume"`
	Satuan         string     `db:"satuan" json:"satuan"`
	NilaiRAB       float64    `db:"nilai_rab" json:"nilaiRab"`
	StatusSifat    string     `db:"status_sifat" json:"statusSifat"`       // 'Reguler' atau 'Mandatori'
	StatusTahapan  string     `db:"status_tahapan" json:"statusTahapan"`   // 'draft_rkp', dll
	SumberDanaID   *uuid.UUID `db:"sumber_dana_id" json:"sumberDanaId"`    // Pointer karena bisa null
	SumberDanaName *string    `db:"sumber_dana_name" json:"sumberDanaName,omitempty"` // Hasil JOIN
	ApprovedBy     *uuid.UUID `db:"approved_by" json:"approvedBy"`
	ApprovedAt     *time.Time `db:"approved_at" json:"approvedAt"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy      *uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt      *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
}

// RequestUsulanProyek adalah format JSON untuk Create dan Update
type RequestUsulanProyek struct {
	ID            uuid.UUID  `json:"id" swaggerignore:"true"`
	TahunAnggaran int        `json:"tahunAnggaran" validate:"required" example:"2025"`
	NamaProyek    string     `json:"namaProyek" validate:"required" example:"Pembangunan Gorong-Gorong"`
	Lokasi        string     `json:"lokasi" validate:"required" example:"Dusun Sukamaju RT 01"`
	Volume        float64    `json:"volume" example:"150.5"`
	Satuan        string     `json:"satuan" example:"Meter"`
	NilaiRAB      float64    `json:"nilaiRab" validate:"required" example:"45000000"`
	StatusSifat   string     `json:"statusSifat" validate:"required" example:"Reguler"`
	StatusTahapan string     `json:"statusTahapan"` // <-- TAMBAHKAN BARIS INI
	SumberDanaID  *uuid.UUID `json:"sumberDanaId" example:"masukkan-uuid-sumber-dana-disini"`
	UserID        uuid.UUID  `json:"-"` // Dari JWT
}

func (u *UsulanProyek) NewUsulanProyekFormat(req RequestUsulanProyek) (newData UsulanProyek) {
	now := time.Now()
	if req.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newData = UsulanProyek{
			ID:            newID,
			TahunAnggaran: req.TahunAnggaran,
			NamaProyek:    req.NamaProyek,
			Lokasi:        req.Lokasi,
			Volume:        req.Volume,
			Satuan:        req.Satuan,
			NilaiRAB:      req.NilaiRAB,
			StatusSifat:   req.StatusSifat,
			StatusTahapan: "RKP", // Default awal
			SumberDanaID:  req.SumberDanaID,
			CreatedBy:     &req.UserID,
			CreatedAt:     &now,
		}
	} else {
		newData = UsulanProyek{
			ID:            req.ID,
			TahunAnggaran: req.TahunAnggaran,
			NamaProyek:    req.NamaProyek,
			Lokasi:        req.Lokasi,
			Volume:        req.Volume,
			Satuan:        req.Satuan,
			NilaiRAB:      req.NilaiRAB,
			StatusSifat:   req.StatusSifat,
			StatusTahapan: req.StatusTahapan, // <-- TAMBAHKAN BARIS INI
			SumberDanaID:  req.SumberDanaID,
			UpdatedBy:     &req.UserID,
			UpdatedAt:     &now,
		}
	}
	return
}

func (u *UsulanProyek) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	u.IsDeleted = true
	u.UpdatedAt = &now
	u.UpdatedBy = &userID
	u.DeletedAt = &now
}