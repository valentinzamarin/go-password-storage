package components

import (
	"fyne.io/fyne/v2/widget"
)

func CreateInputField(placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	return entry
}
