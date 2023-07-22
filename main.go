package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	folderRoutes "github.com/chebro/host-a-hack/internals/routes"
	"github.com/gofiber/fiber/v2"
)

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
	api := app.Group("api")
	folderRoutes.FolderRoutes(api)
	app.Listen(":3000")
}
