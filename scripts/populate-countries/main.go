package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// RestCountriesResponse representa a resposta da API REST Countries
type RestCountriesResponse struct {
	Name struct {
		Common   string            `json:"common"`
		Official string            `json:"official"`
		NativeName map[string]struct {
			Official string `json:"official"`
			Common   string `json:"common"`
		} `json:"nativeName"`
	} `json:"name"`
	Cca2        string   `json:"cca2"`
	Region      string   `json:"region"`
	Subregion   string   `json:"subregion"`
	Languages   map[string]string `json:"languages"`
	Latlng      []float64 `json:"latlng"`
	Area        float64   `json:"area"`
	Population  int64     `json:"population"`
	Flag        string    `json:"flag"`
	Currencies  map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	CapitalInfo struct {
		Latlng []float64 `json:"latlng"`
	} `json:"capitalInfo"`
}

// Top 50 países de comércio exterior do Brasil
var topCountryCodes = []string{
	"CN", "US", "AR", "NL", "CL", "DE", "JP", "IN", "MX", "ES",
	"IT", "FR", "GB", "BE", "KR", "RU", "CA", "PE", "CO", "TH",
	"PY", "UY", "VE", "SA", "AE", "TR", "PL", "ZA", "MY", "SG",
	"ID", "PH", "VN", "PT", "CH", "AT", "SE", "NO", "DK", "FI",
	"IE", "GR", "CZ", "RO", "HU", "AU", "NZ", "EG", "MA", "NG",
}

func main() {
	// Conecta ao PostgreSQL
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "bgc")
	dbPass := getEnv("DB_PASS", "bgc")
	dbName := getEnv("DB_NAME", "bgc")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Testa conexão
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	// Busca dados dos países
	countries, err := fetchCountries(topCountryCodes)
	if err != nil {
		log.Fatal("Failed to fetch countries:", err)
	}

	log.Printf("Fetched %d countries from REST Countries API\n", len(countries))

	// Insere no banco
	inserted := 0
	for _, country := range countries {
		if err := insertCountry(db, country); err != nil {
			log.Printf("Warning: Failed to insert country %s: %v\n", country.Cca2, err)
			continue
		}
		inserted++
	}

	log.Printf("Successfully inserted %d countries into database\n", inserted)
}

func fetchCountries(codes []string) ([]RestCountriesResponse, error) {
	var countries []RestCountriesResponse

	for _, code := range codes {
		url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s", code)

		log.Printf("Fetching data for country: %s\n", code)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s: %w", code, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Warning: Received status %d for country %s\n", resp.StatusCode, code)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response for %s: %w", code, err)
		}

		var result []RestCountriesResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON for %s: %w", code, err)
		}

		if len(result) > 0 {
			countries = append(countries, result[0])
		}

		// Rate limiting (respeita API pública)
		time.Sleep(100 * time.Millisecond)
	}

	return countries, nil
}

func insertCountry(db *sql.DB, country RestCountriesResponse) error {
	// Calcula distância aproximada de Brasília (simplificado)
	distance := calculateDistance(country)

	// Pega primeiro currency code
	currencyCode := ""
	for code := range country.Currencies {
		currencyCode = code
		break
	}

	// Converte languages map para array
	var languages []string
	for _, lang := range country.Languages {
		languages = append(languages, lang)
	}

	query := `
		INSERT INTO public.countries_metadata (
			code, name_pt, name_en, region, subregion,
			distance_brazil_km, latitude, longitude,
			population, flag_emoji, currency_code, languages,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12,
			now(), now()
		)
		ON CONFLICT (code) DO UPDATE SET
			name_pt = EXCLUDED.name_pt,
			name_en = EXCLUDED.name_en,
			region = EXCLUDED.region,
			subregion = EXCLUDED.subregion,
			distance_brazil_km = EXCLUDED.distance_brazil_km,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			population = EXCLUDED.population,
			flag_emoji = EXCLUDED.flag_emoji,
			currency_code = EXCLUDED.currency_code,
			languages = EXCLUDED.languages,
			updated_at = now()
	`

	var lat, lon float64
	if len(country.Latlng) >= 2 {
		lat = country.Latlng[0]
		lon = country.Latlng[1]
	}

	_, err := db.Exec(query,
		country.Cca2,
		country.Name.Common, // Usamos Common como PT (simplificado)
		country.Name.Official,
		country.Region,
		country.Subregion,
		distance,
		lat,
		lon,
		country.Population,
		country.Flag,
		currencyCode,
		languages,
	)

	return err
}

// calculateDistance calcula distância aproximada de Brasília ao país
// Simplificado: usa latitude/longitude para estimar
func calculateDistance(country RestCountriesResponse) int {
	// Brasília: -15.7975, -47.8919
	brasiliaLat := -15.7975
	brasiliaLon := -47.8919

	var targetLat, targetLon float64

	// Usa capital coordinates se disponível, senão usa country center
	if len(country.CapitalInfo.Latlng) >= 2 {
		targetLat = country.CapitalInfo.Latlng[0]
		targetLon = country.CapitalInfo.Latlng[1]
	} else if len(country.Latlng) >= 2 {
		targetLat = country.Latlng[0]
		targetLon = country.Latlng[1]
	}

	// Fórmula de Haversine (simplificada)
	const earthRadius = 6371 // km

	dLat := deg2rad(targetLat - brasiliaLat)
	dLon := deg2rad(targetLon - brasiliaLon)

	a := sin(dLat/2)*sin(dLat/2) +
		cos(deg2rad(brasiliaLat))*cos(deg2rad(targetLat))*
			sin(dLon/2)*sin(dLon/2)

	c := 2 * atan2(sqrt(a), sqrt(1-a))
	distance := earthRadius * c

	return int(distance)
}

func deg2rad(deg float64) float64 {
	return deg * (3.14159265359 / 180)
}

func sin(x float64) float64 {
	// Usa aproximação de Taylor para sin
	x = x - float64(int(x/(2*3.14159265359)))*(2*3.14159265359)
	result := x
	term := x
	for i := 1; i <= 10; i++ {
		term *= -x * x / float64((2*i)*(2*i+1))
		result += term
	}
	return result
}

func cos(x float64) float64 {
	// cos(x) = sin(x + π/2)
	return sin(x + 3.14159265359/2)
}

func sqrt(x float64) float64 {
	if x == 0 {
		return 0
	}
	// Método de Newton
	z := x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

func atan2(y, x float64) float64 {
	// Aproximação simples de atan2
	if x > 0 {
		return atan(y / x)
	}
	if x < 0 && y >= 0 {
		return atan(y/x) + 3.14159265359
	}
	if x < 0 && y < 0 {
		return atan(y/x) - 3.14159265359
	}
	if x == 0 && y > 0 {
		return 3.14159265359 / 2
	}
	if x == 0 && y < 0 {
		return -3.14159265359 / 2
	}
	return 0
}

func atan(x float64) float64 {
	// Série de Taylor para atan
	if x > 1 {
		return 3.14159265359/2 - atan(1/x)
	}
	if x < -1 {
		return -3.14159265359/2 - atan(1/x)
	}
	result := x
	term := x
	for i := 1; i <= 20; i++ {
		term *= -x * x * float64(2*i-1) / float64(2*i+1)
		result += term
	}
	return result
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
