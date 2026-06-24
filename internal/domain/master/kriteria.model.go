package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type Kriteria struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Kode      *string    `db:"kode" json:"kode"`
	Nama      string     `db:"nama" json:"nama"`
	Bobot     float64    `db:"bobot" json:"bobot"`
	Jenis     string     `db:"jenis" json:"jenis"`
	IsActive  bool       `db:"is_active" json:"isActive"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestKriteriaFormat struct {
	ID     uuid.UUID `json:"id"`
	Kode   string    `json:"kode" validate:"required" example:"C1"`
	Nama   string    `json:"nama" validate:"required"`
	Bobot  float64   `json:"bobot" validate:"required"`
	Jenis  string    `json:"jenis" validate:"required,oneof=benefit cost"`
	UserID uuid.UUID `json:"-"`
}

var ColumnMapKriteria = map[string]interface{}{
	"id":        "id",
	"kode":      "kode",
	"nama":      "nama",
	"bobot":     "bobot",
	"jenis":     "jenis",
	"isActive":  "is_active",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}

func (k *Kriteria) NewKriteriaFormat(reqFormat RequestKriteriaFormat) (newKriteria Kriteria, err error) {
	now := time.Now()
	if reqFormat.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newKriteria = Kriteria{
			ID:        newID,
			Kode:      &reqFormat.Kode,
			Nama:      reqFormat.Nama,
			Bobot:     reqFormat.Bobot,
			Jenis:     reqFormat.Jenis,
			IsActive:  true,
			CreatedAt: &now,
			CreatedBy: &reqFormat.UserID,
		}
	} else {
		newKriteria = Kriteria{
			ID:        reqFormat.ID,
			Kode:      &reqFormat.Kode,
			Nama:      reqFormat.Nama,
			Bobot:     reqFormat.Bobot,
			Jenis:     reqFormat.Jenis,
			UpdatedAt: &now,
			UpdatedBy: &reqFormat.UserID,
		}
	}
	return
}

func (k *Kriteria) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	k.IsDeleted = true
	k.UpdatedAt = &now
	k.UpdatedBy = &userID
	k.DeletedAt = &now
}