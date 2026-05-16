package auth

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"sipropeda-backend/infras"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	menuQuery = struct {
		Select, SelectDTO, Insert, Update, Count, Urutan string
	}{
		Select: `SELECT id, name, link, description, icon, permission_label, action, level, seq, parent_id, created_at, updated_at, created_by, updated_by, is_deleted from menus `,
		SelectDTO: `SELECT m.id, m.name, m.link, m.description, m.icon, m.permission_label, m.action, m.level, m.seq, m.parent_id, 
            m.created_at, m.updated_at, m.created_by, m.updated_by, m.is_deleted, cm.name parent_menu from menus m
            LEFT JOIN menus cm ON m.parent_id = cm.id `,
		Insert: `Insert into menus(id, name, link, description, icon, permission_label, action, level, seq, parent_id, created_at, created_by) 
            values (:id, :name, :link, :description, :icon, :permission_label, :action, :level, :seq, :parent_id, :created_at, :created_by)`,
		Update: `Update menus set 
            id=:id, name=:name, link=:link, description=:description, icon=:icon,
            permission_label=:permission_label, action=:action, level=:level, seq=:seq,
            parent_id=:parent_id, updated_at=:updated_at, updated_by=:updated_by, is_deleted=:is_deleted `,
		Count:  `Select count(id) from menus`,
		Urutan: `select coalesce(max(seq),0) + 1 as urutan from menus mu `,
	}

	menuRoleQuery = struct {
		Select, SelectDTO, SelectDTOTrx, Insert, Update, UpdatePermission, InsertBulk, InsertBulkPlaceholder string
	}{
		Select: `SELECT id, menu_id, role_id, permission, commodity_id, created_at from menu_roles `,
		SelectDTO: `SELECT mr.id, mr.menu_id, m.name, m.link, m.description, m.icon, m.level, m.seq, m.permission_label, mr.permission 
            FROM menu_roles mr JOIN menus m on mr.menu_id = m.id `,
		SelectDTOTrx: ` SELECT mr.id, m.id menu_id, m.name, m.link, m.description, m.icon, m.level, m.seq, 
            m.permission_label, m.action, mr.permission FROM menus m
            LEFT JOIN (
                SELECT mr.id, mr.menu_id, mr.permission from menu_roles mr
                WHERE mr.role_id = $1
                AND (case when $2 = '' then mr.commodity_id IS NULL else mr.commodity_id::varchar=$2 end)
            )mr on mr.menu_id = m.id `,
		Insert: `Insert into menu_roles(id, menu_id, role_id, permission, commodity_id, created_at) values
                     (:id, :menu_id, :role_id, :permission, :commodity_id, :created_at)`,
		Update:                `Update menu_roles set id=:id, menu_id=:menu_id, role_id=:role_id, commodity_id=:commodity_id `,
		UpdatePermission:      `Update menu_roles set id=:id, permission=:permission `,
		InsertBulk:            `INSERT INTO public.menu_roles(id, menu_id, role_id, permission, commodity_id, created_at) VALUES `,
		InsertBulkPlaceholder: ` (:id, :menu_id, :role_id, :permission, :commodity_id, :created_at) `,
	}
)

type MenuRepository interface {
	GetAllMenu() (dataMenu []Menu, err error)
	ResolveAll(req model.StandardRequest) (dataMenu pagination.Response, err error)
	ResolveMenuByRoleID(req MenuRequest) (data []MenuResponse, err error)
	ResolveMenuByParentID(req MenuRequest) (data []MenuResponse, err error)
	ResolveMenuByRoleIDTrx(req MenuRequest) (data []MenuResponseTrx, err error)
	ResolveMenuByParentIDTrx(req MenuRequest) (data []MenuResponseTrx, err error)
	CreateMenu(menu Menu) error
	UpdateMenu(menu Menu) error
	ResolveMenuByID(id uuid.UUID) (menu Menu, err error)
	ResolveMenuRoleByID(id uuid.UUID) (MenuRole MenuRole, err error)
	UpdateMenuRole(MenuRole MenuRole) error
	UpdatePermission(MenuRole MenuRole) error
	CreateBulkMenuRole(req []MenuRole) error
}

type MenuRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideMenuRepositoryPostgreSQL(db *infras.PostgresqlConn) *MenuRepositoryPostgreSQL {
	s := new(MenuRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuByRoleID(req MenuRequest) (data []MenuResponse, err error) {
	// PERBAIKAN: Inisialisasi slice agar tidak mereturn null ketika kosong
	data = make([]MenuResponse, 0)

	criteria := ` WHERE m.level = 1 AND coalesce(m.is_deleted,false)=false AND mr.role_id::varchar=$1 
        AND mr.permission ilike '%VIEW%' `
	if req.CommodityId != "" {
		criteria += fmt.Sprintf(" AND mr.commodity_id='%v' ", req.CommodityId)
	} else {
		criteria += " AND mr.commodity_id IS NULL "
	}
	criteria += " ORDER BY m.seq ASC"

	err = r.DB.Read.Select(&data, menuRoleQuery.SelectDTO+criteria, req.RoleId)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuByParentID(req MenuRequest) (data []MenuResponse, err error) {
	// PERBAIKAN: Inisialisasi slice agar tidak mereturn null ketika kosong
	data = make([]MenuResponse, 0)

	criteria := ` WHERE coalesce(m.is_deleted,false)=false AND mr.role_id::varchar=$1 AND m.parent_id::varchar=$2 
        AND mr.permission ilike '%VIEW%' `
	if req.CommodityId != "" {
		criteria += fmt.Sprintf(" AND mr.commodity_id='%v' ", req.CommodityId)
	} else {
		criteria += " AND mr.commodity_id IS NULL "
	}
	criteria += " ORDER BY m.seq ASC"

	err = r.DB.Read.Select(&data, menuRoleQuery.SelectDTO+criteria, req.RoleId, req.ParentId)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuByRoleIDTrx(req MenuRequest) (data []MenuResponseTrx, err error) {
	// PERBAIKAN: Inisialisasi slice agar tidak mereturn null ketika kosong
	data = make([]MenuResponseTrx, 0)

	criteria := ` WHERE m.level = 1 AND coalesce(m.is_deleted,false)=false ORDER BY m.seq ASC `
	err = r.DB.Read.Select(&data, menuRoleQuery.SelectDTOTrx+criteria, req.RoleId, req.CommodityId)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuByParentIDTrx(req MenuRequest) (data []MenuResponseTrx, err error) {
	// PERBAIKAN: Inisialisasi slice agar tidak mereturn null ketika kosong
	data = make([]MenuResponseTrx, 0)

	criteria := ` WHERE coalesce(m.is_deleted,false) = false AND m.parent_id = $3 ORDER BY m.seq ASC `
	err = r.DB.Read.Select(&data, menuRoleQuery.SelectDTOTrx+criteria, req.RoleId, req.CommodityId, req.ParentId)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) GetAllMenu() (dataMenu []Menu, err error) {
	// PERBAIKAN: Inisialisasi slice agar tidak mereturn null ketika kosong
	dataMenu = make([]Menu, 0)

	criteria := ` WHERE COALESCE(is_deleted, false) = false ORDER BY seq ASC `
	err = r.DB.Read.Select(&dataMenu, menuQuery.Select+criteria)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (dataMenu pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" where coalesce(m.is_deleted, false) = ? ")
	searchParams = append(searchParams, false)

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND concat(m.name, m.link, m.description, cm.name) ilike ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind("select count(*) from(" + menuQuery.SelectDTO + searchRoleBuff.String() + ")x")
	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		log.Println(err)
		return
	}

	if totalData < 1 {
		dataMenu.Items = make([]interface{}, 0)
		return
	}

	searchRoleBuff.WriteString("order by " + ColumnMappMenu[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchMenuQuery := searchRoleBuff.String()
	searchMenuQuery = r.DB.Read.Rebind(menuQuery.SelectDTO + searchMenuQuery)
	rows, err := r.DB.Read.Queryx(searchMenuQuery, searchParams...)
	if err != nil {
		return
	}
	
	var items []interface{}
	for rows.Next() {
		var menu MenuDTO
		err = rows.StructScan(&menu)
		if err != nil {
			return
		}
		items = append(items, menu)
	}

	dataMenu.Items = items
	dataMenu.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *MenuRepositoryPostgreSQL) CreateMenu(menu Menu) error {
	stmt, err := r.DB.Read.PrepareNamed(menuQuery.Insert)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(menu)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuByID(id uuid.UUID) (menu Menu, err error) {
	err = r.DB.Read.Get(&menu, menuQuery.Select+" WHERE id=$1 AND coalesce(is_deleted, false) = false ", id)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) UpdateMenu(menu Menu) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateMenu(tx, menu); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateMenu(tx *sqlx.Tx, menu Menu) (err error) {
	stmt, err := tx.PrepareNamed(menuQuery.Update + " WHERE id=:id")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(menu)
	if err != nil {
		log.Println(err)
	}
	return
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuRoleByID(id uuid.UUID) (MenuRole MenuRole, err error) {
	err = r.DB.Read.Get(&MenuRole, menuRoleQuery.Select+" WHERE id=$1 AND coalesce(is_deleted, false) = false ", id)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) UpdateMenuRole(MenuRole MenuRole) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateMenuRole(tx, MenuRole); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateMenuRole(tx *sqlx.Tx, MenuRole MenuRole) (err error) {
	stmt, err := tx.PrepareNamed(menuRoleQuery.Update + " WHERE id=:id")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(MenuRole)
	if err != nil {
		log.Println(err)
	}
	return
}

func (r *MenuRepositoryPostgreSQL) UpdatePermission(MenuRole MenuRole) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdatePermission(tx, MenuRole); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdatePermission(tx *sqlx.Tx, MenuRole MenuRole) (err error) {
	stmt, err := tx.PrepareNamed(menuRoleQuery.UpdatePermission + " WHERE id=:id")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(MenuRole)
	if err != nil {
		log.Println(err)
	}
	return
}

func (u *MenuRepositoryPostgreSQL) CreateBulkMenuRole(req []MenuRole) error {
	return u.DB.WithTransaction(func(db *sqlx.Tx, e chan error) {
		ids := make([]string, 0)
		for _, d := range req {
			ids = append(ids, d.ID)
		}

		var roleId string
		var commodityId *int
		if len(req) > 0 {
			roleId = req[0].RoleId
			commodityId = req[0].CommodityId
		}

		if err := u.txDeleteDetailNotIn(db, roleId, commodityId, ids); err != nil {
			e <- err
			return
		}

		if err := txCreateMenuRole(db, req); err != nil {
			log.Println("err create menu role : ", err)
			e <- err
			return
		}

		e <- nil
	})
}

func txCreateMenuRole(tx *sqlx.Tx, details []MenuRole) (err error) {
	if len(details) == 0 {
		return
	}
	query, args, err := ComposeBulkUpsertMenuRoleQuery(details)
	if err != nil {
		return
	}

	query = tx.Rebind(query)
	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func ComposeBulkUpsertMenuRoleQuery(details []MenuRole) (qResult string, params []interface{}, err error) {
	values := []string{}
	for _, d := range details {
		param := map[string]interface{}{
			"id":           d.ID,
			"menu_id":      d.MenuId,
			"role_id":      d.RoleId,
			"permission":   d.Permission,
			"commodity_id": d.CommodityId,
			"created_at":   d.CreatedAt,
		}
		q, args, err := sqlx.Named(menuRoleQuery.InsertBulkPlaceholder, param)
		if err != nil {
			return qResult, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	qResult = fmt.Sprintf(`%v %v 
                        ON CONFLICT (id) 
                        DO UPDATE SET permission=EXCLUDED.permission `, menuRoleQuery.InsertBulk, strings.Join(values, ","))
	return
}

func (r *MenuRepositoryPostgreSQL) UrutanMenuByID(req UrutanRequest) (urutan Urutan, err error) {
	err = r.DB.Read.Get(&urutan, menuQuery.Urutan+" WHERE m.is_deleted=false and m.level=$1 and m.parent_id=$2", req.Level, req.IdParent)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) txDeleteDetailNotIn(tx *sqlx.Tx, roleId string, commodityId *int, ids []string) (err error) {
	q := "DELETE FROM menu_roles WHERE role_id = ? AND id NOT IN (?) "
	if commodityId != nil {
		q += fmt.Sprintf(" AND commodity_id = '%v' ", *commodityId)
	} else {
		q += " AND commodity_id IS NULL "
	}

	query, args, err := sqlx.In(q, roleId, ids)
	query = tx.Rebind(query)

	if err != nil {
		log.Println(err)
		return
	}

	res, err := r.DB.Write.Exec(query, args...)
	if err != nil {
		return
	}

	_, err = res.RowsAffected()
	if err != nil {
		return
	}
	return
}