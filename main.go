package main

import (
	// built-in
	"log"
	"net/http"
	"net/smtp"
	"os"

	// "text/template"

	// 3rd party packages
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Global Variables
var auth smtp.Auth
var emailTo string
var emailFrom string
var emailPass string
var emailHost string
var emailPort string
var emailAddress string

type email struct {
	CC string `json:"cc"`
	Name string `json:"name"`
	Message string `json:"message"`
}

type response struct {
	email
	Success string `json:"success"`
}

func helloWorld(c* gin.Context){
	c.String(http.StatusOK, "hello world")
}

func handleEmail(c* gin.Context) {
	var newEmail email
	
	if err := c.BindJSON(&newEmail); err != nil {
		return
	}

	// load email heading
	// load email body
	
	// address, auth, from, to, message
	if err := smtp.SendMail(
		emailAddress,
		auth,
		emailFrom,
		[]string{emailTo},
		[]byte("Hello World")); err != nil{
			// handle failing to send email message
			log.Fatal(err)
			c.IndentedJSON(http.StatusBadRequest, response{
				email: newEmail, 
				Success: "failed"})
		}

	// return original body
	c.IndentedJSON(http.StatusAccepted, response{
		email: newEmail,
		Success: "success"})
}

func main(){
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	
	// Configure environment variables / global variables
	emailTo = os.Getenv("TO_ADDR")
	emailFrom = os.Getenv("FROM_ADDR")
	emailPass = os.Getenv("FROM_PASS")
	emailHost = os.Getenv("SMTP_HOST")
	emailPort = os.Getenv("SMTP_PORT")
	emailAddress = emailHost + ":" + emailPort
	
	// Authenticate server
	auth = smtp.PlainAuth("", emailFrom, emailPass, emailHost)
	
	// Gin configuration
	router := gin.Default()
	
	// // router.Use(cors.New(cors.Config{
	// // 	AllowOrigins: []string{os.Getenv("ALLOWED_ORIGIN"),"http://"},
	// // 	AllowMethods: []string{"POST"},
	// // 	AllowOriginFunc: func(origin string) bool {
    // //         return origin == "https://github.com"
    // //     },
	// // }))

	router.GET("/", helloWorld)
	router.POST("/email", handleEmail)

	router.Run(":" + os.Getenv("PORT"))
}