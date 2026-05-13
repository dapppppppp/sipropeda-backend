package auth

import (
	"time"
	"github.com/gofrs/uuid"
)

// Role merepresentasikan tabel roles di database
type Role struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
}

// RequestRole adalah format JSON yang dikirim dari Frontend
type RequestRole struct {
	ID          uuid.UUID `json:"id" swaggerignore:"true"`
	Name        string    `json:"name" validate:"required" example:"Admin Desa"`
	Description string    `json:"description" example:"Hak akses penuh untuk mengelola data master"`
}

func (r *Role) NewRoleFormat(req RequestRole) (newRole Role) {
	now := time.Now()
	if req.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newRole = Role{
			ID:          newID,
			Name:        req.Name,
			Description: req.Description,
			CreatedAt:   &now,
		}
	} else {
		newRole = Role{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
			UpdatedAt:   &now,
		}
	}
	return
}

func (r *Role) SoftDelete() {
	now := time.Now()
	r.IsDeleted = true
	r.UpdatedAt = &now
	r.DeletedAt = &now
}