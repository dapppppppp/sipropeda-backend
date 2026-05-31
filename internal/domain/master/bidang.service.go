package master

import "github.com/gofrs/uuid"

type BidangPembangunanService interface {
	Create(req RequestBidangPembangunan) error
	ResolveAll() ([]BidangPembangunan, error)
	ResolveByID(id uuid.UUID) (BidangPembangunan, error)
	Update(id string, req RequestBidangPembangunan) error
	Delete(id string, userID uuid.UUID) error
}

type bidangPembangunanService struct {
	repo BidangPembangunanRepository
}

func ProvideBidangPembangunanService(repo BidangPembangunanRepository) BidangPembangunanService {
	return &bidangPembangunanService{repo: repo}
}

func (s *bidangPembangunanService) Create(req RequestBidangPembangunan) error {
	newData := (&BidangPembangunan{}).NewBidangPembangunanFormat(req)
	return s.repo.Create(newData)
}

func (s *bidangPembangunanService) ResolveAll() ([]BidangPembangunan, error) {
	return s.repo.ResolveAll()
}

func (s *bidangPembangunanService) ResolveByID(id uuid.UUID) (BidangPembangunan, error) {
	return s.repo.ResolveByID(id)
}

func (s *bidangPembangunanService) Update(id string, req RequestBidangPembangunan) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	req.ID = parsedID
	updatedData := (&BidangPembangunan{}).NewBidangPembangunanFormat(req)
	return s.repo.Update(updatedData)
}

func (s *bidangPembangunanService) Delete(id string, userID uuid.UUID) error {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	data := BidangPembangunan{ID: parsedID}
	data.SoftDelete(userID)
	return s.repo.Delete(data)
}