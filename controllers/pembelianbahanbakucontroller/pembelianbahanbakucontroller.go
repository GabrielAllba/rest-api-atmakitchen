package pembelianBahanBakubakucontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(c *gin.Context) {
	var pembelianBahanBaku models.PembelianBahanBaku

	
	if err := c.BindJSON(&pembelianBahanBaku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	if err := models.DB.Create(&pembelianBahanBaku).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat pembelian bahan baku"})
		return
	}

	
	var bahan models.Bahan
	if err := models.DB.First(&bahan, pembelianBahanBaku.BahanId).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data bahan"})
		return
	}

	
	bahan.Stok += pembelianBahanBaku.Jumlah
	if err := models.DB.Save(&bahan).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate stok bahan"})
		return
	}

	
	if err := models.DB.Preload("Bahan").First(&pembelianBahanBaku, pembelianBahanBaku.Id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data bahan"})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"pembelian_bahan_baku": pembelianBahanBaku})
}


func Show(c *gin.Context) {
    var pembelian_bahan_baku models.PembelianBahanBaku
    id := c.Param("id")

    if err := models.DB.Preload("Bahan").First(&pembelian_bahan_baku, id).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "pembelian_bahan_baku not found"})
            return
        default:
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }
    
    c.JSON(http.StatusOK, gin.H{"pembelian_bahan_baku": pembelian_bahan_baku})
}

func Index(c *gin.Context) {
	var pembelianBahanBaku []models.PembelianBahanBaku
	if err := models.DB.Preload("Bahan").Find(&pembelianBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pembelian_bahan_baku": pembelianBahanBaku})
}

func GetByDateRange(c *gin.Context) {
	var pembelianBahanBaku []models.PembelianBahanBaku
	
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	searchQuery := c.Query("query")

	
	layout := "2006-01-02"
	from, err := time.Parse(layout, fromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from_date format"})
		return
	}
	to, err := time.Parse(layout, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to_date format"})
		return
	}

	query := models.DB.Where("tanggal_pembelian >= ? AND tanggal_pembelian <= ?", from, to).Preload("Bahan")

	if searchQuery != "" {
		query = query.Joins("JOIN bahans ON bahans.id = pembelian_bahan_bakus.bahan_id").
			Where("bahans.nama LIKE ? OR bahans.stok LIKE ? OR bahans.satuan LIKE ? OR pembelian_bahan_bakus.jumlah LIKE ? OR pembelian_bahan_bakus.keterangan LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if err := query.Find(&pembelianBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pembelian_bahan_baku": pembelianBahanBaku})
	
	
}


func Delete(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var pembelianBahanBaku models.PembelianBahanBaku
	if err := models.DB.First(&pembelianBahanBaku, intId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	
	var bahan models.Bahan
	if err := models.DB.First(&bahan, pembelianBahanBaku.BahanId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load related Bahan"})
		return
	}

	
	bahan.Stok -= pembelianBahanBaku.Jumlah
	if bahan.Stok < 0 {
		bahan.Stok = 0 
		
	}
	if err := models.DB.Save(&bahan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Bahan stock"})
		return
	}

	
	if err := models.DB.Delete(&pembelianBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}


func Update(c *gin.Context) {
	var pembelianBahanBaku models.PembelianBahanBaku
	id := c.Param("id")

	// Bind JSON input to the pembelianBahanBaku model
	if err := c.BindJSON(&pembelianBahanBaku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find and update the existing record
	if err := models.DB.Model(&models.PembelianBahanBaku{}).Where("id = ?", id).Updates(pembelianBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui pembelian bahan baku"})
		return
	}

	// Recalculate the stock for the related Bahan
	var totalPembelian float64
	if err := models.DB.Model(&models.PembelianBahanBaku{}).
		Where("bahan_id = ?", pembelianBahanBaku.BahanId).
		Select("SUM(jumlah)").
		Row().
		Scan(&totalPembelian); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghitung total pembelian bahan"})
		return
	}

	// Update the stock of the related Bahan
	if err := models.DB.Model(&models.Bahan{}).Where("id = ?", pembelianBahanBaku.BahanId).Update("stok", totalPembelian).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui stok bahan"})
		return
	}

	// Preload the related Bahan data and return the updated pembelian_bahan_baku
	if err := models.DB.Preload("Bahan").First(&pembelianBahanBaku, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memuat data bahan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pembelian_bahan_baku": pembelianBahanBaku})
}
