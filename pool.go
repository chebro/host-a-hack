package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type ContainerInfo struct {
	id   string
	port int
  session string
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
	items     []*ContainerInfo
	available []*ContainerInfo
	size      int
	mutex     sync.Mutex
}

func NewContainerPool(size int) *ContainerPool {
	pool := &ContainerPool{
		items:     make([]*ContainerInfo, size),
		available: make([]*ContainerInfo, size),
		size:      size,
	}

	pool.AddContainers(size)

	ReloadNginx()

	return pool
}

func (pool *ContainerPool) DisposeContainerPool() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	args := []string{"stop"}
	for _, container := range pool.items {
		os.Remove(container.NginxConfPath())
		args = append(args, container.id)
	}

	exec.Command("docker", args...).Run()
}

func (pool *ContainerPool) AddContainers(count int) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for i := 1; i <= count; i++ {
		container := CreateContainer()
		pool.items = append(pool.items, container)
		pool.available = append(pool.available, container)
	}
}

func (pool *ContainerPool) GetOne() *ContainerInfo {
	pool.mutex.Lock()

	if len(pool.available) > 0 {
		container := pool.available[0]
		pool.available = pool.available[1:]
		pool.mutex.Unlock()
		go pool.AddContainers(1)
		return container
	}

	pool.mutex.Unlock()
	return nil
}
