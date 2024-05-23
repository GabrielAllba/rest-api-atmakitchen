package pengeluaranlaincontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//create, index, show, delete, search, update

func Create(c *gin.Context) {
	var pengeluaranLain models.PengeluaranLain

	if err := c.BindJSON(&pengeluaranLain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if pengeluaranLain.Harga == 0 || pengeluaranLain.Deskripsi == "" || pengeluaranLain.Metode == "" || pengeluaranLain.TanggalPengeluaran == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": pengeluaranLain})
		return
	}

	if err := models.DB.Create(&pengeluaranLain).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal create pengeluaran lain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pengeluaran_lain": pengeluaranLain})
}

func Index(c *gin.Context) {
	var pengeluaranLain []models.PengeluaranLain
	models.DB.Find(&pengeluaranLain)
	c.JSON(http.StatusOK, gin.H{"pengeluaran_lain": pengeluaranLain})
}

func Show(c *gin.Context) {
	var pengeluaranLain models.PengeluaranLain
	id := c.Param("id")

	if err := models.DB.First(&pengeluaranLain, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Pengeluaran_Lain not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"pengeluaran_lain": pengeluaranLain})
}

func Delete(c *gin.Context) {
	var pengeluaranLain models.PengeluaranLain
	id := c.Param("id")

	if err := models.DB.First(&pengeluaranLain, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Pengeluaran Lain not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	models.DB.Delete(&pengeluaranLain)
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil Delete Pengeluaran Lain"})
}

func Search(c *gin.Context) {
	query := c.Query("query")
	var pengeluaranLain []models.PengeluaranLain

	query = strings.ToLower(query)
	result := models.DB.Where("LOWER(deskripsi) LIKE ? OR LOWER(metode) LIKE ? OR harga LIKE ? OR tanggal_pengeluaran LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	if err := result.Find(&pengeluaranLain).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "Pengeluaran Lain Not Found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"pengeluaran_lain": pengeluaranLain})
}

func Update(c *gin.Context) {
	var pengeluaranLain models.PengeluaranLain
	id := c.Param("id")

	if err := models.DB.First(&pengeluaranLain, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Pengeluaran Lain Not Found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	if err := c.BindJSON(&pengeluaranLain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if pengeluaranLain.Deskripsi == "" || pengeluaranLain.Harga == 0 || pengeluaranLain.Metode == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua inputan terisi!"})
		return
	}

	models.DB.Save(&pengeluaranLain)
	c.JSON(http.StatusOK, gin.H{"pengeluaran_lain": pengeluaranLain})
}
