package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

type ContainerPool struct {
	items []ContainerInfo
	size  int
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
