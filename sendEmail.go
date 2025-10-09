package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func sendEmail(to, subject, otp string) error {
	from := os.Getenv("SMTP_EMAIL")        // e.g. soberlyhigh@gmail.com
	password := os.Getenv("SMTP_PASSWORD") // Gmail app password

	body := fmt.Sprintf(`
        <html>
        <body style="font-family: sans-serif; background-color: #111; color: #eee; padding: 20px;">
            <h2>Hello,</h2>
            <p>Your verification code is:</p>
            <h1 style="color: #00ffff;">%s</h1>
            <p>This code will expire in 5 minutes.</p>
            <br>
            <p>Thanks,<br>LocalLoop Team</p>
        </body>
        </html>
    `, otp)

	e := email.NewEmail()
	e.From = fmt.Sprintf("LocalLoop <%s>", from)
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)

	err := e.SendWithTLS("smtp.gmail.com:465", smtp.PlainAuth("", from, password, "smtp.gmail.com"), &tls.Config{
		ServerName: "smtp.gmail.com",
	})
	if err != nil {
		log.Printf("Email send error: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
