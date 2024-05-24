package hamperscontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var hampers []models.Hampers
	models.DB.Find(&hampers)
	c.JSON(http.StatusOK, gin.H{"hampers": hampers})
}

func GetLatestHampersID(c *gin.Context) {
    var latestHampersID int

    // Query the database to get the latest hampers ID
    if err := models.DB.Model(&models.Hampers{}).Select("id").Order("id DESC").Limit(1).Scan(&latestHampersID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest hampers ID"})
        return
    }

    // Return the latest hampers ID
    c.JSON(http.StatusOK, gin.H{"latest_hampers_id": latestHampersID})
}



func Create(c *gin.Context) {
    id := c.PostForm("id")
    fixId, err := strconv.ParseInt(id, 10, 32)
    hampersNameStr := c.PostForm("hampers_name")
    dailyQuotaStr := c.PostForm("daily_quota")
    deskripsiStr := c.PostForm("deskripsi")
    priceStr := c.PostForm("price")
    stockStr := c.PostForm("stock")
    

    // Convert form field values to appropriate types
    hampersPrice, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
        return
    }
    dailyQuota, err := strconv.ParseFloat(dailyQuotaStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid daily quota value"})
        return
    }

    stock, err := strconv.ParseFloat(stockStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid daily quota value"})
        return
    }
    
    

	// Handle file upload for photo
    file, fileHeader, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Generate a unique filename for the uploaded file
	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	filePath := filepath.Join("images", filename) // Assuming "images" folder exists in your project directory
	fullFilePath := filepath.Join("/", filePath)  // Full file path including the leading "/"

	// Save the uploaded file
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fixFullFilePath := strings.Replace(fullFilePath, "\\", "/", -1)

    

    // Create a new Product instance
    hampers := models.Hampers{
        Id: int(fixId),
        HampersName: hampersNameStr,
        Price: hampersPrice,
        Deskripsi: deskripsiStr,
        Stock: stock,
        DailyQuota: dailyQuota,
        Photo: fixFullFilePath,
    }

    // Save the hampers to the database
    if err := models.DB.Create(&hampers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hampers"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"hampers": hampers})
}

func Update (c *gin.Context){
    id := c.Param("id")
    fixId, err := strconv.ParseInt(id, 10, 32)
    hampersNameStr := c.PostForm("hampers_name")
    dailyQuotaStr := c.PostForm("daily_quota")
    deskripsiStr := c.PostForm("deskripsi")
    priceStr := c.PostForm("price")
    stockStr := c.PostForm("stock")
    

    // Convert form field values to appropriate types
    hampersPrice, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
        return
    }
    dailyQuota, err := strconv.ParseFloat(dailyQuotaStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid daily quota value"})
        return
    }

    stock, err := strconv.ParseFloat(stockStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid daily quota value"})
        return
    }
    
     // Find the existing product in the database
    var existingHampers models.Hampers
    if err := models.DB.First(&existingHampers, id).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "Hampers not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }
    

	// Handle file upload for photo
    file, fileHeader, err := c.Request.FormFile("photo")
	filePathFix := "" // Declare filePathFix outside the if block


	if fileHeader != nil {    
    	defer file.Close()    
    	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
    	filePath := filepath.Join("images", filename) 
    	fullFilePath := filepath.Join("/", filePath) 
    	filePathFix = strings.Replace(fullFilePath, "\\", "/", -1)
    	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
	        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        	return
    	}
	} else {
    	filePathFix = existingHampers.Photo
	}

    

    // Create a new Product instance
    hampers := models.Hampers{
        Id: int(fixId),
        HampersName: hampersNameStr,
        Price: hampersPrice,
        Deskripsi: deskripsiStr,
        Stock: stock,
        DailyQuota: dailyQuota,
        Photo: filePathFix,
    }

    // Save the updated hampers to the database
    if err := models.DB.Save(&hampers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hampers"})
        return
    }

    // Respond with the updated product
    c.JSON(http.StatusOK, gin.H{"hampers": hampers})
}

func CreateDetail(c *gin.Context){
    hampersId := c.Param("id")
    fixHampersId, err := strconv.ParseInt(hampersId, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hampers id value"})
        return
    }
    productIdStr := c.PostForm("product_id")
    quantityStr := c.PostForm("jumlah")

    fixProductId, err := strconv.ParseInt(productIdStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id value"})
        return
    }
    fixQuantity, err := strconv.ParseInt(quantityStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity value"})
        return
    }

    detailHampers := models.DetailHampers{
        HampersId: int(fixHampersId),
        ProductId: int(fixProductId),
        Jumlah:    int(fixQuantity),
    }

    // Save the hampers to the database
    if err := models.DB.Create(&detailHampers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create detail hampers"})
        return
    }

    // Preload Hampers, Product, and ProductType
    if err := models.DB.Preload("Hampers").Preload("Product").Preload("Product.ProductType").First(&detailHampers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload associated models"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"detail_hampers": detailHampers})
}

func Search(c *gin.Context) {
	query := c.Query("query")
	var hampers []models.Hampers

	
	query = strings.ToLower(query)
	result := models.DB.Where("LOWER(hampers_name) LIKE ? OR LOWER(price) LIKE ? OR LOWER(deskripsi) LIKE ? OR LOWER(stock) LIKE ? OR LOWER(daily_quota) LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", query+"%", "%"+query+"%")
	

	
	if err := result.Find(&hampers).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "No hampers found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"hampers": hampers})
}

func Delete(c *gin.Context) {
    id := c.Param("id")
    hampersID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hampers ID"})
        return
    }

    
    if err := models.DB.Where("hampers_id = ?", hampersID).Delete(&models.DetailHampers{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete detail hampers"})
        return
    }

    
    if err := models.DB.Delete(&models.Hampers{}, hampersID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hampers"})
        return
    }

    
    c.JSON(http.StatusOK, gin.H{"message": "Hampers deleted successfully"})
}

func Show(c *gin.Context) {
    var hampers models.Hampers
    id := c.Param("id")

    if err := models.DB.First(&hampers, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Hampers not found"})
            return
        } else {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }

    var detailHampers []models.DetailHampers
    if err := models.DB.Where("hampers_id = ?", hampers.Id).Preload("Product").Preload("Product.ProductType").Preload("Product.Consignation").Find(&detailHampers).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
        return
    }

    response := gin.H{
        "hampers": gin.H{
            "id":           hampers.Id,
            "hampers_name":         hampers.HampersName,
            "deskripsi":  hampers.Deskripsi,
            "daily_quota":   hampers.DailyQuota,
            "photo":   hampers.Photo,
            "price":   hampers.Price,
            "stock":   hampers.Stock,
            "produk_hampers": detailHampers,
        },
    }

    c.JSON(http.StatusOK, response)
}


func DeleteDetailHampers(c *gin.Context){
    id := c.Param("id")
    detailHampersId, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hampers ID"})
        return
    }

    
    if err := models.DB.Delete(&models.DetailHampers{}, detailHampersId).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete detail hampers"})
        return
    }

    
    c.JSON(http.StatusOK, gin.H{"message": "Detail hampers deleted successfully"})
}

func UpdateDetail(c *gin.Context) {
    hampersId := c.Param("id")
    productIdStr := c.PostForm("product_id")
    quantityStr := c.PostForm("jumlah")

    fixHampersId, err := strconv.ParseInt(hampersId, 10, 32)
    if(err != nil){
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
    }
    fixProductId, err := strconv.ParseInt(productIdStr, 10, 32)
    fixQuantity, err := strconv.ParseInt(quantityStr, 10, 32)
    

    // Find the existing detail hampers entry
    var detailHampers models.DetailHampers
     if err := models.DB.Where("hampers_id = ? AND product_id = ?", int(fixHampersId), int(fixProductId)).First(&detailHampers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query detail hampers"})
        return
    }

    // Update the quantity
    detailHampers.Jumlah = int(fixQuantity)

    // Save the updated detail hampers to the database
    if err := models.DB.Save(&detailHampers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update detail hampers"})
        return
    }

    // Preload Hampers, Product, and ProductType
    if err := models.DB.Preload("Hampers").Preload("Product").Preload("Product.ProductType").First(&detailHampers).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload associated models"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"detail_hampers": detailHampers})
}
