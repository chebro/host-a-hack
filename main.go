package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

var sockPath = "/tmp/hostahack.sock"
var pool *ContainerPool

func main() {

	app := fiber.New()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan

		app.Shutdown()
	}()

	app.Static("/", "./public")
	app.Route("/", SetupRoutes)

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	if err := os.Chmod(sockPath, 0777); err != nil {
		fmt.Println(err)
		return
	}

	pool = NewContainerPool(3)
	defer pool.DisposeContainerPool()

	if err := app.Listener(listener); err != nil {
		fmt.Println(err)
	}
}
