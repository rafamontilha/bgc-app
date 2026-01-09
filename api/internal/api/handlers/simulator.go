package handlers

import (
	"net/http"
	"time"

	"bgc-app/internal/business/destination"

	"github.com/gin-gonic/gin"
)

// SimulatorHandler handler para endpoints do simulador
type SimulatorHandler struct {
	service destination.ServiceInterface
}

// NewSimulatorHandler cria uma nova instância do handler
func NewSimulatorHandler(service destination.ServiceInterface) *SimulatorHandler {
	return &SimulatorHandler{
		service: service,
	}
}

// SimulateDestinations simula destinos de exportação
// @Summary Simula destinos de exportação
// @Description Retorna recomendações de destinos de exportação baseado em NCM
// @Tags Simulator
// @Accept json
// @Produce json
// @Param request body destination.SimulatorRequest true "Dados da simulação"
// @Success 200 {object} destination.SimulatorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} ErrorResponse
// @Router /v1/simulator/destinations [post]
func (h *SimulatorHandler) SimulateDestinations(c *gin.Context) {
	startTime := time.Now()

	var req destination.SimulatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid_request",
			Message: "Formato de requisição inválido",
			Details: err.Error(),
		})
		return
	}

	// Valida request
	if err := req.ValidateSimulatorRequest(); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Chama service
	resp, err := h.service.RecommendDestinations(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Adiciona tempo de processamento
	resp.Metadata.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	// Retorna sucesso
	c.JSON(http.StatusOK, resp)
}

// handleError trata erros de forma consistente
func (h *SimulatorHandler) handleError(c *gin.Context, err error) {
	switch err {
	case destination.ErrInvalidNCM, destination.ErrInvalidVolume, destination.ErrInvalidMaxResults:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "validation_error",
			Message: err.Error(),
		})
	case destination.ErrNCMNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "ncm_not_found",
			Message: "NCM não encontrado na base de dados",
		})
	case destination.ErrNoDataAvailable:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "no_data_available",
			Message: "Dados não disponíveis para o NCM solicitado",
		})
	case destination.ErrInsufficientData:
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error: "insufficient_data",
			Message: "Dados insuficientes para gerar recomendações",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal_error",
			Message: "Erro interno ao processar requisição",
			Details: err.Error(),
		})
	}
}

// ErrorResponse estrutura padrão de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
