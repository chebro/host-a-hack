package main

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	api := router.Group("/api")
	LandingRoutes(api)
	FolderRoutes(api)
	TtydRoutes(api)
	router.Post("/portreport", PortReportHandler)
}

func LandingRoutes(router fiber.Router) {
	router.Get("/", GetLanding)
}

func FolderRoutes(router fiber.Router) {
	folder := router.Group("/folder")
	folder.Post("/", PostFolder)
}

func TtydRoutes(router fiber.Router) {
	ttyd := router.Group("/ttyd")
	ttyd.Get("/", GetTtyd)
}
