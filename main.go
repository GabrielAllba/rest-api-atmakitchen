package main

import (
	customerauthcontroller "backend-atmakitchen/controllers/customercontroller"
	"backend-atmakitchen/controllers/rolecontroller"
	"backend-atmakitchen/initializers"
	"backend-atmakitchen/middleware"
	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	models.ConnectDatabase()
}

func main() {
	r := gin.Default()

	// customer
	user := r.Group("/api/customer")
	{
		user.POST("/signup", customerauthcontroller.Signup)
		user.POST("/login", customerauthcontroller.Login)
		user.POST("/logout", customerauthcontroller.Logout)
		
	}

	// admin
	admin := r.Group("/api/admin")
	{
		role := admin.Group("/role")
		{
			role.POST("/", rolecontroller.Create);
		}
	}

	// check auth middleware
	r.GET("/api/validate", middleware.RequireAuth, customerauthcontroller.Validate)


	r.Run()

}