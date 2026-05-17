package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Name      string     `db:"nama" json:"name"`     // <-- BACA DB: 'nama', KIRIM JSON: 'name'
	Email     string     `db:"email" json:"email"`
	Foto      *string    `db:"foto" json:"foto"`
	Username  string     `db:"username" json:"username"`
	Password  string     `db:"password" json:"-"`
	RoleID    uuid.UUID  `db:"role_id" json:"roleId"`
	RoleName  *string    `db:"role_name" json:"roleName,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestUserFormat struct {
	ID       uuid.UUID `json:"id" swaggerignore:"true"`
	Name     string    `json:"name" validate:"required"` // <-- MENERIMA JSON: 'name' dari Frontend
	Email    string    `json:"email" validate:"required,email"`
	Username string    `json:"username" validate:"required"`
	Password string    `json:"password"` 
	RoleID   uuid.UUID `json:"roleId" validate:"required"`
}

type ResetPasswordRequest struct {
	ID          string `json:"id" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (u *User) NewUserFormat(reqFormat RequestUserFormat, hashedPassword string) (newUser User) {
	now := time.Now()
	if reqFormat.ID == uuid.Nil {
		newID, _ := uuid.NewV4()
		newUser = User{
			ID:        newID,
			Name:      reqFormat.Name, 
			Email:     reqFormat.Email,
			Username:  reqFormat.Username,
			Password:  hashedPassword,
			RoleID:    reqFormat.RoleID,
			CreatedAt: &now,
		}
	} else {
		newUser = User{
			ID:        reqFormat.ID,
			Name:      reqFormat.Name,
			Email:     reqFormat.Email,
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
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	RoleID   string `json:"roleId"`
	RoleName string `json:"roleName"`
	User     User   `json:"user"`
}