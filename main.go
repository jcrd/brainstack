package main

import (
	"embed"
	"flag"
	"path/filepath"

	cs "github.com/AndreiTelteu/wails-configstore"
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

	configStore, err := cs.NewConfigStore("brainstack")
	if err != nil {
		println("Failed to initialize config store", err.Error())
		return
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "brainstack",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "8839ea48-dd0e-475b-81a0-cb3b652defd5",
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
		Bind: []interface{}{
			app,
			configStore,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
