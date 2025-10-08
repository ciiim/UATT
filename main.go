package main

import (
	bsd_testtool "bsd_testtool/backend"
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	manager := &bsd_testtool.GlobalManager

	if err := manager.Init(""); err != nil {
		fmt.Printf("manager init error: %s\n", err.Error())
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "bsd_testtool",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 200, G: 200, B: 200, A: 1},
		OnStartup:        manager.Startup,
		Bind: []any{
			manager,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
