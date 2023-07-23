package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func ReloadNginx() {
	err := exec.Command("sudo", "systemctl", "reload", "nginx").Run()
	if err != nil {
		panic(err)
	}
}

func (container *ContainerInfo) NginxConfPath() string {
	return "/etc/nginx/sites-available/location_" + container.id + ".conf"
}

func (container *ContainerInfo) WebLinkConfPath() string {
	return "/etc/nginx/sites-available/" + container.id + ".hostahack.xyz"
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
	WriteStringToFile(container.NginxConfPath(), defaultConfig)
}

func (container *ContainerInfo) GenerateWebLinkConfig() error {
	if _, err := os.Stat(container.WebLinkConfPath()); err == nil {
		// if file exists, delete the file
		err := os.Remove(container.WebLinkConfPath())
		if err != nil {
			return err
		}
	} else if len(container.open_ports) == 0 {
		// if file doesn't exist and there are no open ports, do nothing
		return nil
	} else {
		// create a new config file
		os.Create(container.WebLinkConfPath())
		// for each port append the config to a file
		for _, port := range container.open_ports {
			link := sha256.Sum256([]byte(container.id + strconv.Itoa(port)))
			short := hex.EncodeToString(link[:5])
			SaveWebLinkConfig(short, container.ip, strconv.Itoa(port), container.WebLinkConfPath())
		}
		return nil
	}
	return nil
}
