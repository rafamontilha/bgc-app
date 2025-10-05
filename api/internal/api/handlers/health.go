package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bgc-app/internal/business/health"
)

type HealthHandler struct {
	service health.Service
}

func NewHealthHandler(service health.Service) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	status := h.service.GetHealthStatus()
	c.JSON(http.StatusOK, status)
}
