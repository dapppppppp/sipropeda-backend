package master

import (
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type KriteriaService interface {
	Create(req RequestKriteriaFormat) error
	ResolveAll() ([]Kriteria, error)
	ResolvePaging(req model.StandardRequest) (pagination.Response, error)
	ResolveByID(id uuid.UUID) (Kriteria, error)
	Update(id string, req RequestKriteriaFormat) error
	Delete(id string, userID uuid.UUID) error
}

type kriteriaService struct {
	repository KriteriaRepository
}

func ProvideKriteriaService(repository KriteriaRepository) KriteriaService {
	return &kriteriaService{repository: repository}
}

func (s *kriteriaService) Create(req RequestKriteriaFormat) error {
	newKriteria, _ := (&Kriteria{}).NewKriteriaFormat(req)
	return s.repository.Create(newKriteria)
}

func (s *kriteriaService) ResolveAll() ([]Kriteria, error) {
	return s.repository.ResolveAll()
}

func (s *kriteriaService) ResolvePaging(req model.StandardRequest) (pagination.Response, error) {
	return s.repository.ResolvePaging(req)
}

func (s *kriteriaService) ResolveByID(id uuid.UUID) (Kriteria, error) {
	return s.repository.ResolveByID(id)
}

func (s *kriteriaService) Update(id string, req RequestKriteriaFormat) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	req.ID = parsedID
	updatedKriteria, _ := (&Kriteria{}).NewKriteriaFormat(req)

	return s.repository.Update(updatedKriteria)
}

func (s *kriteriaService) Delete(id string, userID uuid.UUID) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	kriteria := Kriteria{ID: parsedID}
	kriteria.SoftDelete(userID)

	return s.repository.Delete(kriteria)
}