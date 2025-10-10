package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func Paginate(c *gin.Context, db *gorm.DB, model interface{}, out interface{}) (*PaginatedResponse, error) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	var totalItems int64
	db.Model(model).Count(&totalItems)

	offset := (page - 1) * perPage
	err := db.Offset(offset).Limit(perPage).Find(out).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	meta := Meta{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	return &PaginatedResponse{Data: out, Meta: meta}, nil
}
