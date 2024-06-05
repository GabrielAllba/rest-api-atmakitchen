package laporanpenjualancontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductSale struct {
	ID         int     `json:"id"`
	Produk     string  `json:"produk"`
	Kuantitas  float64 `json:"kuantitas"`
	Harga      float64 `json:"harga"`
	JumlahUang float64 `json:"jumlah_uang"`
}

func GetByMonthAndYear(c *gin.Context) {
	// Retrieve the month and year from the request path parameters
	month := c.Param("month")
	year := c.Param("year")

	if month == "" || year == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Month and year are required"})
		return
	}

	// Combine month and year to create a date string
	dateStr := year + "-" + month + "-01"

	// Parse the date string to get the start date
	startDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM for month and YYYY for year"})
		return
	}
	endDate := startDate.AddDate(0, 1, 0) // One month later

	var sales []ProductSale

	// Query the database to retrieve transaction details for the specified month and year
	err = models.DB.Raw(`
		SELECT 
			td.id, 
			COALESCE(p.name, h.hampers_name) AS produk, 
			COALESCE(td.product_quantity, td.hampers_quantity) AS kuantitas, 
			COALESCE(td.product_price, td.hampers_price) AS harga, 
			COALESCE(td.product_quantity * td.product_price, td.hampers_quantity * td.hampers_price) AS jumlah_uang
		FROM 
			transaction_details td
		LEFT JOIN 
			products p ON p.id = td.product_id
		LEFT JOIN 
			hampers h ON h.id = td.hampers_id
		WHERE 
			td.tanggal_pengiriman >= ? AND td.tanggal_pengiriman < ?
	`, startDate, endDate).Scan(&sales).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve product sales"})
		return
	}

	// Return the product sales data as JSON
	c.JSON(http.StatusOK, sales)
}
