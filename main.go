package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mail "github.com/go-mail/mail"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Msg struct {
	Email   string `json:"email" form:"email"`
	Message string `json:"message" form:"message"`
	Name    string `json:"name" form:"name"`
}

func main() {
	app_INIT()
	engine := html.NewFileSystem(http.Dir("./views"), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Post("/emailmsg", func(c *fiber.Ctx) error {
		msg := new(Msg)
		// Binds the request body to the Person struct
		if err := c.BodyParser(&msg); err != nil {
			fmt.Println("ERror in parsing of body", err)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		name, email, message := msg.Name, msg.Email, msg.Message

		sendMail(email, name, message)
		return c.SendStatus(200)
	})

	app.Static("/static", "./static")

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func sendMail(eAddr, name, msg string) error {
	m := mail.NewMessage()
	emailENV := os.Getenv("EMAIL")

	m.SetHeader("From", eAddr)
	m.SetHeader("To", eAddr)
	m.SetHeader("Subject", "Portfolio Contact Form Message from, "+eAddr)
	m.SetBody("text/html", msg+" | "+name)

	d := mail.NewDialer("smtp.gmail.com", 587, emailENV, os.Getenv("EMAILPASS"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func app_INIT() {

}

func GetFHEnvVal(envkey string) string {
	BASE_API_ENV := os.Getenv("BASE_API_ENV")
	return os.Getenv(BASE_API_ENV + envkey)
}
