package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type Menu struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	Link            string     `db:"link" json:"link"`
	Icon            *string    `db:"icon" json:"icon"`
	Description     *string    `db:"description" json:"description"`
	PermissionLabel *string    `db:"permission_label" json:"permissionLabel"`
	Action          *string    `db:"action" json:"action"`
	Level           int        `db:"level" json:"level"`
	Seq             int        `db:"seq" json:"seq"`
	ParentID        *uuid.UUID `db:"parent_id" json:"parentId"`
	CreatedAt       *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	IsDeleted       bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestMenuFormat struct {
	ID              uuid.UUID  `json:"id" swaggerignore:"true"`
	Name            string     `json:"name" validate:"required"`
	Link            string     `json:"link" validate:"required"`
	Icon            *string    `json:"icon"`
	Description     *string    `json:"description"`
	PermissionLabel *string    `json:"permissionLabel"`
	Action          *string    `json:"action"`
	Level           int        `json:"level"`
	Seq             int        `json:"seq"`
	ParentID        *uuid.UUID `json:"parentId"`
}

type MenuRole struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	MenuID     uuid.UUID  `db:"menu_id" json:"menuId"`
	RoleID     uuid.UUID  `db:"role_id" json:"roleId"`
	Permission *string    `db:"permission" json:"permission"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
}

// MenuResponse untuk output berjenjang (Tree) di Frontend
type MenuResponse struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	MenuID          uuid.UUID      `db:"menu_id" json:"menuId"`
	RoleID          uuid.UUID      `db:"role_id" json:"roleId"`
	Name            string         `db:"name" json:"name"`
	Link            string         `db:"link" json:"link"`
	Icon            *string        `db:"icon" json:"icon"`
	Level           int            `db:"level" json:"level"`
	Seq             int            `db:"seq" json:"seq"`
	PermissionLabel *string        `db:"permission_label" json:"permissionLabel"`
	Permission      *string        `db:"permission" json:"permission"`
	PermissionList  []string       `json:"permissionList"`
	Children        []MenuResponse `json:"children"`
}

type RequestBulkMenuRole struct {
	RoleID string                `json:"roleId" validate:"required"`
	Data   []RequestMenuRoleItem `json:"data"`
}

type RequestMenuRoleItem struct {
	MenuID     string  `json:"menuId"`
	Permission *string `json:"permission"`
}