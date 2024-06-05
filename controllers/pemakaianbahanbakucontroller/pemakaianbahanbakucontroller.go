package pemakaianBahanBakuController

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(c *gin.Context) {
	var pemakaianBahanBaku models.PemakaianBahanBaku

	if err := c.BindJSON(&pemakaianBahanBaku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Create(&pemakaianBahanBaku).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat pemakaian bahan baku"})
		return
	}


	if err := models.DB.Preload("Bahan").First(&pemakaianBahanBaku, pemakaianBahanBaku.Id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data bahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pemakaian_bahan_baku": pemakaianBahanBaku})
}

func Show(c *gin.Context) {
	var pemakaianBahanBaku models.PemakaianBahanBaku
	id := c.Param("id")

	if err := models.DB.Preload("Bahan").First(&pemakaianBahanBaku, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "pemakaian_bahan_baku not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"pemakaian_bahan_baku": pemakaianBahanBaku})
}

func Index(c *gin.Context) {
	var pemakaianBahanBaku []models.PemakaianBahanBaku
	if err := models.DB.Preload("Bahan").Preload("TransactionDetail").Find(&pemakaianBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pemakaian_bahan_baku": pemakaianBahanBaku})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var pemakaianBahanBaku models.PemakaianBahanBaku
	if err := models.DB.First(&pemakaianBahanBaku, intId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var bahan models.Bahan
	if err := models.DB.First(&bahan, pemakaianBahanBaku.BahanId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load related Bahan"})
		return
	}

	bahan.Stok += pemakaianBahanBaku.Jumlah
	if err := models.DB.Save(&bahan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Bahan stock"})
		return
	}

	if err := models.DB.Delete(&pemakaianBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

func Update(c *gin.Context) {
	var pemakaianBahanBaku models.PemakaianBahanBaku
	id := c.Param("id")

	if err := c.BindJSON(&pemakaianBahanBaku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingPemakaian models.PemakaianBahanBaku
	if err := models.DB.First(&existingPemakaian, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Record not found"})
		return
	}

	var bahan models.Bahan
	if err := models.DB.First(&bahan, existingPemakaian.BahanId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data bahan"})
		return
	}

	bahan.Stok += existingPemakaian.Jumlah - pemakaianBahanBaku.Jumlah
	if bahan.Stok < 0 {
		bahan.Stok = 0
	}

	if err := models.DB.Save(&bahan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui stok bahan"})
		return
	}

	if err := models.DB.Model(&existingPemakaian).Updates(pemakaianBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui pemakaian bahan baku"})
		return
	}

	if err := models.DB.Preload("Bahan").First(&pemakaianBahanBaku, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data bahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pemakaian_bahan_baku": pemakaianBahanBaku})
}
