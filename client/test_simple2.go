package main

import (
	"fmt"
	 "io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
)

func main() {
	// Goroutine correcte
	go func() { 
		fmt.Println("this is a test")
	}()
	
	a := app.New()
	w := a.NewWindow("Hello world")
	w.Resize(fyne.NewSize(900, 600))
	openProj := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
        if err != nil {
            fmt.Println("Error:", err)
            a.Quit()
            return
        }
        if reader == nil {
            fmt.Println("No file selected")
            return
        }

        fmt.Println("Selected:", reader.URI().Path())

        // Read file content
        data, err := io.ReadAll(reader)
        defer reader.Close()

        if err != nil {
            fmt.Println("Error reading file:", err)
            a.Quit()
            return
        }

        fmt.Println("Content:", string(data))

       

    }, w)

	hello := widget.NewLabel("Hello Fyne!")
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	

	content := container.NewVBox(
		hello,
		input,
		widget.NewButton("Open Project!", func() {
			openProj.Show();
		}),
	)
	w.SetContent(content)
	
	w.ShowAndRun()  
}