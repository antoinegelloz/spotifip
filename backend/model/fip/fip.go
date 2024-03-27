package fip

type Fip struct {
	StationName    string `json:"stationName"`
	DelayToRefresh int    `json:"delayToRefresh"`
	Migrated       bool   `json:"migrated"`
	Now            Now    `json:"now"`
	Next           Next   `json:"next"`
}
type FirstLine struct {
	Title string `json:"title"`
	ID    any    `json:"id"`
	Path  any    `json:"path"`
}
type SecondLine struct {
	Title string `json:"title"`
	ID    any    `json:"id"`
	Path  any    `json:"path"`
}
type ThirdLine struct {
	Title any `json:"title"`
	ID    any `json:"id"`
	Path  any `json:"path"`
}
type Release struct {
	Label     string `json:"label"`
	Title     string `json:"title"`
	Reference any    `json:"reference"`
}
type Song struct {
	ID      string  `json:"id"`
	Year    int     `json:"year"`
	Release Release `json:"release"`
}
type Sources struct {
	URL           string `json:"url"`
	BroadcastType string `json:"broadcastType"`
	Format        string `json:"format"`
	Bitrate       int    `json:"bitrate"`
}
type Media struct {
	StartTime int       `json:"startTime"`
	EndTime   int       `json:"endTime"`
	Sources   []Sources `json:"sources"`
}
type Card struct {
	Model     string `json:"model"`
	Src       string `json:"src"`
	WebpSrc   string `json:"webpSrc"`
	Legend    any    `json:"legend"`
	Copyright any    `json:"copyright"`
	Author    any    `json:"author"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Preview   string `json:"preview"`
	ID        string `json:"id"`
	Type      string `json:"type"`
	Preset    string `json:"preset"`
}
type Player struct {
	Model     string `json:"model"`
	Src       string `json:"src"`
	WebpSrc   string `json:"webpSrc"`
	Legend    any    `json:"legend"`
	Copyright any    `json:"copyright"`
	Author    any    `json:"author"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Preview   string `json:"preview"`
	ID        string `json:"id"`
	Type      string `json:"type"`
	Preset    string `json:"preset"`
}
type Visuals struct {
	Card   Card   `json:"card"`
	Player Player `json:"player"`
}
type Now struct {
	PrintProgMusic bool       `json:"printProgMusic"`
	StartTime      int        `json:"startTime"`
	EndTime        int        `json:"endTime"`
	Producer       string     `json:"producer"`
	FirstLine      FirstLine  `json:"firstLine"`
	SecondLine     SecondLine `json:"secondLine"`
	ThirdLine      ThirdLine  `json:"thirdLine"`
	Song           Song       `json:"song"`
	Media          Media      `json:"media"`
	LocalRadios    []any      `json:"localRadios"`
	Visuals        Visuals    `json:"visuals"`
}
type Next struct {
	PrintProgMusic bool       `json:"printProgMusic"`
	StartTime      int        `json:"startTime"`
	EndTime        int        `json:"endTime"`
	Producer       string     `json:"producer"`
	FirstLine      FirstLine  `json:"firstLine"`
	SecondLine     SecondLine `json:"secondLine"`
	ThirdLine      ThirdLine  `json:"thirdLine"`
	Song           any        `json:"song"`
}
