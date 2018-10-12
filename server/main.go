package main

import (
	"log"

	"github.com/marcusolsson/tui-go"
)

type ChatClient struct {
	loginView tui.Widget
	chatView  tui.Widget
	container tui.UI
}

func main() {

	root := tui.NewVBox()

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
