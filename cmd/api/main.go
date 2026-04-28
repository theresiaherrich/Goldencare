package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/theresiaherrich/Goldencare/internal/bootstrap"
	"github.com/theresiaherrich/Goldencare/internal/route"
)

func main() {
	log.Println("Starting Goldencare API...")

	log.Println("Initializing dependency container...")
	container, err := bootstrap.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()
	log.Println("Container initialized successfully")

	app := fiber.New(fiber.Config{
		AppName: container.Config.AppName,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	route.RegisterRoutes(app, container)

	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		<-sigch
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	log.Printf("Server starting on port %s in %s mode", container.Config.AppPort, container.Config.AppEnv)
	log.Println("About to call app.Listen()...")
	if err := app.Listen(":" + container.Config.AppPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
