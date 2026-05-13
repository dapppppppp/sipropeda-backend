package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Password  string     `db:"password" json:"-"`
	RoleID    uuid.UUID  `db:"role_id" json:"roleId"` // <-- Pakai RoleID
	RoleName  *string    `db:"role_name" json:"roleName,omitempty"` // <-- Hasil JOIN
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestUserFormat struct {
	ID       uuid.UUID `json:"id" swaggerignore:"true"`
	Username string    `json:"username" validate:"required" example:"kepala_desa_1"`
	Password string    `json:"password" example:"rahasia123"` 
	RoleID   uuid.UUID `json:"roleId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"` // <-- Harus UUID dari tabel roles
}

func (u *User) NewUserFormat(reqFormat RequestUserFormat, hashedPassword string) (newUser User) {
	now := time.Now()
	if reqFormat.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newUser = User{
			ID:        newID,
			Username:  reqFormat.Username,
			Password:  hashedPassword,
			RoleID:    reqFormat.RoleID,
			CreatedAt: &now,
		}
	} else {
		newUser = User{
			ID:        reqFormat.ID,
			Username:  reqFormat.Username,
			Password:  hashedPassword,
			RoleID:    reqFormat.RoleID,
			UpdatedAt: &now,
		}
	}
	return
}

func (u *User) SoftDelete() {
	now := time.Now()
	u.IsDeleted = true
	u.UpdatedAt = &now
	u.DeletedAt = &now
}

type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"admin_desa"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	RoleID   string `json:"roleId"`
	RoleName string `json:"roleName"`
}