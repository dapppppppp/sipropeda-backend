package transaction

import "github.com/gofrs/uuid"

type PenilaianUsulanService interface {
	SaveBulk(req RequestBulkPenilaian) error
	GetByUsulanID(usulanID string) ([]PenilaianUsulan, error)
}

type penilaianUsulanService struct {
	repo PenilaianUsulanRepository
}

func ProvidePenilaianUsulanService(repo PenilaianUsulanRepository) PenilaianUsulanService {
	return &penilaianUsulanService{repo: repo}
}

func (s *penilaianUsulanService) SaveBulk(req RequestBulkPenilaian) error {
	var listData []PenilaianUsulan

	// Konversi Request ke Model Database
	for _, reqDetail := range req.Data {
		listData = append(listData, PenilaianUsulan{
			UsulanID:   req.UsulanID,
			KriteriaID: reqDetail.KriteriaID,
			NilaiInput: reqDetail.NilaiInput,
		})
	}

	return s.repo.SaveBulk(req.UsulanID, listData)
}

func (s *penilaianUsulanService) GetByUsulanID(usulanID string) ([]PenilaianUsulan, error) {
	parsedID, err := uuid.FromString(usulanID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByUsulanID(parsedID)
}