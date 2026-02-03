package views

/*

this file - simple add form with 4 fields

*/
import (
	"context"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/gui/views/components"
	"password-storage/internal/gui/views/dto/request"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type AddPasswordView struct {
	passwordService  interfaces.PasswordService
	window           fyne.Window
	inputURL         *widget.Entry
	inputLogin       *widget.Entry
	inputPassword    *widget.Entry
	inputDescription *widget.Entry
}

func NewAddPasswordView(passwordService interfaces.PasswordService, window fyne.Window) *AddPasswordView {
	return &AddPasswordView{
		passwordService:  passwordService,
		window:           window,
		inputURL:         components.CreateInputField("URL"),
		inputLogin:       components.CreateInputField("Login"),
		inputPassword:    components.CreateInputField("Password"),
		inputDescription: components.CreateInputField("Description"),
	}
}

func (v *AddPasswordView) Render() fyne.CanvasObject {
	submitButton := widget.NewButton("Add password", v.handleSubmit)

	return container.NewVBox(
		v.inputURL,
		v.inputLogin,
		v.inputPassword,
		v.inputDescription,
		submitButton,
	)
}

func (v *AddPasswordView) handleSubmit() {

	addPasswordRequest := request.AddPasswordRequest{
		URL:         v.inputURL.Text,
		Login:       v.inputLogin.Text,
		Password:    v.inputPassword.Text,
		Description: v.inputDescription.Text,
	}
	addPwCmd, err := addPasswordRequest.ToAddPasswordCommand()
	if err != nil {
		dialog.ShowError(err, v.window)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceErr := v.passwordService.AddNewPassword(ctx, addPwCmd)
	if serviceErr != nil {
		dialog.ShowError(serviceErr, v.window)
		return
	}

	dialog.ShowInformation("success", "Password added", v.window)

	v.clearForm()
}

func (v *AddPasswordView) clearForm() {
	v.inputURL.SetText("")
	v.inputLogin.SetText("")
	v.inputPassword.SetText("")
	v.inputDescription.SetText("")
}
