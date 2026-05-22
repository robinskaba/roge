package main

import (
	"github.com/robinskaba/roge/internal/app"
	"github.com/robinskaba/roge/internal/cmd"
)

func main() {
	// try removing old executable if exists
	app.TryRemovingOldVersion()
	cmd.Execute()
}
