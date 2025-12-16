package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateNewProjectWindow(a fyne.App) fyne.Window {
	baseSize := fyne.NewSize(300, 500)
	w := a.NewWindow("New Project")
	w.Resize(baseSize)

	title := canvas.NewText("Create a new project!", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24 // <-- la taille que tu veux

	inputName := widget.NewEntry()
	inputName.SetPlaceHolder("Project name")

	dirLabel := widget.NewLabel("No folder selected")
	selectDirButton := widget.NewButton("Select folder", func() {
		w.Resize(fyne.NewSize(600, 400))
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			w.Resize(baseSize)
			if err != nil {
				fmt.Println("Error:", err)

				return
			}
			if uri == nil {
				dirLabel.SetText("No folder selected")
				return
			}

			dirLabel.SetText(uri.Path()) // affiche le dossier choisi
		}, w).Show()
	})
	content := container.NewVBox(
		title,
		inputName,
		inputDescription,
		selectDirButton,
		widget.NewButton("Create project!", func() {
			fmt.Println("Create a project,", inputName.Text)
		}),
	)
	w.SetContent(content)

	return w
}
