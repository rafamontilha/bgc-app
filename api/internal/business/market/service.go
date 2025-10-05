package market

import (
	"fmt"

	"bgc-app/internal/config"
)

type Service interface {
	CalculateMarketSize(req MarketSizeRequest) (*MarketSizeResponse, error)
}

type service struct {
	repo   Repository
	config *config.AppConfig
}

func NewService(repo Repository, cfg *config.AppConfig) Service {
	return &service{
		repo:   repo,
		config: cfg,
	}
}

func (s *service) CalculateMarketSize(req MarketSizeRequest) (*MarketSizeResponse, error) {
	chapters := []string{}
	if req.Metric == "SAM" || req.Metric == "SOM" {
		if len(s.config.ScopeChapters) == 0 {
			return nil, fmt.Errorf("server misconfigured: scope chapters empty")
		}
		chapters = s.config.ScopeChapters
	}

	items, err := s.repo.GetMarketDataByYearRange(req.YearFrom, req.YearTo, chapters, req.NCMChapter)
	if err != nil {
		return nil, err
	}

	processedItems := make([]MarketItem, 0, len(items))
	for _, item := range items {
		mi := item

		switch req.Metric {
		case "TAM", "SAM":
		case "SOM":
			switch req.Scenario {
			case "aggressive":
				mi.ValorUSD = item.ValorUSD * s.config.SOMAggressive
			default:
				mi.ValorUSD = item.ValorUSD * s.config.SOMBase
			}
		default:
			return nil, fmt.Errorf("invalid metric; use TAM|SAM|SOM")
		}

		processedItems = append(processedItems, mi)
	}

	return &MarketSizeResponse{
		Metric:   req.Metric,
		Scenario: req.Scenario,
		Items:    processedItems,
	}, nil
}
