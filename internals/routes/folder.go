package folderRoutes

import (
	folderHandler "github.com/chebro/host-a-hack/internals/handlers"
	"github.com/gofiber/fiber/v2"
)

func FolderRoutes (router fiber.Router) {
  folder := router.Group("/folder")
  folder.Post("/", folderHandler.PostFolder)
}
