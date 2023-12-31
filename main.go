package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	engine := html.New("views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(cors.New())
	setupRoutes(app)

	connClosed := make(chan bool)

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)
		signal.Notify(sig, os.Kill)
		<-sig

		// TODO: Perform required shutdown procedures

		if err := app.Shutdown(); err != nil {
			log.Println("Unable to kill the application. Error: ", err.Error())
		}
		connClosed <- true
	}()
	serverPort, exists := os.LookupEnv("PORT")
	if !exists {
		serverPort = "8080"
	}
	if err := app.Listen(":" + serverPort); err != nil {
		log.Println("Unable to start the server. Error: ", err.Error())
	}
	<-connClosed
}
