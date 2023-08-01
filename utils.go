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
		logger.Println(err)
		return err
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(content)
	if err != nil {
		logger.Println(err)
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
