package customerauthcontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context){
	var user models.User

	// string
	user.Email = c.PostForm("email")
	user.Name = c.PostForm("name")
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	bornDateString := c.PostForm("born_date")
	user.PhoneNumber = c.PostForm("phone_number")
	totalPointStr := c.PostForm("total_point")
	roleIdStr := c.PostForm("role_id")

	// validitas semua input terisi
	if user.Email == "" || user.Name == "" || user.Username == "" || user.Password == "" ||
		bornDateString == "" || user.PhoneNumber == "" || totalPointStr == "" || roleIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pastikan semua input terisi"})
		return
	}

	// born date
	
	bornDate, err := time.Parse("2006-01-02", bornDateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal lahir tidak sesuai"})
		return
	}
	user.BornDate = bornDate

	// total point	
	totalPoint, err := strconv.Atoi(totalPointStr)
	if(err != nil){
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total point invalid"})
		return
	}
	user.TotalPoint = totalPoint

	// role id
	roleId, err := strconv.Atoi(roleIdStr)
	if err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format role id tidak benar"})
		return
	}
	// check role id in database
	var role models.Role
	if err := models.DB.Where("id = ?", roleId).First(&role).Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tersebut tidak tersedia"})
		return
	}
	user.RoleId = roleId

	// password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Failed to hash password",
		})
		return
	}

	// existing username and email
	var existingUser models.User
	if err := models.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username tersebut sudah ada"})
		return
	}
	if err := models.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tersebut sudah ada"})
		return
	}


	
	
	user = models.User{Name: user.Name, Email: user.Email, Password: string(hash), Username: user.Username, BornDate: user.BornDate, PhoneNumber: user.PhoneNumber, TotalPoint: user.TotalPoint, RoleId: user.RoleId}
	result := models.DB.Create(&user)

	if result.Error != nil{
		c.JSON(http.StatusBadRequest, gin.H{"user": "Gagal membuat user"})
		return
	}

	var returnUser models.User
	models.DB.Preload("Role").First(&returnUser, "id = ?", user.Id)
	

	c.JSON(http.StatusOK, gin.H{"User": returnUser})
}


func Login(c *gin.Context){
	var req_user models.User

	// Bind other form data fields
	req_user.Email = c.PostForm("email")
	req_user.Password = c.PostForm("password")

	if req_user.Email == "" || req_user.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email dan password tidak boleh kosong"})
		return
	}

	
	var user models.User
	models.DB.Preload("Role").First(&user, "email = ?", req_user.Email)

	if user.Id == 0{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Email tidak ditemukan",
		})

		return
	}

	if user.Role.Name != "Customer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User role bukan Customer",
		})
		return
    }

	err :=	bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req_user.Password))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Invalid password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to create token",
		})

		return
	}

	// send it back

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "","", false, true)
		
	c.JSON(http.StatusOK, gin.H{
		"token":tokenString,
		"user":user,
	})
}
func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func Validate( c *gin.Context){
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message" : user,
	})
}