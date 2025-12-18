package ui

import (
	"fmt"
	"image/color"
	"leafcli/models"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowProjectCreationSuccess(w fyne.Window, projectPath string, onClose func()) {
	d := dialog.NewInformation(
		"Project created",
		"Project successfully created at:\n"+projectPath,
		w,
	)

	if onClose != nil {
		d.SetOnClosed(onClose)
	}

	d.Show()
}

func CreateNewProjectWindow(a fyne.App) fyne.Window {
	baseSize := fyne.NewSize(300, 500)
	w := a.NewWindow("LeafCLI: New Project")
	w.Resize(baseSize)

	var new_p models.Project
	var dirPath string

	title := canvas.NewText("Create a new project!", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 24 // <-- la taille que tu veux

	inputName := widget.NewEntry()
	inputName.SetPlaceHolder("Project name")

	inputDescription := widget.NewEntry()
	inputDescription.SetPlaceHolder("Project description")

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
			dirPath = uri.Path()
			dirLabel.SetText(uri.Path()) // affiche le dossier choisi
		}, w).Show()
	})
	createButton := widget.NewButton("Create project", func() {
		new_p.Name = inputName.Text
		new_p.Description = inputDescription.Text
		new_p.Description = inputDescription.Text

		if new_p.Name == "" {
			dialog.ShowError(fmt.Errorf("project name is required"), w)
			return
		}

		if dirPath == "" {
			dialog.ShowError(fmt.Errorf("please select a folder"), w)
			return
		}

		err := new_p.CreateProjectDirectory(dirPath)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		projectPath := filepath.Join(dirPath, new_p.Name)
		ShowProjectCreationSuccess(w, projectPath, func() { w.Close() })

	})

	content := container.NewVBox(
		title,
		inputName,
		inputDescription,
		dirLabel,
		selectDirButton,
		createButton,
	)
	w.SetContent(content)

	return w
}
