package transaction

import (
	"time"

	"github.com/gofrs/uuid"
)

type PaguAnggaran struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	Tahun          int        `db:"tahun" json:"tahun"`
	SumberDanaID   uuid.UUID  `db:"sumber_dana_id" json:"sumberDanaId"`
	SumberDanaName *string    `db:"sumber_dana_name" json:"sumberDanaName,omitempty"` // Hasil JOIN
	JumlahPagu     float64    `db:"jumlah_pagu" json:"jumlahPagu"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy      *uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt      *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestPaguAnggaran struct {
	ID           uuid.UUID `json:"id" swaggerignore:"true"`
	Tahun        int       `json:"tahun" validate:"required" example:"2025"`
	SumberDanaID uuid.UUID `json:"sumberDanaId" validate:"required" example:"masukkan-uuid-sumber-dana-disini"`
	JumlahPagu   float64   `json:"jumlahPagu" validate:"required" example:"500000000"`
	UserID       uuid.UUID `json:"-"` // Dari JWT
}

func (p *PaguAnggaran) NewPaguAnggaranFormat(req RequestPaguAnggaran) (newData PaguAnggaran) {
	now := time.Now()
	if req.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newData = PaguAnggaran{
			ID:           newID,
			Tahun:        req.Tahun,
			SumberDanaID: req.SumberDanaID,
			JumlahPagu:   req.JumlahPagu,
			CreatedBy:    &req.UserID,
			CreatedAt:    &now,
		}
	} else {
		newData = PaguAnggaran{
			ID:           req.ID,
			Tahun:        req.Tahun,
			SumberDanaID: req.SumberDanaID,
			JumlahPagu:   req.JumlahPagu,
			UpdatedBy:    &req.UserID,
			UpdatedAt:    &now,
		}
	}
	return
}

func (p *PaguAnggaran) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	p.IsDeleted = true
	p.UpdatedAt = &now
	p.UpdatedBy = &userID
	p.DeletedAt = &now
}