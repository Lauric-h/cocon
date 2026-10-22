package view

type WeekNavigation struct {
	Label       string
	PreviousURL string
	CurrentURL  string
	NextURL     string
	Query       string
	IsCurrent   bool
}
