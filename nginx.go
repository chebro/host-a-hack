package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

func ReloadNginx() {
	cmd := exec.Command("sudo", "systemctl", "reload", "nginx")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (container *ContainerInfo) NginxConfPath() string {
	return "/etc/nginx/sites-available/location_" + container.id + ".conf"
}

func (container *ContainerInfo) WebLinkConfPath() string {
	return "/etc/nginx/sites-enabled/" + container.id + ".hostahack.xyz"
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
	proxy_pass              http://%s:7681/;
}`, container.id, container.ip)
	WriteStringToFile(container.NginxConfPath(), defaultConfig)
}

func (container *ContainerInfo) SaveWebLinkNginxConf() map[int]string {
	// port -> subdomain link
	portMap := make(map[int]string)

	full_config := ""
	for _, port := range container.open_ports {
		hash := sha256.Sum256([]byte(container.id + strconv.Itoa(port)))
		short := hex.EncodeToString(hash[:5])
		subdomain := fmt.Sprintf("%s.hostahack.xyz", short)

		full_config += fmt.Sprintf(`
	server {
		listen          80;
		server_name     %s;
		return 301      https://%s$request_uri;
	}
	server {
	        listen 443;
	        server_name                     %s;
	        ssl_certificate                 /etc/letsencrypt/live/hostahack.xyz/fullchain.pem;
	        ssl_certificate_key             /etc/letsencrypt/live/hostahack.xyz/privkey.pem;
	        ssl                             on;

	        ssl_session_cache               builtin:1000    shared:SSL:10m;
	        ssl_protocols                   TLSv1.2         TLSv1.3;
	        ssl_ciphers                     HIGH:!aNULL:!eNULL:!EXPORT:!CAMELLIA:!DES:!MD5:!PSK:!RC4;
	        ssl_prefer_server_ciphers       on;

	        location / {
	                proxy_http_version      1.1;
	                proxy_set_header        Host $host;
	                proxy_set_header        X-Real-IP $remote_addr;
	                proxy_set_header        X-Forwarded-Proto $scheme;
	                proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
	                proxy_set_header        Upgrade $http_upgrade;
	                proxy_set_header        Connection "upgrade";
	                proxy_pass              http://%s:%d;
	                proxy_read_timeout      90;
		}
	}
	    `, subdomain, subdomain, subdomain, container.ip, port)
		portMap[port] = subdomain
	}

	WriteStringToFile(container.WebLinkConfPath(), full_config)
	ReloadNginx()
	return portMap
}
