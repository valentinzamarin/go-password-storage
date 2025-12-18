package views

import (
	"log"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/gui/views/dto/request"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type SetupView struct {
	authService interfaces.AuthService
	window      fyne.Window
	onSuccess   func()
}

func NewSetupView(authService interfaces.AuthService, window fyne.Window, onSuccess func()) *SetupView {
	return &SetupView{
		authService: authService,
		window:      window,
		onSuccess:   onSuccess,
	}
}

func (v *SetupView) Show() {
	passwordEntry := widget.NewPasswordEntry()
	confirmEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Create a strong master password...")
	confirmEntry.SetPlaceHolder("Confirm password...")

	items := []*widget.FormItem{
		widget.NewFormItem("Master Password", passwordEntry),
		widget.NewFormItem("Confirm Password", confirmEntry),
	}

	formDialog := dialog.NewForm("Setup", "Create", "Cancel", items, func(confirmed bool) {
		if !confirmed {
			v.window.Close()
			return
		}

		createReq := &request.CreateMasterPasswordRequest{
			Password:        passwordEntry.Text,
			ConfirmPassword: confirmEntry.Text,
		}

		cmd, err := createReq.ToCreateMasterPasswordCommand()
		if err != nil {
			dialog.ShowError(err, v.window)
			v.Show()
			return
		}

		err = v.authService.CreateMasterPassword(cmd)
		if err != nil {
			log.Println("Failed to create master password:", err)
			dialog.ShowError(err, v.window)
			v.Show()
			return
		}

		v.onSuccess()
	}, v.window)

	formDialog.Resize(fyne.NewSize(450, 200))
	formDialog.Show()
}
