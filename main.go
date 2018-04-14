package main

import (
	"./game"
	"./ui2d"
)

func main() {
	ui := &ui2d.UI2d{}
	game.Run(ui)
}
