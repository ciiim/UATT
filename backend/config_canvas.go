package bsd_testtool

type CanvasConfigBase struct {
	CanvasName string     `json:"CanvasName"`
	Data       CanvasData `json:"Data"`
}

type CanvasComponent struct {
	ID        string `json:"ID"`
	Type      string `json:"Type"`
	Label     string `json:"Label"`
	AttachApp string `json:"AttachApp"`
}

type CanvasComponentConnection struct {
	FromID string `json:"FromID"`
	ToID   string `jsonl:"ToID"`
}

type CanvasData struct {
	ComponentList []CanvasComponent           `json:"ComponentList"`
	Connections   []CanvasComponentConnection `json:"Connections"`
}
