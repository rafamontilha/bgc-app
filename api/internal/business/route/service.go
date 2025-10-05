package route

import (
	"database/sql"
	"fmt"
	"strings"

	"bgc-app/internal/config"
)

type Service interface {
	CompareRoutes(req RouteCompareRequest) (*RouteCompareResponse, error)
}

type service struct {
	repo    Repository
	weights config.PartnerWeights
	tariffs *config.TariffScenarios
}

func NewService(repo Repository, weights config.PartnerWeights, tariffs *config.TariffScenarios) Service {
	return &service{
		repo:    repo,
		weights: weights,
		tariffs: tariffs,
	}
}

func (s *service) CompareRoutes(req RouteCompareRequest) (*RouteCompareResponse, error) {
	tam, err := s.repo.GetTAMByYearAndChapter(req.Year, req.NCMChapter)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sem dados para ano/capítulo: year=%d chapter=%s", req.Year, req.NCMChapter)
		}
		return nil, err
	}

	partners := append([]string{req.From}, req.Alts...)
	weights := s.calculateWeights(partners, req.NCMChapter, req.From, req.Alts)

	scn, hasScenario := s.tariffs.Scenarios[req.Scenario]
	tariffApplied := false

	results := make([]RouteItem, 0, len(partners))
	adjustedTotal := 0.0

	for _, p := range partners {
		share := weights[p]
		est := share * tam

		factor := 1.0
		if hasScenario {
			factor = s.factorFor(scn, req.Year, req.NCMChapter, p)
			if factor != 1.0 {
				tariffApplied = true
			}
		}
		est = est * factor
		adjustedTotal += est

		results = append(results, RouteItem{
			Partner:      p,
			Share:        share,
			Factor:       factor,
			EstimatedUSD: est,
		})
	}

	return &RouteCompareResponse{
		Year:             req.Year,
		NCMChapter:       req.NCMChapter,
		Basis:            "TAM (mview)",
		TAMTotalUSD:      tam,
		From:             req.From,
		Alts:             req.Alts,
		TariffScenario:   req.Scenario,
		TariffApplied:    tariffApplied,
		AdjustedTotalUSD: adjustedTotal,
		Note:             "stub com pesos + fatores de tarifa; substituir por dados reais por parceiro em próxima onda",
		Results:          results,
	}, nil
}

func (s *service) calculateWeights(partners []string, chapter, from string, alts []string) map[string]float64 {
	weights := map[string]float64{}

	if s.weights != nil {
		if w, ok := s.weights[chapter]; ok {
			for k, v := range w {
				weights[strings.ToUpper(k)] = v
			}
		}
		if len(weights) == 0 {
			if w, ok := s.weights["default"]; ok {
				for k, v := range w {
					weights[strings.ToUpper(k)] = v
				}
			}
		}
	}

	if len(weights) == 0 {
		weights[from] = 0.40
		if len(alts) > 0 {
			rem := 0.60 / float64(len(alts))
			for _, a := range alts {
				weights[a] = rem
			}
		}
	}

	sum := 0.0
	for _, p := range partners {
		sum += weights[p]
	}
	if sum == 0 {
		eq := 1.0 / float64(len(partners))
		for _, p := range partners {
			weights[p] = eq
		}
		sum = 1.0
	}
	for k, v := range weights {
		weights[k] = v / sum
	}

	return weights
}

func (s *service) factorFor(scn config.TariffScenario, year int, chapter, partner string) float64 {
	p := strings.ToUpper(partner)
	chap := chapter
	ys := fmt.Sprintf("%d", year)

	if y, ok := scn.Years[ys]; ok {
		if mp, ok := y.Chapters[chap]; ok {
			if f, ok := mp[p]; ok {
				return f
			}
		}
		if f, ok := y.Default[p]; ok {
			return f
		}
	}
	if mp, ok := scn.Chapters[chap]; ok {
		if f, ok := mp[p]; ok {
			return f
		}
	}
	if f, ok := scn.Default[p]; ok {
		return f
	}
	return 1.0
}
