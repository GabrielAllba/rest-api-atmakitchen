package emailcontroller

import (
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "gabrielallbasy@gmail.com"
const CONFIG_AUTH_EMAIL = "gabrielallbasy@gmail.com"
const CONFIG_AUTH_PASSWORD = "titu yeti vmrw uwcu"

func SendEmail(c *gin.Context) {

    emailReceiver := c.PostForm("email_receiver")

    mailer := gomail.NewMessage()
    mailer.SetHeader("From", CONFIG_SENDER_NAME)
    mailer.SetHeader("To", emailReceiver)
    mailer.SetHeader("Subject", "Selamat Anda Menang!")
    mailer.SetBody("text/html", "Hello, <b>have a nice day</b>")

    dialer := gomail.NewDialer(
        CONFIG_SMTP_HOST,
        CONFIG_SMTP_PORT,
        CONFIG_AUTH_EMAIL,
        CONFIG_AUTH_PASSWORD,
    )

    if err := dialer.DialAndSend(mailer); err != nil {
        log.Fatal(err.Error())
        c.JSON(500, gin.H{"error": "Failed to send email"})
        return
    }

    c.JSON(200, gin.H{"message": "Email sent successfully"})
}
