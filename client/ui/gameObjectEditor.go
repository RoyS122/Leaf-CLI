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
	"fyne.io/fyne/v2/widget"
)

func buildGameObjectEditor(p models.Project, w fyne.Window, goj models.GameObject) fyne.CanvasObject {
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

		utils.OpenLuaExternal(models.LoadSettings().CodeEditor, goj.Scripts[scriptSelect.Selected])
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
			sprName := selectField.Selected
			for _, spr := range sprites {
				if spr.Name == sprName {
					goj.Sprite = filepath.Join(spr.Directory, "info.lsp")
					break
				}
			}
			goj.Save()
		}),
	)
}
