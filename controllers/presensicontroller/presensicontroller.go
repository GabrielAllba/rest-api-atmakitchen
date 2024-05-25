package presensicontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(c *gin.Context) {
	var presensi models.Presensi

	if err := c.BindJSON(&presensi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Create(&presensi).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat presensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"presensi": presensi})
}

func Index(c *gin.Context) {
	var presensi []models.Presensi
	models.DB.Find(&presensi)
	c.JSON(http.StatusOK, gin.H{"presensi": presensi})
}

func Update(c *gin.Context) {
	var presensi models.Presensi
    id := c.Param("id")
		if err := models.DB.First(&presensi, id).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "presensi not found"})
				return
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
				return
			}
		}
	
		if err := c.BindJSON(&presensi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		models.DB.Save(&presensi)
		c.JSON(http.StatusOK, gin.H{"presensi": presensi})
	}

// func Show(c *gin.Context) {
// 	var bahan models.Bahan
// 	id := c.Param("id")
// 	if err := models.DB.First(&bahan, id).Error; err != nil {
// 		switch err {
// 		case gorm.ErrRecordNotFound:
// 			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Bahan not found"})
// 			return
// 		default:
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
// }

// func Delete(c *gin.Context) {
// 	var bahan models.Bahan
// 	id := c.Param("id")
// 	if err := models.DB.First(&bahan, id).Error; err != nil {
// 		switch err {
// 		case gorm.ErrRecordNotFound:
// 			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Bahan not found"})
// 			return
// 		default:
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
// 			return
// 		}
// 	}
// 	models.DB.Delete(&bahan)
// 	c.JSON(http.StatusOK, gin.H{"message": "Bahan deleted"})
// }

// func Search(c *gin.Context) {
// 	query := c.Query("query")
// 	var bahans []models.Bahan

	
// 	query = strings.ToLower(query)
// 	result := models.DB.Where("LOWER(nama) LIKE ?  OR stok LIKE ? OR satuan LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	

	
// 	if err := result.Find(&bahans).Error; err != nil {
// 		switch err {
// 		case gorm.ErrRecordNotFound:
// 			c.JSON(http.StatusNotFound, gin.H{"message": "No bahan found"})
// 			return
// 		default:
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"bahans": bahans})
// }

// func Update(c *gin.Context) {
// 	var bahan models.Bahan
// 	id := c.Param("id")
// 	if err := models.DB.First(&bahan, id).Error; err != nil {
// 		switch err {
// 		case gorm.ErrRecordNotFound:
// 			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Bahan not found"})
// 			return
// 		default:
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
// 			return
// 		}
// 	}

// 	if err := c.BindJSON(&bahan); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if bahan.Nama == ""  || bahan.Satuan == ""  {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
// 		return
// 	}

// 	models.DB.Save(&bahan)
// 	c.JSON(http.StatusOK, gin.H{"bahan": bahan})
// }

