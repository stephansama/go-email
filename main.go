package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"text/template"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Global Variables
var auth smtp.Auth
var emailTo string
var emailFrom string
var emailPass string
var emailHost string
var emailPort string
var emailAddress string

// POST route @/email
// cc to optional email
// accept JSON
// format using HTML template

type email struct {
	CC string `json:"cc"`
	Name string `json:"name"`
	Message string `json:"message"`
}

func sendEmail(c* gin.Context) {
	var newEmail email
	
	if err := c.BindJSON(&newEmail); err != nil {
		return
	}

	// parse CC email
	// validCC,ccError := mail.ParseAddress(newEmail.CC)
	// if ccError != nil { return }

	var emailSubject bytes.Buffer
	subjectFilter, _ := template.New("").Parse("Subject: Email from {.}\r\n")
	_ = subjectFilter.Execute(&emailSubject, newEmail.Name)
	emailMime := "MIME-version: 1.0\r\n"

	newLine := "\r\n"

	emailHeading := []byte(
		"From: Stephan Randle <" + emailFrom + ">" + newLine +
		"To: Stephan Randle <" + emailTo + ">" + newLine +
		"CC: " + newEmail.CC + newLine +
		"Subject: " + newEmail.Name + newLine)

	emailMessage := []byte(emailSubject.String() + emailMime + newEmail.Message)

	
	smtp.SendMail(
		emailAddress,
		auth,
		emailFrom,
		[]string{emailFrom},
		emailMessage)
		
	fmt.Println(emailHeading)

	// return original body
	c.IndentedJSON(http.StatusAccepted, newEmail)
}

func main(){
	router := gin.Default()

	// authenticate server
	auth = smtp.PlainAuth("", emailFrom, emailPass, emailHost)
	// load / configure environment variables
	emailTo = os.Getenv("TO_ADDR")
	emailFrom = os.Getenv("FROM_ADDR")
	emailPass = os.Getenv("FROM_PASS")
	emailHost = os.Getenv("SMTP_HOST")
	emailPort = os.Getenv("SMTP_PORT")
	emailAddress = emailHost + ":" + emailPort

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowMethods: []string{"POST"},
	}))

	router.POST("/email", sendEmail)

	router.Run(os.Getenv("PORT"))
}