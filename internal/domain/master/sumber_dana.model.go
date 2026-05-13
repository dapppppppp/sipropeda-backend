package master

import (
	"time"

	"github.com/gofrs/uuid"
)

// SumberDana merepresentasikan tabel sumber_dana di database
type SumberDana struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	NamaSumber string     `db:"nama_sumber" json:"namaSumber"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
}

// RequestSumberDana adalah format JSON yang dikirim dari Frontend
type RequestSumberDana struct {
	ID         uuid.UUID `json:"id" swaggerignore:"true"`
	NamaSumber string    `json:"namaSumber" validate:"required" example:"Dana Desa (DD)"`
}

// NewSumberDanaFormat memproses pembuatan atau update data Sumber Dana
func (s *SumberDana) NewSumberDanaFormat(reqFormat RequestSumberDana) (newSumberDana SumberDana) {
	now := time.Now()
	if reqFormat.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newSumberDana = SumberDana{
			ID:         newID,
			NamaSumber: reqFormat.NamaSumber,
			CreatedAt:  &now,
		}
	} else {
		newSumberDana = SumberDana{
			ID:         reqFormat.ID,
			NamaSumber: reqFormat.NamaSumber,
			UpdatedAt:  &now,
		}
	}
	return
}

// SoftDelete menandai data Sumber Dana sebagai terhapus
func (s *SumberDana) SoftDelete() {
	now := time.Now()
	s.IsDeleted = true
	s.UpdatedAt = &now
	s.DeletedAt = &now
}