package main

import (
	adminauthcontroller "backend-atmakitchen/controllers/admincontroller"
	customerauthcontroller "backend-atmakitchen/controllers/customercontroller"
	moauthcontrollerauthcontroller "backend-atmakitchen/controllers/mocontroller"
	ownerauthcontroller "backend-atmakitchen/controllers/ownercontroller"
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
	
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
        c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
		c.Next()
	})

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

		admin.POST("/login", adminauthcontroller.Login)
		admin.POST("/logout", adminauthcontroller.Logout)
	}

	// mo
	mo := r.Group("/api/mo")
	{			
		mo.POST("/login", moauthcontrollerauthcontroller.Login)
		mo.POST("/logout", moauthcontrollerauthcontroller.Logout)	
	}


	// owner
	owner := r.Group("/api/owner")
	{			
		owner.POST("/login", ownerauthcontroller.Login)
		owner.POST("/logout", ownerauthcontroller.Logout)	
	}

	// check auth middleware
	r.GET("/api/validates", middleware.RequireAuth, customerauthcontroller.Validate)


	r.Run()

}