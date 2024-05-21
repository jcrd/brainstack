package main

import (
	"embed"
	"flag"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	configPath := configdir.LocalConfig("brainstack")
	if err := configdir.MakePath(configPath); err != nil {
		panic(err)
	}

	dbPath := flag.String("db", filepath.Join(configPath, "main.db"), "Path to DB file")
	flag.Parse()

	// Create an instance of the app structure
	app := NewApp(*dbPath)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "brainstack",
		Width:  525,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
