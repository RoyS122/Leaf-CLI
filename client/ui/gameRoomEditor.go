package ui

import (
	"fmt"
	"leafcli/models"
	"leafcli/ui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowAddObjectInstance(
	w fyne.Window,
	lGO []models.GameObject,
	onCreate func(name string),
) {
	listGONames := make([]string, 0, len(lGO))
	for _, goObj := range lGO {
		listGONames = append(listGONames, goObj.Name)
	}

	GameObjectSelect := widget.NewSelect(listGONames, func(selected string) {})
	// GameObjectSelect.SetPlaceHolder("GameObject name")

	dialog.NewForm(
		"New Object Instance",
		"Create",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", GameObjectSelect),
		},
		func(confirm bool) {
			if !confirm {
				return
			}
			if GameObjectSelect.Selected == "" {
				dialog.ShowError(
					fmt.Errorf("asset name cannot be empty"),
					w,
				)
				return
			}
			onCreate(GameObjectSelect.Selected)
		},
		w,
	).Show()
}

func buildRoomEditor(p models.Project, w fyne.Window, room models.Room) fyne.CanvasObject {

	name := widget.NewEntry()
	name.SetText(room.Name)

	// Charger tous les GameObjects disponibles
	gameObjects := models.LoadGameObjects(p.GameObjects)
	var gObjNames []string
	for _, gObj := range gameObjects {
		gObjNames = append(gObjNames, gObj.Name)
	}

	addInstanceButton := widget.NewButton("Add Instance", func() {
		ShowAddObjectInstance(w, gameObjects, func(selected string) {
			var gameobjectslist models.GOList = gameObjects
			room.GOInstances = append(room.GOInstances, models.NewInstanceFromGO(gameobjectslist.GetGOByName(selected), 0, 0))
			fmt.Println("Added instance of:", selected)
			room.Save()
		})
	})

	savebutton := widget.NewButton("Save", func() {
		room.Name = name.Text
		if err := room.Save(); err != nil {
			dialog.ShowError(err, w)
		}
	})

	viewport := widgets.BuildRoomViewport(room)
	form := widget.NewForm(
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Instances", addInstanceButton),
	)
	viewContainer := container.NewGridWrap(
		fyne.NewSize(500, 300),
		viewport,
	)
	centeredView := container.NewCenter(viewContainer)
	return container.NewVBox(
		widget.NewLabelWithStyle("Room", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
		centeredView,
		savebutton,
	)
}
