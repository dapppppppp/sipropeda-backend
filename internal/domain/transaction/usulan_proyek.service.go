package transaction

import (
	"errors"
	"mime/multipart"
	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/xuri/excelize/v2"
)

type UsulanProyekService interface {
	Create(req RequestUsulanProyek) error
	ResolveAll(req model.StandardRequest) (pagination.Response, error)
	ResolveByID(id uuid.UUID) (UsulanProyek, error)
	Update(id string, req RequestUsulanProyek) error
	Delete(id string, userID uuid.UUID) error
	ImportExcelRKP(file multipart.File, tahunAnggaran int) (int, error)
	GetAllData() ([]UsulanProyek, error) // Tambahkan ini di interface
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

func (s *usulanProyekService) ResolveAll(req model.StandardRequest) (pagination.Response, error) {
	return s.repo.ResolveAll(req)
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

func (s *usulanProyekService) GetAllData() ([]UsulanProyek, error) {
	return s.repo.GetAllData()
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

func (s *usulanProyekService) ImportExcelRKP(file multipart.File, tahunAnggaran int) (int, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return 0, errors.New("gagal membaca file excel")
	}
	defer f.Close()

	var targetSheet string
	for _, name := range f.GetSheetList() {
		if strings.Contains(strings.ToLower(name), "usulan") || strings.Contains(strings.ToLower(name), "rkp") {
			targetSheet = name
			break
		}
	}
	if targetSheet == "" {
		targetSheet = f.GetSheetName(0)
	}

	rows, err := f.GetRows(targetSheet)
	if err != nil {
		return 0, errors.New("gagal membaca baris pada sheet data")
	}

	var dataImport []UsulanProyek
	currentBidang := ""

	for i, row := range rows {
		if i < 4 {
			continue
		}

		nonEmptyCols := 0
		for _, col := range row {
			if strings.TrimSpace(col) != "" {
				nonEmptyCols++
			}
		}

		for colIdx := 0; colIdx <= 2; colIdx++ {
			if colIdx < len(row) {
				val := strings.TrimSpace(row[colIdx])
				if strings.Contains(strings.ToLower(val), "bidang") && len(val) > 6 {
					currentBidang = val
					break
				}
			}
		}

		// Baris proyek yang valid biasanya memiliki setidaknya 3 kolom terisi (No, Kegiatan, RAB/Lokasi)
		// Jika kurang dari 3, kemungkinan besar itu adalah judul, nama daerah, atau tanda tangan.
		if nonEmptyCols < 3 {
			continue
		}

		namaProyek := ""
		lokasi := "Desa"
		var nilaiRab float64 = 0
		sumberDana := ""

		for colIdx := 2; colIdx <= 6; colIdx++ {
			if colIdx < len(row) {
				val := strings.TrimSpace(row[colIdx])
				valLower := strings.ToLower(val)
				if strings.Contains(valLower, "jenis kegiatan") || strings.Contains(valLower, "usulan") {
					continue
				}

				// Filter kata-kata yang bukan usulan proyek (misal: header, nama desa, dll)
				isInvalid := false
				invalidPrefixes := []string{
					"pemerintah", "lampiran", "rencana kerja", "rkp", "tahun anggaran",
					"jumlah", "total", "mengetahui", "disetujui", "kepala desa", "sekretaris",
				}
				for _, kw := range invalidPrefixes {
					if strings.HasPrefix(valLower, kw) {
						isInvalid = true
						break
					}
				}

				// Filter tambahan untuk "Desa XYZ", "Kecamatan ABC", "Kabupaten DEF"
				words := strings.Fields(valLower)
				if len(words) > 0 && len(words) <= 4 {
					fw := words[0]
					if fw == "desa" || fw == "kecamatan" || fw == "kabupaten" || fw == "provinsi" {
						// Pengecualian jika mengandung kata kerja/kegiatan yang valid
						if !strings.Contains(valLower, "wisata") && !strings.Contains(valLower, "pengembangan") {
							isInvalid = true
						}
					}
				}

				if isInvalid {
					continue
				}

				if len(val) > len(namaProyek) && len(val) > 5 {
					namaProyek = val
				}
			}
		}

		if namaProyek == "" || strings.Contains(strings.ToLower(namaProyek), "bidang") {
			continue
		}

		for colIdx := 4; colIdx <= 10; colIdx++ {
			if colIdx < len(row) {
				val := strings.TrimSpace(row[colIdx])
				valLower := strings.ToLower(val)
				if valLower == "desa" || strings.HasPrefix(valLower, "dusun") || strings.HasPrefix(valLower, "rt") || strings.HasPrefix(valLower, "rw") {
					lokasi = val
					break
				}
			}
		}

		for colIdx := len(row) - 1; colIdx >= len(row)-6 && colIdx >= 0; colIdx-- {
			val := strings.TrimSpace(row[colIdx])
			valLower := strings.ToLower(val)
			if strings.Contains(valLower, "dd") || strings.Contains(valLower, "add") || strings.Contains(valLower, "pad") || strings.Contains(valLower, "swadaya") || strings.Contains(valLower, "dll") || strings.Contains(valLower, "bkk") {
				sumberDana = val
				break
			}
		}

		for _, val := range row {
			valClean := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(val), "Rp", ""), ".", ""), ",", "")
			if rab, err := strconv.ParseFloat(strings.TrimSpace(valClean), 64); err == nil {
				if rab > nilaiRab {
					nilaiRab = rab
				}
			}
		}

		if nilaiRab == 0 {
			continue
		}

		newID, _ := uuid.NewV4()
		usulan := UsulanProyek{
			ID:            newID,
			NamaProyek:    namaProyek,
			Lokasi:        lokasi,
			NilaiRAB:      nilaiRab,
			TahunAnggaran: tahunAnggaran,
			StatusTahapan: currentBidang, 
			StatusSifat:   sumberDana,    
		}
		dataImport = append(dataImport, usulan)
	}

	if len(dataImport) == 0 {
		return 0, errors.New("tidak ada data proyek yang berhasil diekstrak")
	}

	err = s.repo.BulkInsert(dataImport)
	if err != nil {
		return 0, err
	}
	return len(dataImport), nil
}