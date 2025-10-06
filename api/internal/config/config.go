package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port               string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	ScopeChapters      []string
	SOMBase            float64
	SOMAggressive      float64
	PartnerWeightsFile string
	TariffScenariosFile string
}

type PartnerWeights map[string]map[string]float64

type TariffScenario struct {
	Default  map[string]float64            `yaml:"default"`
	Chapters map[string]map[string]float64 `yaml:"chapters"`
	Years    map[string]struct {
		Default  map[string]float64            `yaml:"default"`
		Chapters map[string]map[string]float64 `yaml:"chapters"`
	} `yaml:"years"`
}

type TariffScenarios struct {
	Scenarios map[string]TariffScenario `yaml:"scenarios"`
}

func LoadConfig() *AppConfig {
	cfg := &AppConfig{
		Port:                getenv("PORT", "8080"),
		DBHost:              getenv("PGHOST", "localhost"),
		DBPort:              getenv("PGPORT", "5432"),
		DBUser:              getenv("PGUSER", "postgres"),
		DBPassword:          getenv("PGPASSWORD", ""),
		DBName:              getenv("PGDATABASE", "postgres"),
		ScopeChapters:       []string{"02", "08", "84", "85"},
		SOMBase:             0.015,
		SOMAggressive:       0.03,
		PartnerWeightsFile:  getenv("PARTNER_WEIGHTS_FILE", "./config/partners_stub.yaml"),
		TariffScenariosFile: getenv("TARIFF_SCENARIOS_FILE", "./config/tariff_scenarios.yaml"),
	}

	if v := getenv("SCOPE_CHAPTERS", ""); v != "" {
		parts := strings.Split(v, ",")
		clean := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if len(p) == 1 {
				p = "0" + p
			}
			if len(p) >= 2 {
				clean = append(clean, p[:2])
			}
		}
		if len(clean) > 0 {
			cfg.ScopeChapters = clean
		}
	}

	if v := getenv("SOM_BASE", ""); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.SOMBase = f
		}
	}
	if v := getenv("SOM_AGGRESSIVE", ""); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.SOMAggressive = f
		}
	}

	return cfg
}

func LoadPartnerWeights(path string) PartnerWeights {
	var doc struct {
		Partners map[string]map[string]float64 `yaml:"partners"`
	}

	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("partner weights not found (%s): using defaults", path)
		return nil
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		log.Printf("failed to parse partner weights: %v", err)
		return nil
	}
	if len(doc.Partners) == 0 {
		log.Printf("partner weights file is empty: %s", path)
		return nil
	}

	out := make(PartnerWeights)
	for chapterKey, partnersMap := range doc.Partners {
		chKey := strings.TrimSpace(chapterKey)
		if chKey != "default" && len(chKey) == 1 {
			chKey = "0" + chKey
		}
		if _, ok := out[chKey]; !ok {
			out[chKey] = make(map[string]float64)
		}
		for partnerCode, weight := range partnersMap {
			p := strings.ToUpper(strings.TrimSpace(partnerCode))
			out[chKey][p] = weight
		}
	}
	return out
}

func LoadTariffScenarios(path string) *TariffScenarios {
	var tariffs TariffScenarios
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("tariff scenarios not found (%s): continuing without tariffs", path)
		return &tariffs
	}
	if err := yaml.Unmarshal(b, &tariffs); err != nil {
		log.Printf("failed to parse tariff scenarios: %v", err)
	}
	return &tariffs
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
