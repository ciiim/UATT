package bsd_testtool

type CanvasConfig struct {
	CanvasFileName string     `json:"CanvasName"`
	Data           CanvasData `json:"Data"`
}

type Position struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type CanvasComponent struct {
	ID        string   `json:"ID"`
	Type      string   `json:"Type"`
	Label     string   `json:"Label"`
	AttachApp string   `json:"AttachApp"`
	Value     string   `json:"Value"`
	Position  Position `json:"Position"`
}

type CanvasComponentConnection struct {
	FromID string `json:"FromID"`
	ToID   string `jsonl:"ToID"`
}

type CanvasData struct {
	ComponentList []CanvasComponent           `json:"ComponentList"`
	Connections   []CanvasComponentConnection `json:"Connections"`
}
