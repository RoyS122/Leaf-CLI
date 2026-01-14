package ui

import (
	"leafcli/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateSettingsWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("LeafCLI: Engine Settings")
	w.Resize(fyne.NewSize(400, 300))

	settings := models.LoadSettings()

	title := canvas.NewText("Engine Settings", nil)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = 22

	// Input pour dossier par défaut des projets
	projectDirEntry := widget.NewEntry()
	projectDirEntry.SetText(settings.CodeEditor)
	projectDirEntry.SetPlaceHolder("Code editor for script edits")
	selectDirBtn := widget.NewButton("Select code editor", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri != nil {
				projectDirEntry.SetText(uri.Path())
			}
		}, w).Show()

	})

	// Input pour la langue
	languageSelect := widget.NewSelect([]string{"en"}, func(s string) {
		settings.Language = s
	})
	languageSelect.SetSelected(settings.Language)

	// Bouton de sauvegarde
	saveBtn := widget.NewButton("Save Settings", func() {
		settings.CodeEditor = projectDirEntry.Text
		settings.Language = languageSelect.Selected

		if err := models.SaveSettings(&settings); err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Settings Saved", "Engine settings have been saved successfully.", w)
		w.Close()
	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Code editor:"),
		projectDirEntry,
		selectDirBtn,
		widget.NewLabel("Language:"),
		languageSelect,
		saveBtn,
	)

	w.SetContent(content)
	return w
}
