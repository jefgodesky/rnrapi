package helpers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"strconv"
)

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageStr := c.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)
		page = max(page, 1)

		maxPageSizeStr := os.Getenv("MAX_PAGE_SIZE")
		if maxPageSizeStr == "" {
			maxPageSizeStr = "500"
		}
		maxPageSize, _ := strconv.Atoi(maxPageSizeStr)

		pageSizeStr := c.DefaultQuery("page_size", "50")
		pageSize, _ := strconv.Atoi(pageSizeStr)
		pageSize = min(pageSize, maxPageSize)
		pageSize = max(pageSize, 1)

		c.Set("page", page)
		c.Set("page_size", pageSize)

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
