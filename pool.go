package main

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	id         string
	ip         string
	session    string
	open_ports []int
}

func CreateContainer() *ContainerInfo {
	cmd := exec.Command("docker", "run", "-d", "--rm", "-v", sockPath+":/tmp/server.sock", "hostahack:latest")
	stdout, err := cmd.Output()

	if err != nil {
		logger.Fatalln(err)
	}

	container := &ContainerInfo{
		id: strings.TrimSpace(string(stdout)),
	}

	ctx := context.Background()

	var dockerCli, _ = client.NewClientWithOpts(client.FromEnv)
	cjson, err := dockerCli.ContainerInspect(ctx, container.id)
	if err != nil {
		panic(err)
	}

	container.ip = cjson.NetworkSettings.IPAddress

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
		size: size,
	}

	pool.AddContainers(size)

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

	logger.Println("Stopping containers")
	out, err := exec.Command("docker", args...).Output()
	if err == nil {
		logger.Println(string(out))
	}
}

func (pool *ContainerPool) AddContainers(count int) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for i := 1; i <= count; i++ {
		container := CreateContainer()
		pool.items = append(pool.items, container)
		pool.available = append(pool.available, container)
		logger.Println("Container added " + container.id)
	}

	ReloadNginx()
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

func (pool *ContainerPool) GetContainerById(container_id string) *ContainerInfo {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for _, container := range pool.items {
		if strings.HasPrefix(container.id, container_id) {
			return container
		}
	}

	return nil
}

func (pool *ContainerPool) RemoveContainerById(container_id string) {
	for idx, container := range pool.items {
		if strings.HasPrefix(container.id, container_id) {
			os.Remove(container.NginxConfPath())
			pool.items[idx] = pool.items[len(pool.items)-1]
			pool.items = pool.items[:len(pool.items)-1]
			return
		}
	}
}
