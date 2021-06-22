package services

import (
	"bytes"
	"echoApp/conf"
	"echoApp/model"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"
)

func SendMailForNewReviews(receiverEmails []string, reviews []model.AppReview, groupName string, keywords string, mailConfig conf.MailConfig) {
	fmt.Printf("\nSending mail for %s Group to : %s", groupName, strings.Join(receiverEmails, " & "))

	if mailConfig.SendMail == "N" {
		fmt.Println("Skipping sending mail")
		return
	}

	// SMTP server configuration.
	from := mailConfig.User
	password := mailConfig.Password
	smtpHost := mailConfig.Host
	smtpPort := mailConfig.Port

	auth := smtp.PlainAuth("", from, password, smtpHost)

	emailTemplate, errs := template.ParseFiles("public/views/mails/subscribed_reviews.html")
	if errs != nil {
		log.Printf("template parse : %v",errs)
	}

	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: New Reviews Added For Group " + groupName +"!\r\n", headers)))

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
	err := smtp.SendMail(smtpHost + ":" + smtpPort, auth, from, receiverEmails, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}