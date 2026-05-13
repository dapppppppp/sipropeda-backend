package master

import "github.com/gofrs/uuid"

type SumberDanaService interface {
	Create(req RequestSumberDana) error
	ResolveAll() ([]SumberDana, error)
	ResolveByID(id uuid.UUID) (SumberDana, error)
	Update(id string, req RequestSumberDana) error
	Delete(id string) error
}

type sumberDanaService struct {
	repo SumberDanaRepository
}

func ProvideSumberDanaService(repo SumberDanaRepository) SumberDanaService {
	return &sumberDanaService{repo: repo}
}

func (s *sumberDanaService) Create(req RequestSumberDana) error {
	newData := (&SumberDana{}).NewSumberDanaFormat(req)
	return s.repo.Create(newData)
}

func (s *sumberDanaService) ResolveAll() ([]SumberDana, error) {
	return s.repo.ResolveAll()
}

func (s *sumberDanaService) ResolveByID(id uuid.UUID) (SumberDana, error) {
	return s.repo.ResolveByID(id)
}

func (s *sumberDanaService) Update(id string, req RequestSumberDana) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	req.ID = parsedID
	updatedData := (&SumberDana{}).NewSumberDanaFormat(req)
	return s.repo.Update(updatedData)
}

func (s *sumberDanaService) Delete(id string) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	data := SumberDana{ID: parsedID}
	data.SoftDelete()
	return s.repo.Delete(data)
}