package main

import (
	adminauthcontroller "backend-atmakitchen/controllers/admincontroller"
	"backend-atmakitchen/controllers/autologincontroller"
	"backend-atmakitchen/controllers/bahancontroller"
	"backend-atmakitchen/controllers/bankcontroller"
	"backend-atmakitchen/controllers/consignationcontroller"
	customerauthcontroller "backend-atmakitchen/controllers/customercontroller"
	"backend-atmakitchen/controllers/emailcontroller"
	"backend-atmakitchen/controllers/hamperscontroller"
	moauthcontrollerauthcontroller "backend-atmakitchen/controllers/mocontroller"
	ownerauthcontroller "backend-atmakitchen/controllers/ownercontroller"
	"backend-atmakitchen/controllers/productcontroller"
	"backend-atmakitchen/controllers/producttypecontroller"
	"backend-atmakitchen/controllers/resepcontroller"
	"backend-atmakitchen/controllers/rolecontroller"
	"backend-atmakitchen/controllers/tokencontroller"
	"backend-atmakitchen/initializers"
	"backend-atmakitchen/middleware"
	"backend-atmakitchen/models"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	models.ConnectDatabase()
}
func getImage(c *gin.Context) {
	filename := c.Param("filename")
	imagePath := filepath.Join("images", filename)
	c.File(imagePath)
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	// r.Use(func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "*")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
	// 	c.Header("Access-Control-Allow-Credentials", "true")
	// 	if c.Request.Method == "OPTIONS" {
	//         c.AbortWithStatus(204)
	//     	return
	// 	}
	// 	c.Next()
	// })

	// login based by role
	autologin := r.Group("/api/autologin")
	{
		autologin.POST("/login", autologincontroller.Login)
		autologin.POST("/logout", autologincontroller.Logout)
	}

	email := r.Group("/api/email")
	{
		email.POST("/send", emailcontroller.SendEmail)
	}

	// product
	product := r.Group("/api/product")
	{
		product.POST("", productcontroller.Create)
		product.GET("", productcontroller.Index)
		product.GET("/:id", productcontroller.Show)
		product.GET("/search", productcontroller.Search)
		product.DELETE("/:id", productcontroller.Delete)
		product.PUT("/:id", productcontroller.Update)
		product.GET("/type", productcontroller.SearchType)
		product.GET("/type/search", productcontroller.SearchProductByType)
	}

	// customer
	user := r.Group("/api/customer")
	{
		user.POST("/signup", customerauthcontroller.Signup)
		user.POST("/login", customerauthcontroller.Login)
		user.POST("/logout", customerauthcontroller.Logout)
		user.GET("/token/validate/:tokenString", customerauthcontroller.Validate)
		user.GET("/user", customerauthcontroller.GetUser)
		user.GET("/users/:role_id", customerauthcontroller.GetUsersByRoleID)
		user.PUT("/updatepassword/:email", customerauthcontroller.UpdatePassword)

	}

	// admin
	admin := r.Group("/api/admin")
	{
		role := admin.Group("/role")
		{
			role.POST("", rolecontroller.Create)
			role.GET("/:id", rolecontroller.Show)
		}

		admin.POST("/login", adminauthcontroller.Login)
		admin.POST("/logout", adminauthcontroller.Logout)
		admin.GET("/token/validate/:tokenString", adminauthcontroller.Validate)
	}

	// mo
	mo := r.Group("/api/mo")
	{
		mo.POST("/login", moauthcontrollerauthcontroller.Login)
		mo.POST("/logout", moauthcontrollerauthcontroller.Logout)
		mo.GET("/token/validate/:tokenString", moauthcontrollerauthcontroller.Validate)
	}

	// owner
	owner := r.Group("/api/owner")
	{
		owner.POST("/login", ownerauthcontroller.Login)
		owner.POST("/logout", ownerauthcontroller.Logout)
		owner.GET("/token/validate/:tokenString", ownerauthcontroller.Validate)
	}

	// product type
	product_type := r.Group("/api/product_type")
	{
		product_type.POST("", producttypecontroller.Create)
		product_type.GET("/", producttypecontroller.Index)
		product_type.GET("/:id", producttypecontroller.Show)
	}

	// consignation
	consignation := r.Group("/api/consignation")
	{
		consignation.POST("", consignationcontroller.Create)
		consignation.GET("", consignationcontroller.Index)
		consignation.GET("/search", consignationcontroller.Search)
		consignation.DELETE("/:id", consignationcontroller.Delete)
		consignation.PUT("/:id", consignationcontroller.Update)

	}

	bank := r.Group("/api/bank")
	{
		bank.POST("", bankcontroller.Create)
		bank.GET("", bankcontroller.Index)

	}

	hampers := r.Group("/api/hampers")
	{
		hampers.GET("", hamperscontroller.Index)
		hampers.GET("/:id", hamperscontroller.Show)
		hampers.PUT("/:id", hamperscontroller.Update)
		hampers.GET("/latest_id", hamperscontroller.GetLatestHampersID)
		hampers.POST("", hamperscontroller.Create)
		hampers.POST("/detail/:id", hamperscontroller.CreateDetail)
		hampers.GET("/search", hamperscontroller.Search)
		hampers.DELETE("/:id", hamperscontroller.Delete)
		hampers.DELETE("/detail/:id", hamperscontroller.DeleteDetailHampers)

	}

	//resep
	resep := r.Group("/api/resep")
	{
		resep.POST("", resepcontroller.Create)
		resep.GET("", resepcontroller.Index)
		resep.GET("/:id", resepcontroller.Show)
		resep.GET("/search", resepcontroller.Search)
		resep.DELETE("/:id", resepcontroller.Delete)
		resep.PUT("/:id", resepcontroller.Update)
		// resep.GET("/type", resepcontroller.SearchType);
		// resep.GET("/type/search", resepcontroller.SearchProductByType);
	}

	//bahan
	bahan := r.Group("/api/bahan")
	{
		bahan.POST("", bahancontroller.Create)
		bahan.GET("", bahancontroller.Index)
		bahan.GET("/:id", bahancontroller.Show)
		bahan.GET("/search", bahancontroller.Search)
		bahan.DELETE("/:id", bahancontroller.Delete)
		bahan.PUT("/:id", bahancontroller.Update)
	}

	roles := r.Group("/api/roles")
	{
		roles.POST("", rolecontroller.Create)
		roles.GET("", rolecontroller.Index)
		roles.PUT("/:id", rolecontroller.Update)
		roles.GET("/:id", rolecontroller.Show)
		roles.DELETE("/:id", rolecontroller.Delete)
	}


	
	token := r.Group("/api/token")
	{
		token.POST("/create/:user_id", tokencontroller.CreateToken)
		token.DELETE("/delete", tokencontroller.DeleteToken)
		token.GET("/check", tokencontroller.CheckToken)
	}

	// Define the route
	r.GET("/api/images/:filename", getImage)

	// check auth middleware
	r.GET("/api/validates", middleware.RequireAuth, customerauthcontroller.Validate)

	r.Run("127.0.0.1:8000")

}
