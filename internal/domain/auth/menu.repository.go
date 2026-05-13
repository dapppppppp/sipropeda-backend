package auth

import (
	"sipropeda-backend/infras"
	"github.com/gofrs/uuid"
)

type MenuRepository interface {
	ResolveAllMenu() ([]Menu, error)
	CreateMenu(data Menu) error
	UpdateMenu(data Menu) error
	DeleteMenu(id uuid.UUID) error
	
	// Fitur Menu Role
	ResolveMenuByRoleID(roleID uuid.UUID, level int, parentID *uuid.UUID) ([]MenuResponse, error)
	SaveBulkMenuRole(roleID uuid.UUID, data []MenuRole) error
}

type menuRepository struct {
	db *infras.PostgresqlConn
}

func ProvideMenuRepository(db *infras.PostgresqlConn) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) ResolveAllMenu() ([]Menu, error) {
	var data []Menu
	query := `SELECT id, name, link, icon, description, permission_label, action, level, seq, parent_id, created_at, updated_at FROM menus WHERE is_deleted = false ORDER BY level ASC, seq ASC`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *menuRepository) CreateMenu(data Menu) error {
	query := `INSERT INTO menus (id, name, link, icon, description, permission_label, action, level, seq, parent_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Write.Exec(query, data.ID, data.Name, data.Link, data.Icon, data.Description, data.PermissionLabel, data.Action, data.Level, data.Seq, data.ParentID)
	return err
}

func (r *menuRepository) UpdateMenu(data Menu) error {
	query := `UPDATE menus SET name=$1, link=$2, icon=$3, description=$4, permission_label=$5, action=$6, level=$7, seq=$8, parent_id=$9, updated_at=NOW() WHERE id=$10`
	_, err := r.db.Write.Exec(query, data.Name, data.Link, data.Icon, data.Description, data.PermissionLabel, data.Action, data.Level, data.Seq, data.ParentID, data.ID)
	return err
}

func (r *menuRepository) DeleteMenu(id uuid.UUID) error {
	query := `UPDATE menus SET is_deleted = true, deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Write.Exec(query, id)
	return err
}

// Ambil hierarki menu berdasarkan Role
func (r *menuRepository) ResolveMenuByRoleID(roleID uuid.UUID, level int, parentID *uuid.UUID) ([]MenuResponse, error) {
	var data []MenuResponse
	query := `
		SELECT mr.id, mr.menu_id, mr.role_id, m.name, m.link, m.icon, m.level, m.seq, m.permission_label, mr.permission 
		FROM menu_roles mr 
		JOIN menus m ON mr.menu_id = m.id 
		WHERE mr.role_id = $1 AND m.is_deleted = false AND m.level = $2
	`
	
	var err error
	if parentID == nil {
		query += ` AND m.parent_id IS NULL ORDER BY m.seq ASC`
		err = r.db.Read.Select(&data, query, roleID, level)
	} else {
		query += ` AND m.parent_id = $3 ORDER BY m.seq ASC`
		err = r.db.Read.Select(&data, query, roleID, level, parentID)
	}

	return data, err
}

// Timpa semua akses menu untuk role tertentu
func (r *menuRepository) SaveBulkMenuRole(roleID uuid.UUID, data []MenuRole) error {
	tx, err := r.db.Write.Beginx()
	if err != nil { return err }

	_, err = tx.Exec(`DELETE FROM menu_roles WHERE role_id = $1`, roleID)
	if err != nil { tx.Rollback(); return err }

	insertQuery := `INSERT INTO menu_roles (id, menu_id, role_id, permission) VALUES ($1, $2, $3, $4)`
	for _, val := range data {
		_, err = tx.Exec(insertQuery, val.ID, val.MenuID, val.RoleID, val.Permission)
		if err != nil { tx.Rollback(); return err }
	}

	return tx.Commit()
}