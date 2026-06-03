package auth

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type AppConfigService interface {
	ResolveByID(id string) (data AppConfig, err error)
	ResolveDTOByID(id string) (data AppConfigDTOByID, err error)
	Update(id string, req RequestAppConfigFormat) error
	UploadFile(r *http.Request) (string, error)
}

type appConfigServiceImpl struct {
	repo AppConfigRepository
}

func ProvideAppConfigService(repo AppConfigRepository) AppConfigService {
	return &appConfigServiceImpl{repo: repo}
}

func (s *appConfigServiceImpl) ResolveByID(id string) (data AppConfig, err error) {
	return s.repo.ResolveByID(id)
}

func (s *appConfigServiceImpl) ResolveDTOByID(id string) (data AppConfigDTOByID, err error) {
	return s.repo.ResolveDTOByID(id)
}

func (s *appConfigServiceImpl) Update(id string, req RequestAppConfigFormat) error {
	data, err := s.repo.ResolveByID(id)
	if err != nil {
		return err
	}

	data.AppConfigFormat(req)
	return s.repo.Update(data)
}

func (s *appConfigServiceImpl) UploadFile(r *http.Request) (string, error) {
	// Batasi ukuran file (Max 10 MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 1. Buat folder files/config-app jika belum ada (Persis seperti CHM)
	uploadDir := filepath.Join("files", "config-app")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// 2. Ambil ekstensi asli dari gambar (misal: .png atau .jpg)
	ext := filepath.Ext(handler.Filename)

	// 3. Generate nama file menggunakan UUID
	// Hasilnya: "16f8e087-364b-4324-bda8-c84cb9393670.png"
	newFileName := uuid.New().String() + ext
	
	// Jalur simpan ke hardisk komputer/server
	savePath := filepath.Join(uploadDir, newFileName)

	// Simpan file ke direktori
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// 4. Return path untuk disimpan ke database
	// Menggunakan forward slash (/) agar format path-nya seragam dan aman dibaca di Frontend
	dbPath := "files/config-app/" + newFileName

	return dbPath, nil
}