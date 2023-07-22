package folderHandler

import (
  "github.com/gofiber/fiber/v2"
)

func PostFolder (c *fiber.Ctx) error {

  return c.JSON(fiber.Map{"status": "success", "message": "uploaded folder"})
}
