package main

import (
	"github.com/Nebula-Pack/repo-emmulator/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/clone", handler.CloneRepo)

	r.Run(":4242")
}
