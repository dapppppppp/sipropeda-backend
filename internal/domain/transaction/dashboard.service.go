package transaction

type DashboardService interface {
	GetStatistik(tahun int) (DashboardData, error)
}

type dashboardService struct {
	repo DashboardRepository
}

func ProvideDashboardService(repo DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetStatistik(tahun int) (DashboardData, error) {
	// Di sini Anda bisa menambahkan logika bisnis tambahan di masa depan jika diperlukan
	// sebelum atau sesudah memanggil repository.
	return s.repo.GetDashboardStatistik(tahun)
}