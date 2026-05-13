package auth

import (
	"strings"
	"github.com/gofrs/uuid"
)

type MenuService interface {
	ResolveAllMenu() ([]Menu, error)
	CreateMenu(req RequestMenuFormat) error
	UpdateMenu(id string, req RequestMenuFormat) error
	DeleteMenu(id string) error
	
	ResolveMenuByRoleID(roleID string) ([]MenuResponse, error)
	SaveBulkMenuRole(req RequestBulkMenuRole) error
}

type menuService struct {
	repo MenuRepository
}

func ProvideMenuService(repo MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) ResolveAllMenu() ([]Menu, error) {
	return s.repo.ResolveAllMenu()
}

func (s *menuService) CreateMenu(req RequestMenuFormat) error {
	newID, _ := uuid.NewV4()
	newMenu := Menu{
		ID:              newID,
		Name:            req.Name,
		Link:            req.Link,
		Icon:            req.Icon,
		Description:     req.Description,
		PermissionLabel: req.PermissionLabel,
		Action:          req.Action,
		Level:           req.Level,
		Seq:             req.Seq,
		ParentID:        req.ParentID,
	}
	return s.repo.CreateMenu(newMenu)
}

func (s *menuService) UpdateMenu(id string, req RequestMenuFormat) error {
	parsedID, _ := uuid.FromString(id)
	updateMenu := Menu{
		ID:              parsedID,
		Name:            req.Name,
		Link:            req.Link,
		Icon:            req.Icon,
		Description:     req.Description,
		PermissionLabel: req.PermissionLabel,
		Action:          req.Action,
		Level:           req.Level,
		Seq:             req.Seq,
		ParentID:        req.ParentID,
	}
	return s.repo.UpdateMenu(updateMenu)
}

func (s *menuService) DeleteMenu(id string) error {
	parsedID, _ := uuid.FromString(id)
	return s.repo.DeleteMenu(parsedID)
}

// Rekursif Tree Builder
func (s *menuService) ResolveMenuByRoleID(roleID string) ([]MenuResponse, error) {
	parsedRoleID, err := uuid.FromString(roleID)
	if err != nil { return nil, err }

	// 1. Ambil Menu Level 1 (Parent Tertinggi)
	parentMenus, err := s.repo.ResolveMenuByRoleID(parsedRoleID, 1, nil)
	if err != nil { return nil, err }

	// 2. Loop & Cari Anaknya (Level 2)
	for i := range parentMenus {
		if parentMenus[i].Permission != nil {
			parentMenus[i].PermissionList = strings.Split(*parentMenus[i].Permission, ",")
		}

		childMenus, _ := s.repo.ResolveMenuByRoleID(parsedRoleID, 2, &parentMenus[i].MenuID)
		
		for j := range childMenus {
			if childMenus[j].Permission != nil {
				childMenus[j].PermissionList = strings.Split(*childMenus[j].Permission, ",")
			}
		}
		parentMenus[i].Children = childMenus
	}

	return parentMenus, nil
}

func (s *menuService) SaveBulkMenuRole(req RequestBulkMenuRole) error {
	roleID, _ := uuid.FromString(req.RoleID)
	var listData []MenuRole

	for _, item := range req.Data {
		newID, _ := uuid.NewV4()
		menuID, _ := uuid.FromString(item.MenuID)
		listData = append(listData, MenuRole{
			ID:         newID,
			MenuID:     menuID,
			RoleID:     roleID,
			Permission: item.Permission,
		})
	}

	return s.repo.SaveBulkMenuRole(roleID, listData)
}