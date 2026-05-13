package auth

import "github.com/gofrs/uuid"

type RoleService interface {
	CreateRole(req RequestRole) error
	ResolveAll() ([]Role, error)
	ResolveByID(id uuid.UUID) (Role, error)
	UpdateRole(id string, req RequestRole) error
	DeleteRole(id string) error
}

type roleService struct {
	repo RoleRepository
}

func ProvideRoleService(repo RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(req RequestRole) error {
	newRole := (&Role{}).NewRoleFormat(req)
	return s.repo.Create(newRole)
}

func (s *roleService) ResolveAll() ([]Role, error) {
	return s.repo.ResolveAll()
}

func (s *roleService) ResolveByID(id uuid.UUID) (Role, error) {
	return s.repo.ResolveByID(id)
}

func (s *roleService) UpdateRole(id string, req RequestRole) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}
	req.ID = parsedID
	updatedRole := (&Role{}).NewRoleFormat(req)
	return s.repo.Update(updatedRole)
}

func (s *roleService) DeleteRole(id string) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}
	role := Role{ID: parsedID}
	role.SoftDelete()
	return s.repo.Delete(role)
}