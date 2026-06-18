package transaction

import (
	"sipropeda-backend/infras"
)

type DashboardRepository interface {
	GetDashboardStatistik(tahun int) (DashboardData, error)
}

type dashboardRepository struct {
	db *infras.PostgresqlConn
}

func ProvideDashboardRepository(db *infras.PostgresqlConn) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetDashboardStatistik(tahun int) (DashboardData, error) {
	var data DashboardData

	// 1. Hitung Total Usulan Tahun Ini
	r.db.Read.Get(&data.TotalUsulan, `SELECT COUNT(id) FROM usulan_proyek WHERE tahun_anggaran = $1 AND is_deleted = false`, tahun)

	// 2. Hitung Usulan yang sudah APBDes
	r.db.Read.Get(&data.UsulanAPBDes, `SELECT COUNT(id) FROM usulan_proyek WHERE tahun_anggaran = $1 AND status_tahapan = 'APBDes' AND is_deleted = false`, tahun)

	// 3. Hitung Kriteria Aktif
	r.db.Read.Get(&data.TotalKriteria, `SELECT COUNT(id) FROM m_kriteria WHERE is_active = true AND is_deleted = false`)

	// 4. Hitung Total Pagu Anggaran Definitif Tahun Ini
	r.db.Read.Get(&data.TotalPagu, `
		SELECT COALESCE(SUM(CASE WHEN pagu_definitif > 0 THEN pagu_definitif ELSE pagu_estimasi END), 0) 
		FROM pagu_anggaran WHERE tahun = $1 AND is_deleted = false`, tahun)

	// 5. Ambil Top 5 Ranking (Prioritas Tertinggi di tahap RAPBDes atau RKP terbaru)
	queryTop5 := `
		SELECT a.ranking, u.nama_proyek, u.lokasi, a.nilai_preferensi_v
		FROM arsip_perankingan a
		JOIN usulan_proyek u ON a.usulan_id = u.id
		WHERE u.tahun_anggaran = $1 AND u.is_deleted = false
		ORDER BY a.ranking ASC
		LIMIT 5
	`
	// Inisialisasi array kosong agar tidak me-return null di JSON jika data masih kosong
	data.Top5Usulan = make([]TopUsulan, 0)
	r.db.Read.Select(&data.Top5Usulan, queryTop5, tahun)

	return data, nil
}