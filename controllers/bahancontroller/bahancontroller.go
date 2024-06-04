package bahancontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(c *gin.Context) {
	var bahan models.Bahan

	if err := c.BindJSON(&bahan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bahan.Nama == "" || bahan.Satuan == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}

	var existingBahan models.Bahan

	if err := models.DB.Where("name = ?", bahan.Nama).First(&existingBahan).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bahan dengan nama tersebut sudah ada"})
		return
	}

	if err := models.DB.Create(&bahan).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat bahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
}

func Index(c *gin.Context) {
	var bahan []models.Bahan
	models.DB.Find(&bahan)
	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
}

func Show(c *gin.Context) {
	var bahan models.Bahan
	id := c.Param("id")
	if err := models.DB.First(&bahan, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Bahan not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
}

func Delete(c *gin.Context) {
	var bahan models.Bahan
	id := c.Param("id")
	if err := models.DB.First(&bahan, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Bahan not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	models.DB.Delete(&bahan)
	c.JSON(http.StatusOK, gin.H{"message": "Bahan deleted"})
}

func Search(c *gin.Context) {
	query := c.Query("query")
	var bahans []models.Bahan

	
	query = strings.ToLower(query)
	result := models.DB.Where("LOWER(nama) LIKE ?  OR stok LIKE ? OR satuan LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	

	
	if err := result.Find(&bahans).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "No bahan found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"bahans": bahans})
}

func Update(c *gin.Context) {
	var bahan models.Bahan
	id := c.Param("id")
	if err := models.DB.First(&bahan, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Bahan not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	if err := c.BindJSON(&bahan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bahan.Nama == ""  || bahan.Satuan == ""  {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}

	models.DB.Save(&bahan)
	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
}
func KurangiStock(c *gin.Context) {
    var bahan models.Bahan
    id := c.Param("id")
    quantityStr := c.Param("quantity")

    // Convert quantity to an integer
    quantity, err := strconv.Atoi(quantityStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
        return
    }

    // Find the bahan record by its ID
    if err := models.DB.Where("id = ?", id).First(&bahan).Error; err != nil {
        
        return
    }

    // Check if the stock is sufficient
    if int(bahan.Stok) < quantity {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
        return
    }

    // Subtract the quantity from the stock
    bahan.Stok -= float64(quantity)

    // Save the updated bahan record back to the database
    if err := models.DB.Save(&bahan).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully", "bahan": bahan})
}
