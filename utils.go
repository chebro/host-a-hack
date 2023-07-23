package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func WriteStringToFile(filename, content string) error {
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

func CopyFolderToContainer(c *fiber.Ctx, savePath string) error {
	store, err := Store.Get(c)
	if err != nil {
		return err
	}
	// copy zip file to container
	container_id := store.Get("container_id").(string)
	containerPath := fmt.Sprintf("%s:/root", container_id)
	err = exec.Command("docker", "cp", savePath, containerPath).Run()
	if err != nil {
		return err
	}
	// unzip the file inside the container
	filePath := fmt.Sprintf("/root/%s.zip", store.ID())
	err = exec.Command("docker", "exec", container_id, "unzip", filePath).Run()
	if err != nil {
		return err
	}
	// delete the zip file from container
	err = exec.Command("docker", "exec", container_id, "rm", filePath).Run()
	if err != nil {
		return err
	}
	// delete zip file from host
	err = os.Remove(savePath)
	if err != nil {
		return err
	}
	return nil
}

func SaveWebLinkConfig(link string, ip string, port string, filename string) string {
	webLink := fmt.Sprintf("%s.hostahack.xyz", link)
	config := fmt.Sprintf(`
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
                proxy_pass              http://%s:%s;
                proxy_read_timeout      90;
	}
}
    `, webLink, webLink, webLink, ip, port)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(config); err != nil {
		panic(err)
	}
	return webLink
}
