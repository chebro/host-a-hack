package main

import (
  "fmt"
	"github.com/gofiber/fiber/v2"
)

func GetLanding (c *fiber.Ctx) error {
  store, err := Store.Get(c)
  if err != nil {
    panic(err)
  }
  return c.JSON(fiber.Map{"status": "success", "id": store.ID()})
}

func PostFolder (c *fiber.Ctx) error {
  return c.JSON(fiber.Map{"status": "success", "message": "uploaded folder"})
}

func GetTtyd (c *fiber.Ctx) error {
  store, err := Store.Get(c)
  if err != nil {
    panic(err)
  }

  container := pool.GetOne()
  container.session = store.ID()

  link := fmt.Sprintf("/ttyd/%s", container.id)
  c.Status(fiber.StatusFound)
  c.Append("Location", link)
  return c.JSON(fiber.Map{"status": "redirect"})
}
