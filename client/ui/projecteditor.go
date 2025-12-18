package ui

import (
	"leafcli/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateProjectEditor(a fyne.App, p models.Project) fyne.Window {
	w := a.NewWindow("LeafCLI: " + p.Name + " " + p.Version.String())
	w.Resize(fyne.NewSize(1000, 600))

	// ---- Categories fixes
	spritesList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("Asset") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// o.(*widget.Label).SetText(p.AssetsByCategory("Sprites")[i].Name)
		},
	)

	gameObjectsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("Asset") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// o.(*widget.Label).SetText(p.AssetsByCategory("GameObjects")[i].Name)
		},
	)

	roomsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("Asset") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// o.(*widget.Label).SetText(p.AssetsByCategory("Rooms")[i].Name)
		},
	)

	// ---- Accordion
	accordion := widget.NewAccordion(
		widget.NewAccordionItem("Sprites", spritesList),
		widget.NewAccordionItem("GameObjects", gameObjectsList),
		widget.NewAccordionItem("Rooms", roomsList),
	)

	// ---- Layout principal
	content := container.NewHSplit(
		accordion,
		widget.NewLabel("Editor area"), // placeholder pour la zone de travail
	)
	content.SetOffset(0.25)

	w.SetContent(content)
	return w
}
