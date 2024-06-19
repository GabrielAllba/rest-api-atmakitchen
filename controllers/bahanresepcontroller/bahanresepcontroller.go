package bahanresepcontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func Index(c *gin.Context) {
	var bahan_resep []models.BahanResep
	models.DB.Find(&bahan_resep).Preload("Product").Preload("Bahan")
	c.JSON(http.StatusOK, gin.H{"bahan_resep": bahan_resep})
}

// func Create(c *gin.Context) {

//     // Extract form fields
//     id := c.PostForm("id")
//     productIDStr := c.PostForm("product_id")
//     bahanIDStr := c.PostForm("bahan_id")
//     quantity := c.PostForm("quantity")
//     unit := c.PostForm("unit")

//     // Convert form field values to appropriate types
    
//     productID, err := strconv.Atoi(productIDStr)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
//         return
//     }

//     bahanID, err := strconv.Atoi(bahanIDStr)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bahan ID"})
//         return
//     }

//     idFix, err := strconv.Atoi(id)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
//         return
//     }
//     quantityFix, err := strconv.Atoi(quantity)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Quantity"})
//         return
//     }

//     // Create a new Product instance
//     bahan_resep := models.BahanResep{
//         Id: idFix,
//         BahanId: bahanID,
//         ProductId:  productID,
//         Quantity: float64(quantityFix),
//         Unit: unit,
         
//     }

//     // Save the product to the database
//     if err := models.DB.Create(&bahan_resep).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resep"})
//         return
//     }

//     models.DB.Preload("Product").Preload("Product.Consignation").Preload("Product.ProductType").First(&bahan_resep)

//     c.JSON(http.StatusOK, gin.H{"bahan_resep": bahan_resep})
// }

func Create(c *gin.Context) {
    var bahanResep models.BahanResep

    if err := c.BindJSON(&bahanResep); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Create(&bahanResep).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create bahanResep"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"bahanResep": bahanResep})
}

func Show(c *gin.Context) {
    var bahan_resep []models.BahanResep
    product_id := c.Param("product_id")

    // Use Where to filter by product_id and then load all matching records
    if err := models.DB.Preload("Bahan").Preload("Product").Where("product_id = ?", product_id).Find(&bahan_resep).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "bahan_resep not found"})
            return
        } else {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }
    c.JSON(http.StatusOK, gin.H{"bahan_resep": bahan_resep})
}
