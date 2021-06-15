package services

import (
	"bytes"
	"echoApp/model"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

func SendMailForNewReviews(recieverEmails []string, reviews []model.AppReview, groupName string, keywords string) {
	// Sender data.
	from := os.Getenv("MAIL_USER")
	password := os.Getenv("MAIL_PASSWORD")

	// Receiver email address.
	to := recieverEmails

	// smtp server configuration.
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")
	auth := smtp.PlainAuth("", from, password, smtpHost)
	emailTemplate, errs := template.ParseFiles("public/views/mails/subscribed_reviews.html")
	if errs != nil {
		log.Printf("template parse : %v",errs)
	}

	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: New Reviews Added For " + groupName +"!\r\n", headers)))

	_ = emailTemplate.Execute(&body, struct {
		Reviews []model.AppReview
		GroupName string
		Keywords string
	}{
		Reviews: reviews,
		GroupName: groupName,
		Keywords: keywords,
	})

	// Sending email.
	err := smtp.SendMail(smtpHost + ":" + smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}