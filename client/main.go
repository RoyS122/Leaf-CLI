package main

import (
	"leafcli/ui"

	"fyne.io/fyne/v2/app"
)

func main() {

	a := app.NewWithID("fr.roys.leaf")

	w := ui.CreateWinSelection(a)
	w.ShowAndRun()
}
