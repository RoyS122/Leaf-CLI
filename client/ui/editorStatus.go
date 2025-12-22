package ui

import "fyne.io/fyne/v2"

type EditorStatus struct {
	editorPanel fyne.Container
	selected    struct {
		cat string
		id  int
	}
}
