package utils

import (
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Menyimpan parameter pagination dari query URL
// contoh: /users?page=2&limit=10&filter=john&sort=-created_at
type PaginationQuery struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Filter string `json:"filter"`
	Sort   string `json:"sort"`
}

// Metadata pagination yang dikirim di response API
// berisi informasi halaman, total data, dll
type PaginationMeta struct {
	Page      int    `json:"page" example:"1"`
	Limit     int    `json:"limit" example:"10"`
	Total     int64  `json:"total" example:"100"`
	TotalPage int    `json:"total_pages" example:"10"`
	Filter    string `json:"filter,omitempty" example:"john"`
	Sort      string `json:"sort,omitempty" example:"-created_at"`
}

// Generic result pagination supaya bisa dipakai untuk tipe data apa saja
// contoh: PaginationResult[UserResponse]
type PaginationResult[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// Mengambil query pagination dari URL lalu mengubahnya
// menjadi struct PaginationQuery yang siap dipakai repository
func ParsePaginationQuery(c *fiber.Ctx) PaginationQuery {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	filter := strings.TrimSpace(c.Query("filter", ""))
	sort := strings.TrimSpace(c.Query("sort", ""))

	return PaginationQuery{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
		Filter: filter,
		Sort:   sort,
	}
}

// Membuat metadata pagination berdasarkan total data dari database
// digunakan untuk response API
func BuildPaginationMeta(q PaginationQuery, total int64) PaginationMeta {

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(q.Limit)))
	}

	return PaginationMeta{
		Page:      q.Page,
		Limit:     q.Limit,
		Total:     total,
		TotalPage: totalPage,
		Filter:    q.Filter,
		Sort:      q.Sort,
	}
}

// Konfigurasi untuk menentukan bagaimana pagination bekerja
// seperti field yang bisa difilter dan disort
type PaginationConfig struct {
	SearchFields      []string
	AllowedSortFields map[string]string
	DefaultSortField  string
	DefaultSortOrder  string
	UseILIKE          bool
}

// Menerapkan filter, sorting, limit dan offset ke query GORM
// supaya repository tidak perlu menulis ulang logic pagination
func ApplyPagination(db *gorm.DB, q PaginationQuery, cfg PaginationConfig) *gorm.DB {
	db = applyFilter(db, q.Filter, cfg)
	db = applySort(db, q.Sort, cfg)
	db = db.Limit(q.Limit).Offset(q.Offset)

	return db
}

// Menghitung total data di database setelah filter diterapkan
// digunakan untuk menentukan total halaman
func CountWithFilter(db *gorm.DB, q PaginationQuery, cfg PaginationConfig) (int64, error) {

	var total int64

	filteredDB := applyFilter(db, q.Filter, cfg)
	if err := filteredDB.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Menambahkan kondisi filter ke query database
// contoh: WHERE name ILIKE '%john%' OR email ILIKE '%john%'
func applyFilter(db *gorm.DB, filter string, cfg PaginationConfig) *gorm.DB {

	if strings.TrimSpace(filter) == "" || len(cfg.SearchFields) == 0 {
		return db
	}

	pattern := "%" + strings.TrimSpace(filter) + "%"
	conditions := make([]string, 0, len(cfg.SearchFields))
	args := make([]interface{}, 0, len(cfg.SearchFields))

	operator := "LIKE"
	if cfg.UseILIKE {
		operator = "ILIKE"
	}

	for _, field := range cfg.SearchFields {
		conditions = append(conditions, field+" "+operator+" ?")
		args = append(args, pattern)
	}

	return db.Where(strings.Join(conditions, " OR "), args...)
}

// Menambahkan sorting ke query database dengan whitelist field
// supaya aman dari sort yang tidak diizinkan
func applySort(db *gorm.DB, sort string, cfg PaginationConfig) *gorm.DB {

	field := cfg.DefaultSortField
	order := strings.ToUpper(strings.TrimSpace(cfg.DefaultSortOrder))
	if order != "ASC" {
		order = "DESC"
	}

	if strings.TrimSpace(sort) != "" {
		rawSort := strings.TrimSpace(sort)
		rawOrder := "ASC"

		if strings.HasPrefix(rawSort, "-") {
			rawOrder = "DESC"
			rawSort = strings.TrimPrefix(rawSort, "-")
		}

		if column, ok := cfg.AllowedSortFields[rawSort]; ok {
			field = column
			order = rawOrder
		}
	}

	return db.Order(field + " " + order)
}
