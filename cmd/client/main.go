package main

import (
	"runtime"

	"github.com/hexagon-0/voxel-game/internal/client/app"
)

func main() {
	runtime.LockOSThread()

	app := app.App{}
	app.Run()
}
