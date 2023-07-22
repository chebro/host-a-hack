package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	folderRoutes "github.com/chebro/host-a-hack/internals/routes"
	"github.com/gofiber/fiber/v2"
)

type ContainerInfo struct {
	id   string
	port int
}

func CreateContainer() *ContainerInfo {
	const min_port = 30000
	const max_port = 60000
	port := rand.Intn(max_port-min_port) + min_port
	cmd := exec.Command("docker", "run", "-d", "--rm", "-p", strconv.Itoa(port)+":7681", "tsl0922/ttyd:alpine")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	container := &ContainerInfo{
		id:   strings.TrimSpace(string(stdout)),
		port: port,
	}

	container.SaveNginxConf()

	return container
}

func writeStringToFile(filename, content string) error {
	// Open the file in write-only mode with file creation permission
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (container *ContainerInfo) SaveNginxConf() {
	defaultConfig := fmt.Sprintf(`location /ttyd/%s {
	proxy_http_version      1.1;
	proxy_set_header        Host $host;
	proxy_set_header        X-Real-IP $remote_addr;
	proxy_set_header        X-Forwarded-Proto $scheme;
	proxy_set_header        X-Forwarded-Proto $scheme;
	proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header        Upgrade $http_upgrade;
	proxy_set_header        Connection "upgrade";
	proxy_pass              http://localhost:%s/;
}`, container.id, strconv.Itoa(container.port))
	writeStringToFile(container.NginxConfPath(), defaultConfig)
}

type ContainerPool struct {
	items []ContainerInfo
	size  int
}

func (container *ContainerInfo) NginxConfPath() string {
	return "/etc/nginx/sites-available/location_" + container.id + ".conf"
}

func ReloadNginx() {
	err := exec.Command("sudo", "systemctl", "reload", "nginx").Run()
	if err != nil {
		panic(err)
	}
}

func NewContainerPool(size int) *ContainerPool {
	pool := &ContainerPool{
		items: make([]ContainerInfo, size),
		size:  size,
	}

	for i := 0; i < size; i++ {
		pool.items[i] = *CreateContainer()
		fmt.Println(pool.items[i].id)
	}

	ReloadNginx()

	return pool
}

func (pool *ContainerPool) DisposeContainerPool() {
	args := []string{"stop"}
	for _, item := range pool.items {
		os.Remove(item.NginxConfPath())
		args = append(args, item.id)
	}

	exec.Command("docker", args...).Run()
}

func main() {
	pool := NewContainerPool(3)
	defer pool.DisposeContainerPool()

	fmt.Print(pool.size)

	app := fiber.New()

	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		app.Shutdown()
	}()

  app.Static("/", "./public")
  api := app.Group("api")
  folderRoutes.FolderRoutes(api)
  app.Listen(":3000")
}
