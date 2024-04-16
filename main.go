package main

import (
	"backend-atmakitchen/initializers"
	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	models.ConnectDatabase()
}

func main() {
	r := gin.Default()
	r.Run()

}