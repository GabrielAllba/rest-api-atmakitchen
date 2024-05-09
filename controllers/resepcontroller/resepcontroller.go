package resepcontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func Create(c *gin.Context) {
    // Parse form data including file upload
    // err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
    // if err != nil {
    //     c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
    //     return
    // }

    // Extract form fields
    resepInstruction := c.PostForm("instruction")
    productIDStr := c.PostForm("product_id")

    // Convert form field values to appropriate types
    
    productID, err := strconv.Atoi(productIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Create a new Product instance
    resep := models.Resep{
        Instruction:     resepInstruction,
        ProductId:  productID, 
    }

    // Save the product to the database
    if err := models.DB.Create(&resep).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resep"})
        return
    }

    models.DB.Preload("Product").Preload("Product.Consignation").Preload("Product.ProductType").First(&resep)

    c.JSON(http.StatusOK, gin.H{"resep": resep})
}


func Update(c *gin.Context) {
    // Extract the product ID from the request parameters
    id := c.Param("id")

    // Convert the ID string to an integer
    resepID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resep ID"})
        return
    }

    // Find the existing product in the database
    var existingResep models.Resep
    if err := models.DB.First(&existingResep, resepID).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "Resep not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }

    
    resepInstruction := c.PostForm("instruction")
    


    // Update the product fields
	
	existingResep.Instruction = resepInstruction
	

    // Save the updated product to the database
    if err := models.DB.Save(&existingResep).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resep"})
        return
    }

    models.DB.First(&existingResep)

    // Respond with the updated product
    c.JSON(http.StatusOK, gin.H{"resep": existingResep})
}



func Index(c *gin.Context) {
	var resep []models.Resep
	models.DB.Find(&resep)
	c.JSON(http.StatusOK, gin.H{"resep": resep})
}

func Show(c *gin.Context) {
    var resep models.Resep
    id := c.Param("id")

    if err := models.DB.Preload("Product").First(&resep, id).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
            return
        default:
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }
    
    c.JSON(http.StatusOK, gin.H{"resep": resep})
}

func Search(c *gin.Context) {
    query := c.Query("query")
    var resep []models.Resep

    // Convert query to lowercase
    query = strings.ToLower(query)

    result := models.DB.Joins("JOIN products ON reseps.product_id = products.id").
        Where("LOWER(reseps.instruction) LIKE ? OR LOWER(products.name) LIKE ?", "%"+query+"%", "%"+query+"%").Preload("Product")

    // Execute the query
    if err := result.Find(&resep).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "No recipes found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }

    // Return the search results
    c.JSON(http.StatusOK, gin.H{"recipes": resep})
}
func Delete(c *gin.Context) {
    // Extract the product ID from the request parameters
    id := c.Param("id")

    // Convert the ID string to an integer
    resepID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resep ID"})
        return
    }

    // Find the product in the database
    var resep models.Resep
    if err := models.DB.First(&resep, resepID).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "Resep not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }



    // Delete the product from the database
    if err := models.DB.Delete(&resep).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resep"})
        return
    }

    // Respond with a success message
    c.JSON(http.StatusOK, gin.H{"message": "Resep deleted successfully"})
}