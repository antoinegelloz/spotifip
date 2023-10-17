package fip

type Fip struct {
	Now  Track   `json:"now"`
	Next []Track `json:"next"`
}

type Track struct {
	FirstLine  string     `json:"firstLine"`
	SecondLine string     `json:"secondLine"`
	ThirdLine  string     `json:"thirdLine"`
	StartTime  int        `json:"startTime"`
	EndTime    int        `json:"endTime"`
	CardVisual CardVisual `json:"cardVisual"`
	Song       Song       `json:"song"`
}

type CardVisual struct {
	Type  string `json:"type"`
	Src   string `json:"src"`
	Model string `json:"model"`
}

type Song struct {
	ID      string  `json:"id"`
	Year    int     `json:"year"`
	Release Release `json:"release"`
}

type Release struct {
	Title string `json:"title"`
	Label string `json:"label"`
}
