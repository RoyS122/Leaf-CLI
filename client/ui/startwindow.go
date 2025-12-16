package ui

import (
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateWinSelection(a fyne.App) fyne.Window {

	w := a.NewWindow("Project selection")
	w.Resize(fyne.NewSize(900, 600))
	openProj := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			fmt.Println("Error:", err)
			a.Quit()
			return
		}
		if reader == nil {
			fmt.Println("No file selected")
			return
		}

		fmt.Println("Selected:", reader.URI().Path())

		// Read file content
		data, err := io.ReadAll(reader)
		defer reader.Close()

		if err != nil {
			fmt.Println("Error reading file:", err)
			a.Quit()
			return
		}

		fmt.Println("Content:", string(data))

	}, w)

	hello := widget.NewLabel("Hello Fyne!")
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	w.SetOnClosed(a.Quit)
	content := container.NewVBox(
		hello,
		input,
		widget.NewButton("Open Project!", func() {
			openProj.Show()
		}),
		widget.NewButton("New Project", func() {
			CreateNewProjectWindow(a).Show()
		}),
	)
	w.SetContent(content)

	return w
}
