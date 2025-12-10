package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) GetHealth(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Service is up and running!"})
}
