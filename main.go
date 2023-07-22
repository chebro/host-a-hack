package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

var pool ContainerPool

func main() {
	pool := NewContainerPool(3)
	defer pool.DisposeContainerPool()

	fmt.Print(pool.size)

	app := fiber.New()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan

		app.Shutdown()
	}()

	app.Static("/", "./public")
  app.Route("/", SetupRoutes)
	app.Listen(":3000")
}
