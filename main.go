package main

import (
	"context"

	"github.com/agejevasv/goclipher/internal/icon"
	"github.com/cdfmlr/ellipsis"
	"github.com/getlantern/systray"
	"golang.design/x/clipboard"
)

const MaxHistory = 10
const EllipsisAt = 42

func main() {
	systray.Run(onReady, nil)
}

func onReady() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Skip", false)

	systray.SetTemplateIcon(icon.Data, icon.Data)

	clips := make([]string, MaxHistory)
	menus := make([]*systray.MenuItem, MaxHistory)

	for i := 0; i < MaxHistory; i++ {
		item := systray.AddMenuItem("", "")
		item.Hide()
		menus[i] = item
	}

	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit Goclipher", "Quit")

	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()

	ch := clipboard.Watch(ctx, clipboard.FmtText)
	currIdx := 0

	for data := range ch {
		if ctx.Value("Skip").(bool) {
			ctx = context.WithValue(ctx, "Skip", false)
			continue
		}

		clips = append([]string{string(data)}, clips[:MaxHistory-1]...)

		for i := 0; i < MaxHistory; i++ {
			if clips[i] != "" {
				menus[i].Uncheck()
				menus[i].SetTitle(ellipsis.Ending(clips[i], EllipsisAt))
				menus[i].Show()
			}

			if i == 0 {
				currIdx = 0
				menus[i].Check()
			}

			go func() {
				<-menus[i].ClickedCh
				ctx = context.WithValue(ctx, "Skip", true)
				menus[currIdx].Uncheck()
				menus[i].Check()
				currIdx = i
				clipboard.Write(clipboard.FmtText, []byte(clips[i]))
			}()
		}
	}
}
