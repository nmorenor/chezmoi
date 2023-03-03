package main

import (
	"github.com/EngoEngine/engo"
	"github.com/nmorenor/chezmoi/scenes"
)

func main() {
	opts := engo.RunOptions{
		Title:          "Chez moi",
		Width:          800,
		Height:         600,
		StandardInputs: true,
	}
	engo.RegisterScene(&scenes.MainScene{})
	engo.RegisterScene(&scenes.SetUserNameScene{})
	engo.RegisterScene(&scenes.SetJoinUserNameScene{})
	engo.RegisterScene(&scenes.SetJoinSessionNameScene{})
	startOptionsScene := &scenes.StartOptionsScene{}
	engo.Run(opts, startOptionsScene)
}
