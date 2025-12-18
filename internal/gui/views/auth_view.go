// internal/gui/views/login_view.go
package views

import (
	"password-storage/internal/app/interfaces"
	"password-storage/internal/gui/views/dto/request"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type LoginView struct {
	authService interfaces.AuthService
	window      fyne.Window
	onSuccess   func()
}

func NewLoginView(authService interfaces.AuthService, window fyne.Window, onSuccess func()) *LoginView {
	return &LoginView{
		authService: authService,
		window:      window,
		onSuccess:   onSuccess,
	}
}

func (v *LoginView) Show() {
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter master password...")

	formDialog := dialog.NewForm("Unlock", "Unlock", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Master Password", passwordEntry),
	}, func(confirmed bool) {
		if !confirmed {
			v.window.Close()
			return
		}

		authReq := &request.AuthenticateRequest{
			Password: passwordEntry.Text,
		}

		cmd, err := authReq.ToCommand()
		if err != nil {
			dialog.ShowError(err, v.window)
			v.Show()
			return
		}

		err = v.authService.Authenticate(cmd)
		if err != nil {
			dialog.ShowError(err, v.window)
			v.Show()
			return
		}

		v.onSuccess()
	}, v.window)

	formDialog.Resize(fyne.NewSize(400, 150))
	formDialog.Show()
}
