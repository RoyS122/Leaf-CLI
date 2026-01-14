package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Section(title string, content fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle(
			title,
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		widget.NewSeparator(),
		content,
	)
}
