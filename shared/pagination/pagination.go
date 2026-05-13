package pagination

import "math"

// Response format standar untuk data yang memiliki halaman
type Response struct {
	Items interface{} `json:"items"`
	Meta  Meta        `json:"meta"`
}

// Meta berisi informasi halaman
type Meta struct {
	TotalData   int `json:"totalData"`
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	TotalPage   int `json:"totalPage"`
}

// CreateMeta menghitung total halaman berdasarkan data
func CreateMeta(totalData, pageSize, pageNumber int) Meta {
	if pageSize == 0 {
		pageSize = 10 // Default jika tidak dikirim
	}
	if pageNumber == 0 {
		pageNumber = 1 // Default halaman pertama
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

	return Meta{
		TotalData:   totalData,
		PerPage:     pageSize,
		CurrentPage: pageNumber,
		TotalPage:   totalPage,
	}
}