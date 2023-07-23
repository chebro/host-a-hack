package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetLanding(c *fiber.Ctx) error {
	store, err := Store.Get(c)
	if err != nil {
		panic(err)
	}
	return c.JSON(fiber.Map{"status": "success", "id": store.ID()})
}

func PostFolder(c *fiber.Ctx) error {
	store, err := Store.Get(c)
	if err != nil {
		panic(err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "Failed to parse form"})
	}

	file, err := form.File["folder"][0].Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "Failed to read folder"})
	}
	defer file.Close()

	savePath := fmt.Sprintf("./uploads/%s.zip", store.ID())
	dst, err := os.Create(savePath)
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "Failed to copy zip file contents"})
	}
	err = CopyFolderToContainer(c, savePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "Failed to upload to container"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Sucessfully uploaded folder"})
}

func GetTtyd(c *fiber.Ctx) error {
	store, err := Store.Get(c)
	if err != nil {
		panic(err)
	}

	var container *ContainerInfo
	if store.Fresh() {
		container = pool.GetOne()
		container.session = store.ID()

		store.Set("container_id", container.id)
	} else {
		container_id := store.Get("container_id").(string)
		container = pool.GetContainerById(container_id)
	}

	redirect := fmt.Sprintf("/ttyd/%s", container.id)
	c.Status(fiber.StatusFound)
	c.Append("Location", redirect)

	store.Save()
	return c.JSON(fiber.Map{"status": "redirect"})
}

func PortReportHandler(ctx *fiber.Ctx) error {
	payload := struct {
		ContainerId string `json:"container_id"`
		OpenPorts   []int  `json:"open_ports"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(payload)
	container := pool.GetContainerById(payload.ContainerId)
	if container == nil {
		return ctx.JSON(fiber.Map{"status": "failure"})
	}

	container.open_ports = payload.OpenPorts

	portMap := container.GenerateWebLinkConfig()

	fmt.Println(portMap)

	return ctx.JSON(fiber.Map{"status": "success"})
}
