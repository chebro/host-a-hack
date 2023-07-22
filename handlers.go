package main

import (
	"fmt"

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

	link := fmt.Sprintf("/ttyd/%s", container.id)
	c.Status(fiber.StatusFound)
	c.Append("Location", link)

	store.Save()
	return c.JSON(fiber.Map{"status": "redirect"})
}
