package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/marcusolsson/tui-go"
)

func main() {
	sidebar := tui.NewVBox(
		tui.NewLabel("        USERS       "),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: "localhost:8123", Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial: ", err)
	}
	defer c.Close()

	done := make(chan struct{})

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	input.OnSubmit(func(e *tui.Entry) {
		if len(e.Text()) == 0 {
			return
		}
		c.WriteMessage(websocket.TextMessage, []byte(e.Text()))
		input.SetText("")

	})

	go func() {
		for {
			select {
			case <-interupt:
				log.Println("interrupt")
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					ui.Quit()
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				ui.Quit()
				return
			}
		}
	}()

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				return
			}
			ui.Update(func() {
				history.Append(tui.NewHBox(
					tui.NewLabel(time.Now().Format("15:04")),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "test"))),
					tui.NewLabel(string(message)),
					tui.NewSpacer(),
				))
			})

		}
	}()
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
