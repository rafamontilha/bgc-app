package route

type RouteCompareRequest struct {
	From       string
	Alts       []string
	Year       int
	NCMChapter string
	Scenario   string
}

type RouteItem struct {
	Partner      string  `json:"partner"`
	Share        float64 `json:"share"`
	Factor       float64 `json:"factor"`
	EstimatedUSD float64 `json:"estimated_usd"`
}

type RouteCompareResponse struct {
	Year             int         `json:"year"`
	NCMChapter       string      `json:"ncm_chapter"`
	Basis            string      `json:"basis"`
	TAMTotalUSD      float64     `json:"tam_total_usd"`
	From             string      `json:"from"`
	Alts             []string    `json:"alts"`
	TariffScenario   string      `json:"tariff_scenario"`
	TariffApplied    bool        `json:"tariff_applied"`
	AdjustedTotalUSD float64     `json:"adjusted_total_usd"`
	Note             string      `json:"note"`
	Results          []RouteItem `json:"results"`
}
