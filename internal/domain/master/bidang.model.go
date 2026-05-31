package master

import (
	"time"

	"github.com/gofrs/uuid"
)

// BidangPembangunan merepresentasikan tabel bidang_pembangunan di database
type BidangPembangunan struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	NamaBidang string     `db:"nama_bidang" json:"namaBidang"`
	CreatedBy  *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedBy  *uuid.UUID `db:"updated_by" json:"updatedBy"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
}

// RequestBidangPembangunan adalah format JSON untuk Create dan Update
type RequestBidangPembangunan struct {
	ID         uuid.UUID `json:"id" swaggerignore:"true"`
	NamaBidang string    `json:"namaBidang" validate:"required" example:"Bidang Pemberdayaan Masyarakat"`
	UserID     uuid.UUID `json:"-"` // Dari JWT
}

func (b *BidangPembangunan) NewBidangPembangunanFormat(req RequestBidangPembangunan) (newData BidangPembangunan) {
	now := time.Now()
	if req.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newData = BidangPembangunan{
			ID:         newID,
			NamaBidang: req.NamaBidang,
			CreatedBy:  &req.UserID,
			CreatedAt:  &now,
		}
	} else {
		newData = BidangPembangunan{
			ID:         req.ID,
			NamaBidang: req.NamaBidang,
			UpdatedBy:  &req.UserID,
			UpdatedAt:  &now,
		}
	}
	return
}

func (b *BidangPembangunan) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &now
	b.UpdatedBy = &userID
	b.DeletedAt = &now
}