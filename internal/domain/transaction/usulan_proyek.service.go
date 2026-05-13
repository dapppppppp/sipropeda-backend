package transaction

import "github.com/gofrs/uuid"

type UsulanProyekService interface {
	Create(req RequestUsulanProyek) error
	ResolveAll() ([]UsulanProyek, error)
	ResolveByID(id uuid.UUID) (UsulanProyek, error)
	Update(id string, req RequestUsulanProyek) error
	Delete(id string, userID uuid.UUID) error
}

type usulanProyekService struct {
	repo UsulanProyekRepository
}

func ProvideUsulanProyekService(repo UsulanProyekRepository) UsulanProyekService {
	return &usulanProyekService{repo: repo}
}

func (s *usulanProyekService) Create(req RequestUsulanProyek) error {
	newData := (&UsulanProyek{}).NewUsulanProyekFormat(req)
	return s.repo.Create(newData)
}

func (s *usulanProyekService) ResolveAll() ([]UsulanProyek, error) {
	return s.repo.ResolveAll()
}

func (s *usulanProyekService) ResolveByID(id uuid.UUID) (UsulanProyek, error) {
	return s.repo.ResolveByID(id)
}

func (s *usulanProyekService) Update(id string, req RequestUsulanProyek) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	req.ID = parsedID
	updatedData := (&UsulanProyek{}).NewUsulanProyekFormat(req)
	return s.repo.Update(updatedData)
}

func (s *usulanProyekService) Delete(id string, userID uuid.UUID) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	data := UsulanProyek{ID: parsedID}
	data.SoftDelete(userID)
	return s.repo.Delete(data)
}