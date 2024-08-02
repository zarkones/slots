package core

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var App = app.New()

var MainW = App.NewWindow("Slot Game")

var WindowSize = fyne.NewSize(800, 600)

func init() {
	MainW.CenterOnScreen()
	MainW.Resize(WindowSize)
}
