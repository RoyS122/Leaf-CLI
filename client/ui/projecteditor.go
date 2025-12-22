package ui

import (
	"fmt"
	"leafcli/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowCreateAssetDialog(
	w fyne.Window,
	assetType string,
	onCreate func(name string),
) {

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Asset name")

	dialog.NewForm(
		"New "+assetType,
		"Create",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
		},
		func(confirm bool) {
			if !confirm {
				return
			}
			if nameEntry.Text == "" {
				dialog.ShowError(
					fmt.Errorf("asset name cannot be empty"),
					w,
				)
				return
			}
			onCreate(nameEntry.Text)
		},
		w,
	).Show()
}

func CreateProjectEditor(a fyne.App, p models.Project) fyne.Window {
	var edStatus EditorStatus
	w := a.NewWindow("LeafCLI: " + p.Name + " " + p.Version.String())
	w.Resize(fyne.NewSize(1000, 600))

	menu := fyne.NewMainMenu(
		fyne.NewMenu(
			"File",
			fyne.NewMenuItem("New Project", func() {
				fmt.Println("New project")
			}),
			fyne.NewMenuItem("Open Project", func() {
				fmt.Println("Open project")
			}),
		),
		fyne.NewMenu(
			"Edit",
			fyne.NewMenuItem("Undo", func() {}),
			fyne.NewMenuItem("Redo", func() {}),
		),
		fyne.NewMenu(
			"Settings",
			fyne.NewMenuItem("Engine Settings", func() {
				wsettings := CreateSettingsWindow(a)
				wsettings.Show()
				wsettings.RequestFocus()
			}),
		),
	)

	w.SetMainMenu(menu)
	editorPanel := container.NewVBox(
		widget.NewLabel("No selection"),
	)
	edStatus.editorPanel = *editorPanel

	var assetList *widget.List

	// ---- Categories
	spritesCat := createCategory(&edStatus, assetList, p, w, "Sprites")
	gameObjectsCat := createCategory(&edStatus, assetList, p, w, "GameObjects")
	roomsCat := createCategory(&edStatus, assetList, p, w, "Rooms")

	leftPanel := container.NewVBox(
		spritesCat,
		gameObjectsCat,
		roomsCat,
	)

	// ---- Layout principal
	content := container.NewHSplit(
		leftPanel,
		container.NewBorder(
			nil, nil, nil, nil,
			&edStatus.editorPanel,
		),
	)
	content.SetOffset(0.20)

	w.SetContent(content)
	return w
}
