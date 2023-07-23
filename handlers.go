package main

import (
	"fmt"
	"os"
  "io"

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
	fmt.Sprintln("hi")
  form, err := c.MultipartForm()
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse form")
  }
  files := form.File["folder"]
  file, err := files[0].Open()
  defer file.Close()
  fmt.Println(file)
  savePath := "./uploads/upload.zip"
  dst, err := os.Create(savePath)
  defer dst.Close()
  f, err := io.Copy(dst, file)
  fmt.Println(f)
	return c.JSON(fiber.Map{"status": "success", "message": "uploaded folder"})

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
