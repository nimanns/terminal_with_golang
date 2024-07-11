package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Simple GUI App")

	label := widget.NewLabel("Hello, World!")
	button := widget.NewButton("Click Me!", func() {
		label.SetText("Button clicked!")
	})

	content := container.NewVBox(
		label,
		button,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(200, 100))
	myWindow.ShowAndRun()
}
