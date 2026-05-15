package main

import (
	"github.com/robinskaba/roge/cmd"
	"github.com/robinskaba/roge/internal/pkg"
)

func main() {
	// try removing old executable if exists
	pkg.TryRemovingOldVersion()
	cmd.Execute()
}
