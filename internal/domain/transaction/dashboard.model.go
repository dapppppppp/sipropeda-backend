package transaction

// DashboardData merepresentasikan ringkasan statistik untuk halaman utama
type DashboardData struct {
	TotalUsulan   int         `json:"totalUsulan"`
	UsulanAPBDes  int         `json:"usulanApbdes"`
	TotalKriteria int         `json:"totalKriteria"`
	TotalPagu     float64     `json:"totalPagu"`
	Top5Usulan    []TopUsulan `json:"top5Usulan"`
}

// TopUsulan merepresentasikan baris data untuk tabel 5 prioritas teratas
type TopUsulan struct {
	Ranking          int     `db:"ranking" json:"ranking"`
	NamaProyek       string  `db:"nama_proyek" json:"namaProyek"`
	Lokasi           string  `db:"lokasi" json:"lokasi"`
	NilaiPreferensiV float64 `db:"nilai_preferensi_v" json:"nilaiPreferensiV"`
}