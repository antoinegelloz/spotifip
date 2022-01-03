package fip

type Fip struct {
	Name           string `json:"name"`
	Now            Now    `json:"now"`
	Next           Next   `json:"next"`
	DelayToRefresh int    `json:"delayToRefresh"`
	Slug           string `json:"slug"`
	Media          Media  `json:"media"`
	Visual         Visual `json:"visual"`
}

type Now struct {
	FirstLine  string  `json:"firstLine"`
	SecondLine string  `json:"secondLine"`
	ThirdLine  string  `json:"thirdLine"`
	Cover      Cover   `json:"cover"`
	Song       Song    `json:"song"`
	NowTime    int     `json:"nowTime"`
	NowPercent float64 `json:"nowPercent"`
}

type Next struct {
	FirstLine  string `json:"firstLine"`
	SecondLine string `json:"secondLine"`
	ThirdLine  string `json:"thirdLine"`
	Cover      Cover  `json:"cover"`
	Song       Song   `json:"song"`
}

type Media struct {
	Sources   []Sources `json:"sources"`
	StartTime int       `json:"startTime"`
	EndTime   int       `json:"endTime"`
}

type Visual struct {
	Src       string `json:"src"`
	WebpSrc   string `json:"webpSrc"`
	Legend    string `json:"legend"`
	Copyright string `json:"copyright"`
	Author    string `json:"author"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Preview   string `json:"preview"`
}

type Cover struct {
	Src       string      `json:"src"`
	WebpSrc   string      `json:"webpSrc"`
	Legend    interface{} `json:"legend"`
	Copyright interface{} `json:"copyright"`
	Author    interface{} `json:"author"`
	Width     int         `json:"width"`
	Height    int         `json:"height"`
	Preview   interface{} `json:"preview"`
}
type Release struct {
	Title string `json:"title"`
	Label string `json:"label"`
}
type Song struct {
	Title   string  `json:"title"`
	Year    int     `json:"year"`
	Release Release `json:"release"`
}

type Sources struct {
	URL           string `json:"url"`
	BroadcastType string `json:"broadcastType"`
	Format        string `json:"format"`
	Bitrate       int    `json:"bitrate"`
}
