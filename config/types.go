package config

type Config struct {
	FramesPerSecond int                     `json:"frames_per_second"`
	Viewports       map[string]HostViewport `json:"viewports"`
	Clients         map[string]Client       `json:"clients"`
}

type HostViewport struct {
	PosX   int `json:"posX"`
	PosY   int `json:"posY"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ClientViewports struct {
	ID       string `json:"id"`
	DisplayX int    `json:"displayX"`
	DisplayY int    `json:"displayY"`
}

type Client struct {
	ID        string            `json:"id"`
	Viewports []ClientViewports `json:"viewports"`
}
