package market

type MarketItem struct {
	Ano        int     `json:"ano"`
	NCMChapter string  `json:"ncm_chapter"`
	ValorUSD   float64 `json:"valor_usd"`
}

type MarketSizeRequest struct {
	Metric     string
	YearFrom   int
	YearTo     int
	NCMChapter string
	Scenario   string
}

type MarketSizeResponse struct {
	Metric   string       `json:"metric"`
	Scenario string       `json:"scenario"`
	Items    []MarketItem `json:"items"`
}
