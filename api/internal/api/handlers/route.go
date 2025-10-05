package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"bgc-app/internal/business/route"
)

type RouteHandler struct {
	service route.Service
}

func NewRouteHandler(service route.Service) *RouteHandler {
	return &RouteHandler{service: service}
}

func (h *RouteHandler) CompareRoutes(c *gin.Context) {
	from := strings.ToUpper(c.DefaultQuery("from", "USA"))
	altsRaw := c.DefaultQuery("alts", "CHN,ARE,SAU,IND")

	alts := make([]string, 0)
	for _, a := range strings.Split(altsRaw, ",") {
		a = strings.TrimSpace(strings.ToUpper(a))
		if a != "" && a != from {
			alts = append(alts, a)
		}
	}

	year, err := strconv.Atoi(c.DefaultQuery("year", "2024"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	chapter := c.Query("ncm_chapter")
	if len(chapter) == 1 {
		chapter = "0" + chapter
	}
	if chapter == "" || len(chapter) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ncm_chapter (2 dígitos) é obrigatório"})
		return
	}

	scenario := c.DefaultQuery("tariff_scenario", "base")

	req := route.RouteCompareRequest{
		From:       from,
		Alts:       alts,
		Year:       year,
		NCMChapter: chapter,
		Scenario:   scenario,
	}

	result, err := h.service.CompareRoutes(req)
	if err != nil {
		if strings.Contains(err.Error(), "sem dados") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
