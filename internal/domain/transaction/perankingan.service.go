package transaction

import (
	"encoding/json"
	"errors"
	"math"
	"sort"

	"sipropeda-backend/shared/model"
	"sipropeda-backend/shared/pagination"

	"github.com/gofrs/uuid"
)

type PerankinganService interface {
	HitungTOPSIS(req RequestHitungTopsis) ([]ArsipPerankingan, error)
	// Interface diperbarui
	GetArsip(req model.StandardRequest, tahun int, tahap string) (pagination.Response, error)
}

type perankinganService struct {
	repo PerankinganRepository
}

func ProvidePerankinganService(repo PerankinganRepository) PerankinganService {
	return &perankinganService{repo: repo}
}

func (s *perankinganService) HitungTOPSIS(req RequestHitungTopsis) ([]ArsipPerankingan, error) {
	belumDinilai, err := s.repo.CountUsulanBelumDinilai(req.TahunAnggaran, req.TahapVersi)
	if err != nil {
		return nil, errors.New("gagal melakukan pengecekan status penilaian usulan")
	}
	
	if belumDinilai > 0 {
		return nil, errors.New("GAGAL: Terdapat usulan proyek yang belum dinilai. Silakan lengkapi semua penilaian terlebih dahulu agar perhitungan akurat!")
	}

	kriteria, err := s.repo.GetKriteriaAktif()
	if err != nil || len(kriteria) == 0 {
		return nil, errors.New("data kriteria tidak ditemukan")
	}

	matriksMentah, err := s.repo.GetMatriksPenilaian(req.TahunAnggaran, req.TahapVersi)
	if err != nil || len(matriksMentah) == 0 {
		return nil, errors.New("data usulan atau penilaian belum diisi untuk tahap ini")
	}

	mapNilai := make(map[string]map[string]float64)
	for _, m := range matriksMentah {
		uID := m.UsulanID.String()
		kID := m.KriteriaID.String()
		if mapNilai[uID] == nil {
			mapNilai[uID] = make(map[string]float64)
		}
		mapNilai[uID][kID] = m.NilaiInput
	}

	pembagi := make(map[string]float64)
	for _, k := range kriteria {
		kID := k.ID.String()
		var totalKuadrat float64
		for _, nilaiUsulan := range mapNilai {
			val := nilaiUsulan[kID]
			totalKuadrat += val * val
		}
		pembagi[kID] = math.Sqrt(totalKuadrat)
	}

	matriksY := make(map[string]map[string]float64)
	for uID, nilaiUsulan := range mapNilai {
		matriksY[uID] = make(map[string]float64)
		for _, k := range kriteria {
			kID := k.ID.String()
			var y float64 = 0
			if pembagi[kID] != 0 {
				y = (nilaiUsulan[kID] / pembagi[kID]) * k.Bobot
			}
			matriksY[uID][kID] = y
		}
	}

	idealPositif := make(map[string]float64)
	idealNegatif := make(map[string]float64)

	for _, k := range kriteria {
		kID := k.ID.String()
		var maxVal, minVal float64
		first := true

		for _, yVal := range matriksY {
			v := yVal[kID]
			if first {
				maxVal, minVal = v, v
				first = false
			} else {
				if v > maxVal { maxVal = v }
				if v < minVal { minVal = v }
			}
		}

		if k.Jenis == "benefit" {
			idealPositif[kID] = maxVal
			idealNegatif[kID] = minVal
		} else { // cost
			idealPositif[kID] = minVal
			idealNegatif[kID] = maxVal
		}
	}

	var hasilAkhir []ArsipPerankingan
	for uID, yVal := range matriksY {
		var totalDPlus, totalDMin float64

		for _, k := range kriteria {
			kID := k.ID.String()
			val := yVal[kID]
			totalDPlus += math.Pow(val-idealPositif[kID], 2)
			totalDMin += math.Pow(val-idealNegatif[kID], 2)
		}

		dPlus := math.Sqrt(totalDPlus)
		dMin := math.Sqrt(totalDMin)

		var nilaiV float64 = 0
		if (dPlus + dMin) != 0 {
			nilaiV = dMin / (dPlus + dMin)
		}

		detail := map[string]interface{}{
			"d_plus": dPlus,
			"d_min":  dMin,
			"nilai_Y": yVal,
		}
		detailJSON, _ := json.Marshal(detail)
		detailStr := string(detailJSON)

		usulanUUID, _ := uuid.FromString(uID)
		hasilAkhir = append(hasilAkhir, ArsipPerankingan{
			UsulanID:         usulanUUID,
			NilaiPreferensiV: nilaiV,
			TahapVersi:       req.TahapVersi,
			DetailKalkulasi:  &detailStr,
		})
	}

	sort.Slice(hasilAkhir, func(i, j int) bool {
		return hasilAkhir[i].NilaiPreferensiV > hasilAkhir[j].NilaiPreferensiV
	})

	for i := range hasilAkhir {
		hasilAkhir[i].Ranking = i + 1
	}

	err = s.repo.SaveHasilPerankingan(hasilAkhir)
	if err != nil {
		return nil, err
	}

	return hasilAkhir, nil
}

// Diperbarui agar memanggil fungsi pagination dari repo
func (s *perankinganService) GetArsip(req model.StandardRequest, tahun int, tahap string) (pagination.Response, error) {
	return s.repo.GetArsip(req, tahun, tahap)
}