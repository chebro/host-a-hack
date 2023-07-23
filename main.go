package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

var pool = NewContainerPool(3)

func main() {
	defer pool.DisposeContainerPool()

	app := fiber.New()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan

		app.Shutdown()
	}()

	app.Static("/", "./public")
	app.Route("/", SetupRoutes)
	err := app.Listen(":3000")
	if err != nil {
		fmt.Println(err)
	}
}
