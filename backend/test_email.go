package main

import (
    "fmt"
    "log"
    "net/smtp"
    "os"
    "github.com/joho/godotenv"
)

func main() {
    godotenv.Load(".env")
    from := os.Getenv("SMTP_USER")
    pass := os.Getenv("SMTP_PASSWORD")
    to := os.Getenv("CONTACT_EMAIL")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")

    subject := "Test Email from Go"
    body := "This is a test email sent from a Go script."
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, pass, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        log.Fatal("Failed to send email:", err)
    }
    fmt.Println("Test email sent successfully!")
} 