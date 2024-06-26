package handler

import (
	"net/http"

	"github.com/Nebula-Pack/repo-emmulator/internal/service"
	"github.com/gin-gonic/gin"
)

type CloneRequest struct {
	Repo string `json:"repo" binding:"required"`
}

func CloneRepo(c *gin.Context) {
	var req CloneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CloneRepository(req.Repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
