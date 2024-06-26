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

	isLua, hasRockspec, scanResponse, err := service.CloneRepository(req.Repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{"isLua": isLua}
	if hasRockspec {
		response["hasRockspec"] = hasRockspec
		response["scanResponse"] = scanResponse
	}

	c.JSON(http.StatusOK, response)
}
