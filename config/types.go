package config

type Config struct {
	FramesPerSecond int             `json:"frames_per_second"`
	Viewports       []HostViewports `json:"viewports"`
	Clients         []Clients       `json:"clients"`
}

type HostViewports struct {
	ID     string `json:"id"`
	PosX   int    `json:"posX"`
	PosY   int    `json:"posY"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type ClientViewports struct {
	ID       string `json:"id"`
	DisplayX int    `json:"displayX"`
	DisplayY int    `json:"displayY"`
}

type Clients struct {
	ID        string            `json:"id"`
	Viewports []ClientViewports `json:"viewports"`
}
