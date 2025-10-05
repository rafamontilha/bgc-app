package route

type Repository interface {
	GetTAMByYearAndChapter(year int, chapter string) (float64, error)
}
