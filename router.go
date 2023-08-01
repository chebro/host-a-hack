package main

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	api := router.Group("/api")
	LandingRoutes(api)
	FolderRoutes(api)
	router.Get("/allocatettyd", GetTtyd)
	router.Post("/portreport", PortReportHandler)
	router.Post("/sessionend", SessionEndHandler)
}

func LandingRoutes(router fiber.Router) {
	router.Get("/", GetLanding)
}

func FolderRoutes(router fiber.Router) {
	folder := router.Group("/folder")
	folder.Post("/", PostFolder)
}
