package ui

import (
	"fmt"
	"leafcli/models"
	"leafcli/utils"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
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

	w := a.NewWindow("LeafCLI: " + p.Name + " " + p.Version.String())
	w.Resize(fyne.NewSize(1000, 600))

	var selected struct {
		cat string
		id  int
	}

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
	)

	w.SetMainMenu(menu)
	editorPanel := container.NewVBox(
		widget.NewLabel("No selection"),
	)
	buildGameObjectEditor := func(goj models.GameObject) fyne.CanvasObject {
		name := widget.NewEntry()
		name.SetText(goj.Name)

		sprites := models.LoadSprites(p.Sprites)
		var nSPs []string
		for _, nSP := range sprites {
			nSPs = append(nSPs, nSP.Name)
		}
		selectField := widget.NewSelect(nSPs, func(s string) {
			fmt.Println("Selected:", s)
		})
		scripts := make([]string, 0, len(goj.Scripts))
		for sname := range goj.Scripts {
			scripts = append(scripts, sname)
		} // liste des scripts pour le GameObject
		scriptSelect := widget.NewSelect(scripts, func(selected string) {
			// OpenLuaEditor(filepath.Join(goObj.Directory, selected))
		})

		addBtn := widget.NewButton("+", func() {
			dialog.ShowEntryDialog("New Script", "Script name:", func(name string) {
				if name == "" {
					return
				}
				// crée le fichier .lua
				filePath := filepath.Join(goj.Directory, name+".lua")
				os.WriteFile(filePath, []byte("-- new script"), 0644)

				// ajoute à la liste et refresh le select
				goj.Scripts[name] = filePath
				scripts = append(scripts, name)
				scriptSelect.Options = scripts
				scriptSelect.Refresh()
				goj.Save()
			}, w)
		})

		editButton := widget.NewButton("Edit", func() {
			if scriptSelect.Selected == "" {
				return
			}
			utils.OpenLuaExternal("C:/Users/morga/AppData/Local/Programs/Microsoft VS Code/Code.exe", goj.Scripts[scriptSelect.Selected])
		})
		editButton.Disable() // désactivé par défaut

		// Sélection du script
		scriptSelect.OnChanged = func(selected string) {
			if selected == "" {
				editButton.Disable()
			} else {
				editButton.Enable()
			}
		}

		// HBox = toujours bouton présent
		scriptRow := container.NewHBox(
			scriptSelect,
			editButton,
			addBtn,
		)

		return container.NewVBox(
			widget.NewLabelWithStyle("Game Object", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewForm(
				widget.NewFormItem("Name", name),
				widget.NewFormItem("Sprite", selectField),
				widget.NewFormItem("Scripts", scriptRow),
			),
			widget.NewButton("Save", func() {
				goj.Name = name.Text
				goj.Save()
			}),
		)
	}

	var assetList *widget.List

	createCategory := func(name string) fyne.CanvasObject {

		switch name {
		case "GameObjects":

			assetList = widget.NewList(
				func() int { return len(p.GameObjects) },
				func() fyne.CanvasObject { return widget.NewLabel("Asset") },
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(models.LoadGameObjects(p.GameObjects)[i].Name)
				},
			)

		case "Sprites":
			assetList = widget.NewList(
				func() int { return len(p.Sprites) },
				func() fyne.CanvasObject { return widget.NewLabel("Asset") },
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(models.LoadSprites(p.Sprites)[i].Name)
				},
			)

		case "Rooms":
			assetList = widget.NewList(
				func() int { return len(p.Rooms) },
				func() fyne.CanvasObject { return widget.NewLabel("Asset") },
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(models.LoadRooms(p.Rooms)[i].Name)

				},
			)
		}
		assetList.OnSelected = func(id widget.ListItemID) {
			selected.cat = name
			selected.id = id
			editorPanel.Objects = nil // clear

			switch selected.cat {
			case "GameObjects":
				var gO models.GameObject = models.LoadGameObjects(p.GameObjects)[selected.id]
				editorPanel.Add(buildGameObjectEditor(gO))
			}
			editorPanel.Refresh()
		}
		scroll := container.NewVScroll(assetList)
		scroll.SetMinSize(fyne.NewSize(200, 120)) // hauteur par défaut

		scroll.Hide() // Liste cachée au départ

		// Header = Label à gauche, boutons à droite
		addBtn := widget.NewButton("+", func() {
			ShowCreateAssetDialog(w, name, func(assetName string) {
				switch name {
				case "Sprites":
					var sp models.Sprite
					sp.Name = assetName
					sp.Directory = filepath.Join(p.Directory, "Sprites")
					utils.EnsureDir(sp.Directory)
					p.Sprites = append(p.Sprites, filepath.Join(sp.Directory, sp.Name+".lsp"))

					sp.Save()
					fmt.Println("Create sprite:", assetName)

				case "GameObjects":
					var gO models.GameObject
					gO.Name = assetName

					gO.Directory = filepath.Join(p.Directory, "GameObjects", gO.Name)
					utils.EnsureDir(gO.Directory)
					p.GameObjects = append(p.GameObjects, filepath.Join(gO.Directory, "info.lgo"))
					gO.Scripts = make(map[string]string)
					gO.Save()

					fmt.Println("Create game object:", assetName)

				case "Rooms":
					var roo models.Room
					roo.Name = assetName
					roo.Directory = filepath.Join(p.Directory, "Rooms")
					utils.EnsureDir(roo.Directory)
					p.Rooms = append(p.Rooms, filepath.Join(roo.Directory, roo.Name+".lro"))
					roo.Save()
					fmt.Println("Create room:", assetName)
				}

				p.UpdateProject()
			})
		})
		toggleBtn := widget.NewButton("▾", func() {
			if name == selected.cat {
				selected.cat = ""
				selected.id = 0
			}
			if scroll.Visible() {
				scroll.Hide()
			} else {
				scroll.Show()
			}
		})

		header := container.NewHBox(
			widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(), // pousse les boutons à droite
			addBtn,
			toggleBtn,
		)

		return container.NewVBox(header, scroll)
	}

	// ---- Categories
	spritesCat := createCategory("Sprites")
	gameObjectsCat := createCategory("GameObjects")
	roomsCat := createCategory("Rooms")

	leftPanel := container.NewVBox(
		spritesCat,
		gameObjectsCat,
		roomsCat,
	)
	// editorPanel := container.NewVBox()

	// ---- Layout principal
	content := container.NewHSplit(
		leftPanel,
		container.NewBorder(
			nil, nil, nil, nil,
			editorPanel,
		),
	)
	content.SetOffset(0.20)

	w.SetContent(content)
	return w
}
