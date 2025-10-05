package market

type Repository interface {
	GetMarketDataByYearRange(yearFrom, yearTo int, chapters []string, ncmChapter string) ([]MarketItem, error)
}
