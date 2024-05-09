package bahancontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(c *gin.Context) {
	var bahan models.Bahan

	if err := c.BindJSON(&bahan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bahan.Nama == "" || bahan.Harga == 0 || bahan.Merk == "" || bahan.Satuan == "" || bahan.Stok == 0 {
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
	var bahan []models.Bahan
	nama := c.Query("nama")
	models.DB.Where("nama LIKE ?", "%"+nama+"%").Find(&bahan)
	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
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

	if bahan.Nama == "" || bahan.Harga == 0 || bahan.Merk == "" || bahan.Satuan == "" || bahan.Stok == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}

	models.DB.Save(&bahan)
	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
}