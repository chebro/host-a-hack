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
	items     []ContainerInfo
	available []*ContainerInfo
	size      int
	mutex     sync.Mutex
}

func NewContainerPool(size int) *ContainerPool {
	pool := &ContainerPool{
		items:     make([]ContainerInfo, size),
		available: make([]*ContainerInfo, size),
		size:      size,
	}

	for i := 0; i < size; i++ {
		pool.items[i] = *CreateContainer()
		pool.available[i] = &pool.items[i]
		fmt.Println(pool.items[i].id)
	}

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

func (pool *ContainerPool) GetOne() *ContainerInfo {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.available) > 0 {
		container := pool.available[0]
		pool.available = pool.available[1:]
		return container
	}

	return nil
}
