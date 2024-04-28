package productcontroller

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

func Create(c *gin.Context) {
    // Parse form data including file upload
    err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
        return
    }

    // Extract form fields
    productName := c.PostForm("name")
    productPriceStr := c.PostForm("price")
    productDescription := c.PostForm("description")
    productStockStr := c.PostForm("stock")
    productDailyQuotaStr := c.PostForm("daily_quota")
    productRewardPoinStr := c.PostForm("reward_poin")
    productStatus := c.PostForm("status")
    productTypeIDStr := c.PostForm("product_type_id")
    consignationIDStr := c.PostForm("consignation_id")

    // Convert form field values to appropriate types
    productPrice, err := strconv.ParseFloat(productPriceStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
        return
    }
    productStock, err := strconv.ParseFloat(productStockStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock value"})
        return
    }
    productDailyQuota, err := strconv.ParseFloat(productDailyQuotaStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid daily quota value"})
        return
    }
    productRewardPoin, err := strconv.Atoi(productRewardPoinStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward poin value"})
        return
    }
    productTypeID, err := strconv.Atoi(productTypeIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product type ID"})
        return
    }

  

    // Initialize consignationID as nil
    var consignationID *int

    // Parse consignation ID if provided
    if consignationIDStr != "null" {
        id, err := strconv.Atoi(consignationIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consignation ID"})
            return
        }
        consignationID = &id
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
    product := models.Product{
        Name:           productName,
        Price:          productPrice,
        Description:    productDescription,
        Photo:          fixFullFilePath,
        Stock:          productStock,
        DailyQuota:     productDailyQuota,
        RewardPoin:     productRewardPoin,
        Status:         productStatus,
        ProductTypeId:  productTypeID,
        ConsignationId: consignationID, 
    }

    // Save the product to the database
    if err := models.DB.Create(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"product": product})
}

func Index(c *gin.Context) {
	var product []models.Product
	models.DB.Find(&product)
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func Show(c *gin.Context) {
    var product models.Product
    id := c.Param("id")

    if err := models.DB.Preload("Consignation").Preload("ProductType").First(&product, id).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
            return
        default:
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }
    
    c.JSON(http.StatusOK, gin.H{"product": product})
}


func Search(c *gin.Context) {
	query := c.Query("query")
	var products []models.Product

	
	query = strings.ToLower(query)
	result := models.DB.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR price LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	result = result.Preload("ProductType")

	
	if err := result.Find(&products).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "No products found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func Delete(c *gin.Context) {
    // Extract the product ID from the request parameters
    id := c.Param("id")

    // Convert the ID string to an integer
    productID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Find the product in the database
    var product models.Product
    if err := models.DB.First(&product, productID).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }



    // Delete the product from the database
    if err := models.DB.Delete(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
        return
    }

    // Respond with a success message
    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func Update(c *gin.Context) {
    // Extract the product ID from the request parameters
    id := c.Param("id")

    // Convert the ID string to an integer
    productID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Find the existing product in the database
    var existingProduct models.Product
    if err := models.DB.First(&existingProduct, productID).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }

    // Parse form data including file upload
    err = c.Request.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
        return
    }

    // Extract form fields
    productName := c.PostForm("name")
    productPriceStr := c.PostForm("price")
    productDescription := c.PostForm("description")
    productStockStr := c.PostForm("stock")
    productDailyQuotaStr := c.PostForm("daily_quota")
    productRewardPoinStr := c.PostForm("reward_poin")
    productStatus := c.PostForm("status")
    productTypeIDStr := c.PostForm("product_type_id")
    consignationIDStr := c.PostForm("consignation_id")

    // Convert form field values to appropriate types
    productPrice, err := strconv.ParseFloat(productPriceStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
        return
    }
    productStock, err := strconv.ParseFloat(productStockStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock value"})
        return
    }
    productDailyQuota, err := strconv.ParseFloat(productDailyQuotaStr, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid daily quota value"})
        return
    }
    productRewardPoin, err := strconv.Atoi(productRewardPoinStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward poin value"})
        return
    }
    productTypeID, err := strconv.Atoi(productTypeIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product type ID"})
        return
    }

    // Initialize consignationID as nil

    var consignationID *int

    // Parse consignation ID if provided
    if consignationIDStr != "null" {
        id, err := strconv.Atoi(consignationIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consignation ID"})
            return
        }
        consignationID = &id
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
    	filePathFix = existingProduct.Photo
	}


    // Update the product fields
    existingProduct.Name = productName
    existingProduct.Price = productPrice
    existingProduct.Description = productDescription
    existingProduct.Stock = productStock
    existingProduct.DailyQuota = productDailyQuota
    existingProduct.RewardPoin = productRewardPoin
    existingProduct.Status = productStatus
    existingProduct.ProductTypeId = productTypeID
    existingProduct.ConsignationId = consignationID
    existingProduct.Photo = filePathFix // Update the photo path

    // Save the updated product to the database
    if err := models.DB.Save(&existingProduct).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
        return
    }

    // Respond with the updated product
    c.JSON(http.StatusOK, gin.H{"product": existingProduct})
}
