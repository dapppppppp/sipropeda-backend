package transaction

import "github.com/gofrs/uuid"

type PaguAnggaranService interface {
	Create(req RequestPaguAnggaran) error
	ResolveAll() ([]PaguAnggaran, error)
	ResolveByID(id uuid.UUID) (PaguAnggaran, error)
	Update(id string, req RequestPaguAnggaran) error
	Delete(id string, userID uuid.UUID) error
}

type paguAnggaranService struct {
	repo PaguAnggaranRepository
}

func ProvidePaguAnggaranService(repo PaguAnggaranRepository) PaguAnggaranService {
	return &paguAnggaranService{repo: repo}
}

func (s *paguAnggaranService) Create(req RequestPaguAnggaran) error {
	newData := (&PaguAnggaran{}).NewPaguAnggaranFormat(req)
	return s.repo.Create(newData)
}

func (s *paguAnggaranService) ResolveAll() ([]PaguAnggaran, error) {
	return s.repo.ResolveAll()
}

func (s *paguAnggaranService) ResolveByID(id uuid.UUID) (PaguAnggaran, error) {
	return s.repo.ResolveByID(id)
}

func (s *paguAnggaranService) Update(id string, req RequestPaguAnggaran) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	req.ID = parsedID
	updatedData := (&PaguAnggaran{}).NewPaguAnggaranFormat(req)
	return s.repo.Update(updatedData)
}

func (s *paguAnggaranService) Delete(id string, userID uuid.UUID) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	data := PaguAnggaran{ID: parsedID}
	data.SoftDelete(userID)
	return s.repo.Delete(data)
}