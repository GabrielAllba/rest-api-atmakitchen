package productcontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
    var product models.Product

    // Bind JSON request body to the Product struct
    if err := c.BindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }

    // Handle file upload for photo
    // file, err := c.FormFile("photo")
    // if err != nil {
    //     c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
    //     return
    // }
    
    // // Generate a unique filename for the uploaded file (you can use UUID or any other method)
    // filename := uuid.New().String() + filepath.Ext(file.Filename)
    // filePath := filepath.Join("images", filename) // Assuming "images" folder exists in your project directory

    // // Save the uploaded file
    // if err := c.SaveUploadedFile(file, filePath); err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
    //     return
    // }
    // product.Photo = filePath

    // If the product type is "Titipan", validate and fetch consignation
    if product.ProductType.Name == "Titipan" {
        var consignation models.Consignation
        if err := models.DB.Where("id = ?", product.ConsignationId).First(&consignation).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Penitip tidak tersedia"})
            return
        }
    }

    // Create the product in the database
    if err := models.DB.Create(&product).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat Product"})
        return
    }

    // Preload associated data and return the created product
    var returnProduct models.Product
    if err := models.DB.Preload("Consignation").Preload("ProductType").First(&returnProduct, product.Id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"product": returnProduct})
}
