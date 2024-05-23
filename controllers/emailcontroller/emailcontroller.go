package emailcontroller

// import (
// 	// "log"

// 	// "github.com/gin-gonic/gin"
// 	// "gopkg.in/gomail.v2"

// 	"fmt"

// 	"net/http"
// 	"time"

// 	"net/smtp"

// 	"github.com/dgrijalva/jwt-go"
// )

// const CONFIG_SMTP_HOST = "smtp.gmail.com"
// const CONFIG_SMTP_PORT = 587
// const CONFIG_SENDER_NAME = "gabrielallbasy@gmail.com"
// const CONFIG_AUTH_EMAIL = "gabrielallbasy@gmail.com"
// const CONFIG_AUTH_PASSWORD = "titu yeti vmrw uwcu"

// // func SendEmail(c *gin.Context) {

// // 	emailReceiver := c.PostForm("email_receiver")

// // 	mailer := gomail.NewMessage()
// // 	mailer.SetHeader("From", CONFIG_SENDER_NAME)
// // 	mailer.SetHeader("To", emailReceiver)
// // 	mailer.SetHeader("Subject", "Selamat Anda Menang!")
// // 	mailer.SetBody("text/html", "Hello, <b>have a nice day</b>")

// // 	dialer := gomail.NewDialer(
// // 		CONFIG_SMTP_HOST,
// // 		CONFIG_SMTP_PORT,
// // 		CONFIG_AUTH_EMAIL,
// // 		CONFIG_AUTH_PASSWORD,
// // 	)

// // 	if err := dialer.DialAndSend(mailer); err != nil {
// // 		log.Fatal(err.Error())
// // 		c.JSON(500, gin.H{"error": "Failed to send email"})
// // 		return
// // 	}

// // 	c.JSON(200, gin.H{"message": "Email sent successfully"})
// // }

// func requestPasswordReset(w http.ResponseWriter, r *http.Request) {
// 	email := r.FormValue("email")

// 	var userID int
// 	err := db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	// Create token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     time.Now().Add(time.Hour * 24).Unix(),
// 	})

// 	tokenString, err := token.SignedString([]byte(secretKey))
// 	if err != nil {
// 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Store token in database
// 	_, err = db.Exec("INSERT INTO password_resets (user_id, token, created_at) VALUES (?, ?, ?)", userID, tokenString, time.Now())
// 	if err != nil {
// 		http.Error(w, "Error storing token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Send email

// 	e := email.NewEmail()
// 	e.From = "no-reply@example.com"
// 	e.To = []string{email}
// 	e.Subject = "Password Reset Request"
// 	e.Text = []byte(fmt.Sprintf("Please click the following link to reset your password: http://localhost:3000/reset-password?token=%s", tokenString))
// 	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "your_email@gmail.com", "your_email_password", "smtp.gmail.com"))
// 	if err != nil {
// 		http.Error(w, "Error sending email", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte("Password reset email sent"))
// }

// func resetPassword(w http.ResponseWriter, r *http.Request) {
// 	tokenString := r.URL.Query().Get("token")
// 	newPassword := r.FormValue("password")

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(secretKey), nil
// 	})

// 	if err != nil {
// 		http.Error(w, "Invalid token", http.StatusBadRequest)
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		userID := int(claims["user_id"].(float64))

// 		_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", newPassword, userID)
// 		if err != nil {
// 			http.Error(w, "Error updating password", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Write([]byte("Password updated successfully"))
// 	} else {
// 		http.Error(w, "Invalid token", http.StatusBadRequest)
// 	}
// }
