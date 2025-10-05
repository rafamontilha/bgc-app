package health

import (
	"bgc-app/internal/config"
)

type HealthStatus struct {
	Select             int      `json:"select"`
	Status             string   `json:"status"`
	ChaptersOnda1      []string `json:"chapters_onda1"`
	PartnerWeights     bool     `json:"partner_weights"`
	TariffsLoaded      bool     `json:"tariffs_loaded"`
	AvailableScenarios []string `json:"available_scenarios"`
}

type Service interface {
	GetHealthStatus() *HealthStatus
}

type service struct {
	config  *config.AppConfig
	weights config.PartnerWeights
	tariffs *config.TariffScenarios
}

func NewService(cfg *config.AppConfig, weights config.PartnerWeights, tariffs *config.TariffScenarios) Service {
	return &service{
		config:  cfg,
		weights: weights,
		tariffs: tariffs,
	}
}

func (s *service) GetHealthStatus() *HealthStatus {
	scenarios := make([]string, 0, len(s.tariffs.Scenarios))
	for k := range s.tariffs.Scenarios {
		scenarios = append(scenarios, k)
	}

	return &HealthStatus{
		Select:             1,
		Status:             "ok",
		ChaptersOnda1:      s.config.ScopeChapters,
		PartnerWeights:     (s.weights != nil),
		TariffsLoaded:      (len(s.tariffs.Scenarios) > 0),
		AvailableScenarios: scenarios,
	}
}
