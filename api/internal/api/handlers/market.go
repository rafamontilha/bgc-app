package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"bgc-app/internal/business/market"
)

type MarketHandler struct {
	service market.Service
}

func NewMarketHandler(service market.Service) *MarketHandler {
	return &MarketHandler{service: service}
}

func (h *MarketHandler) GetMarketSize(c *gin.Context) {
	metric := strings.ToUpper(c.Query("metric"))
	if metric == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "metric is required [TAM|SAM|SOM]"})
		return
	}

	yearFrom, _ := strconv.Atoi(c.DefaultQuery("year_from", "2020"))
	yearTo, _ := strconv.Atoi(c.DefaultQuery("year_to", "2025"))
	ncmChapter := c.Query("ncm_chapter")
	scenario := strings.ToLower(c.DefaultQuery("scenario", "base"))

	req := market.MarketSizeRequest{
		Metric:     metric,
		YearFrom:   yearFrom,
		YearTo:     yearTo,
		NCMChapter: ncmChapter,
		Scenario:   scenario,
	}

	result, err := h.service.CalculateMarketSize(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
