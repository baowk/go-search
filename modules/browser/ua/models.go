package ua

type UserAgent struct {
	Browser  string `json:"b"`
	Platform string `json:"p"`
	Os       string `json:"os"`
	Ver      string `json:"ver"`
	MainVer  string `json:"mv"`
	UA       string `json:"ua"`
	Mob      bool   `json:"mob"`
	Bot      bool   `json:"bot"`
	// Weight   float64 `json:"weight"`
	// Vendor   string  `json:"vendor"`
}

type Connection struct {
	Downlink      float32 `json:"downlink"`
	DownlinkMax   float32 `json:"downlinkMax"`
	EffectiveType string  `json:"effectiveType"`
	Rtt           int     `json:"rtt"`
	Type          string  `json:"type"`
}

type UserAgentData struct {
	AppName        string     `json:"appName"`
	Connection     Connection `json:"connection"`
	Language       string     `json:"language"`
	Platform       string     `json:"platform"`
	PluginsLength  int        `json:"pluginsLength"`
	ScreenHeight   int        `json:"screenHeight"`
	ScreenWidth    int        `json:"screenWidth"`
	UserAgent      string     `json:"userAgent"`
	Vendor         string     `json:"vendor"`
	ViewportHeight int        `json:"viewportHeight"`
	ViewportWidth  int        `json:"viewportWidth"`
	Weight         float64    `json:"weight"`
	DeviceCategory string     `json:"deviceCategory"`
}
