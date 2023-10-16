package fip

type Fip struct {
	DelayToRefresh int     `json:"delayToRefresh"`
	Now            Track   `json:"now"`
	Migrated       bool    `json:"migrated"`
	Next           []Track `json:"next"`
}

type Track struct {
	FirstLine  string `json:"firstLine"`
	SecondLine string `json:"secondLine"`
	Song       Song   `json:"song"`
}

type Sources struct {
	URL           string `json:"url"`
	BroadcastType string `json:"broadcastType"`
	Format        string `json:"format"`
	Bitrate       int    `json:"bitrate"`
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
