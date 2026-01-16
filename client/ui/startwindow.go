package ui

import (
	"fmt"
	"leafcli/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateWinSelection(a fyne.App) fyne.Window {
	var mainWindow bool = true
	w := a.NewWindow("LeafCLI: Project selection")
	w.Resize(fyne.NewSize(900, 600))
	w.SetOnClosed(func() {
		if mainWindow {
			a.Quit()
		}
	})

	title := widget.NewLabelWithStyle(
		"Hello Game Programmer.",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	emptyLabel := widget.NewLabel("Aucun projet récent")
	emptyLabel.Alignment = fyne.TextAlignCenter

	openButton := widget.NewButton("Open a project", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if reader == nil {
				fmt.Println("No file selected")
				return
			}
			defer reader.Close()
			p, err := models.LoadProjectFromFile(reader.URI().Path())
			if err != nil {
				fmt.Println("Failed to load project:", err)
				return
			}
			editorWin := CreateProjectEditor(a, *p)
			editorWin.Show()
			mainWindow = false
			w.Close()

		}, w).Show()
	})

	createButton := widget.NewButton("Create a project", func() {
		np := CreateNewProjectWindow(a)
		np.Show()
		np.RequestFocus()
	})

	createButton.Importance = widget.HighImportance

	buttons := container.NewGridWithColumns(
		2,
		openButton,
		createButton,
	)

	recentsProjectContent := container.NewVBox(
		emptyLabel,
		buttons,
	)

	recentProjectsSection := Section(
		"Recent projects",
		container.NewCenter(recentsProjectContent),
	)

	content := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewCenter(
			container.NewVBox(
				title,
				recentProjectsSection,
			),
		),
	)

	w.SetContent(content)
	return w
}
