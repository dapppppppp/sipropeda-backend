package transaction

import (
	"bytes"
	"sipropeda-backend/infras"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type PerankinganRepository interface {
	GetKriteriaAktif() ([]KriteriaTopsis, error)
	GetMatriksPenilaian(tahun int, tahap string) ([]MatriksPenilaian, error)
	SaveHasilPerankingan(data []ArsipPerankingan) error
	
	// Diperbarui agar mendukung pagination
	GetArsip(req model.StandardRequest, tahun int, tahap string) (pagination.Response, error)
	
	CountUsulanBelumDinilai(tahun int, tahap string) (int, error)
}

type perankinganRepository struct {
	db *infras.PostgresqlConn
}

func ProvidePerankinganRepository(db *infras.PostgresqlConn) PerankinganRepository {
	return &perankinganRepository{db: db}
}

func (r *perankinganRepository) GetKriteriaAktif() ([]KriteriaTopsis, error) {
	var data []KriteriaTopsis
	query := `SELECT id, kode, bobot, jenis FROM m_kriteria WHERE is_deleted = false ORDER BY kode ASC`
	err := r.db.Read.Select(&data, query)
	return data, err
}

func (r *perankinganRepository) GetMatriksPenilaian(tahun int, tahap string) ([]MatriksPenilaian, error) {
	var data []MatriksPenilaian
	query := `
		SELECT p.usulan_id, p.kriteria_id, p.nilai_input
		FROM penilaian_usulan p
		JOIN usulan_proyek u ON p.usulan_id = u.id
		WHERE u.tahun_anggaran = $1 AND u.status_tahapan::text = $2::text AND u.is_deleted = false
	`
	err := r.db.Read.Select(&data, query, tahun, tahap)
	return data, err
}

func (r *perankinganRepository) SaveHasilPerankingan(data []ArsipPerankingan) error {
	if len(data) == 0 {
		return nil
	}

	deleteQuery := `
		DELETE FROM arsip_perankingan 
		WHERE usulan_id IN (
			SELECT id FROM usulan_proyek WHERE status_tahapan::text = $1::text
		) AND tahap_versi::text = $1::text
	`
	_, err := r.db.Write.Exec(deleteQuery, data[0].TahapVersi)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO arsip_perankingan (id, usulan_id, nilai_preferensi_v, ranking, tahap_versi, detail_kalkulasi)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, val := range data {
		newID, _ := uuid.NewV4()
		_, err = r.db.Write.Exec(insertQuery, newID, val.UsulanID, val.NilaiPreferensiV, val.Ranking, val.TahapVersi, val.DetailKalkulasi)
		if err != nil {
			return err
		}
	}
	return nil
}

// LOGIKA GET ARSIP MENGGUNAKAN PAGINATION DAN DYNAMIC QUERY
func (r *perankinganRepository) GetArsip(req model.StandardRequest, tahun int, tahap string) (data pagination.Response, err error) {
	var searchParams []interface{}
	var filterBuff bytes.Buffer

	// Base Condition
	filterBuff.WriteString(" WHERE u.tahun_anggaran = ? AND a.tahap_versi::text = ? AND u.status_tahapan::text = ? AND coalesce(u.is_deleted, false) = false")
	searchParams = append(searchParams, tahun, tahap, tahap)

	// Filter Keyword
	if req.Keyword != "" {
		filterBuff.WriteString(" AND u.nama_proyek ILIKE ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	selectDto := `
		SELECT a.id, a.usulan_id, u.nama_proyek as usulan_name, a.nilai_preferensi_v, a.ranking, a.tahap_versi, a.created_at
		FROM arsip_perankingan a
		JOIN usulan_proyek u ON a.usulan_id = u.id
	`

	// 1. Menghitung Total Data
	countQuery := r.db.Read.Rebind("SELECT count(*) FROM (" + selectDto + filterBuff.String() + ")x")
	var totalData int
	err = r.db.Read.QueryRowx(countQuery, searchParams...).Scan(&totalData)
	if err != nil {
		return
	}

	if totalData < 1 {
		data.Items = []ArsipPerankingan{} // Slice kosong
		data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
		return data, nil
	}

	// 2. Sorting
	orderByColumn, ok := ColumnMapArsipPerankingan[req.SortBy].(string)
	if !ok {
		orderByColumn = "a.ranking"
	}
	filterBuff.WriteString(" ORDER BY " + orderByColumn + " " + req.SortType)

	// 3. Limit Offset
	offset := (req.PageNumber - 1) * req.PageSize
	filterBuff.WriteString(" LIMIT ? OFFSET ? ")
	searchParams = append(searchParams, req.PageSize, offset)

	searchQuery := r.db.Read.Rebind(selectDto + filterBuff.String())
	rows, err := r.db.Read.Queryx(searchQuery, searchParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	// 4. Ekstraksi ke slice lokal
	var items []ArsipPerankingan
	for rows.Next() {
		var item ArsipPerankingan
		err = rows.StructScan(&item)
		if err != nil {
			return
		}
		items = append(items, item)
	}

	data.Items = items
	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return data, nil
}

func (r *perankinganRepository) CountUsulanBelumDinilai(tahun int, tahap string) (int, error) {
	var count int
	query := `
		SELECT COUNT(id) 
		FROM usulan_proyek 
		WHERE tahun_anggaran = $1 AND status_tahapan::text = $2::text AND is_deleted = false
		AND NOT EXISTS (SELECT 1 FROM penilaian_usulan pu WHERE pu.usulan_id = usulan_proyek.id)
	`
	err := r.db.Read.Get(&count, query, tahun, tahap)
	return count, err
}