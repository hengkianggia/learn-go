package filters

import (
	"learn/internal/pkg/constants"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FilterFunc func(*gorm.DB, *gin.Context) *gorm.DB

func ApplyFilter(db *gorm.DB, c *gin.Context, filters ...FilterFunc) *gorm.DB {
	for _, filter := range filters {
		db = filter(db, c)
	}
	return db
}

func WithStatus() FilterFunc {
	return func(db *gorm.DB, c *gin.Context) *gorm.DB {
		if status := c.Query(constants.FilterEventStatus); status != "" {
			return db.Where("LOWER(status) = ?", strings.ToLower(status))
		}
		return db
	}
}

func WithSearch(columns ...string) FilterFunc {
	return func(db *gorm.DB, c *gin.Context) *gorm.DB {
		if search := c.Query(constants.FilterSearch); search != "" {
			for i, column := range columns {
				conditional := "LOWER(" + column + ") LIKE ?"
				if i == 0 {
					db = db.Where(conditional, "%"+strings.ToLower(search)+"%")
				} else {
					db = db.Or(conditional, "%"+strings.ToLower(search)+"%")
				}
			}
		}
		return db
	}
}

func WithDataRange(column string) FilterFunc {
	return func(db *gorm.DB, c *gin.Context) *gorm.DB {
		start := c.Query(constants.FilterStartDate)
		end := c.Query(constants.FilterEndDate)

		if start != "" && end != "" {
			return db.Where(column+" BETWEEN ? AND ?", start, end)
		}
		if start != "" {
			return db.Where(column+" >= ?", start)
		}
		if end != "" {
			return db.Where(column+" <= ?", end)
		}
		return db
	}
}
