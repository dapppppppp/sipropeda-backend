package transaction

import (
	"errors"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/xuri/excelize/v2"
)

type UsulanProyekService interface {
	Create(req RequestUsulanProyek) error
	ResolveAll() ([]UsulanProyek, error)
	ResolveByID(id uuid.UUID) (UsulanProyek, error)
	Update(id string, req RequestUsulanProyek) error
	Delete(id string, userID uuid.UUID) error
	ImportExcelRKP(file multipart.File, tahunAnggaran int) (int, error)
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

// =================================================================================
// FUZZY COLUMN EXTRACTOR + AUTO MAPPING (SANGAT CERDAS)
// =================================================================================
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
	currentBidang := "" // Stateful Tracker untuk merekam "Bidang"

	for i, row := range rows {
		if i < 4 {
			continue
		}

		// A. DETEKSI BIDANG (Membaca baris judul Bidang dan menyimpannya di memori)
		for colIdx := 0; colIdx <= 2; colIdx++ {
			if colIdx < len(row) {
				val := strings.TrimSpace(row[colIdx])
				if strings.Contains(strings.ToLower(val), "bidang") && len(val) > 6 {
					currentBidang = val
					break
				}
			}
		}

		namaProyek := ""
		lokasi := "Desa"
		var nilaiRab float64 = 0
		sumberDana := ""

		// B. MENCARI NAMA PROYEK
		for colIdx := 2; colIdx <= 6; colIdx++ {
			if colIdx < len(row) {
				val := strings.TrimSpace(row[colIdx])
				valLower := strings.ToLower(val)
				if strings.Contains(valLower, "jenis kegiatan") || strings.Contains(valLower, "usulan") {
					continue
				}
				if len(val) > len(namaProyek) && len(val) > 5 {
					namaProyek = val
				}
			}
		}

		if namaProyek == "" || strings.Contains(strings.ToLower(namaProyek), "bidang") {
			continue // Skip jika ini hanya baris judul bidang
		}

		// C. MENCARI LOKASI
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

		// D. MENCARI SUMBER DANA (Cari di 6 kolom paling kanan)
		for colIdx := len(row) - 1; colIdx >= len(row)-6 && colIdx >= 0; colIdx-- {
			val := strings.TrimSpace(row[colIdx])
			valLower := strings.ToLower(val)
			if strings.Contains(valLower, "dd") || strings.Contains(valLower, "add") || strings.Contains(valLower, "pad") || strings.Contains(valLower, "swadaya") || strings.Contains(valLower, "dll") || strings.Contains(valLower, "bkk") {
				sumberDana = val
				break
			}
		}

		// E. MENCARI RAB
		for _, val := range row {
			valClean := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(val), "Rp", ""), ".", ""), ",", "")
			if rab, err := strconv.ParseFloat(strings.TrimSpace(valClean), 64); err == nil {
				if rab > nilaiRab {
					nilaiRab = rab
				}
			}
		}

		newID, _ := uuid.NewV4()
		usulan := UsulanProyek{
			ID:            newID,
			NamaProyek:    namaProyek,
			Lokasi:        lokasi,
			NilaiRAB:      nilaiRab,
			TahunAnggaran: tahunAnggaran,
			// TRIK CERDAS: Kita titipkan teks mentah Bidang & Sumber Dana ke properti string ini sementara, 
			// untuk nanti diolah oleh SQL Database menjadi UUID asli!
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