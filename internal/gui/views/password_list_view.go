package views

import (
	"fmt"
	command "password-storage/internal/app/command"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/gui/views/dto/mapper"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type GetPasswordsView struct {
	passwordService interfaces.PasswordService
	window          fyne.Window
}

func NewGetPasswordsView(passwordService interfaces.PasswordService, window fyne.Window) *GetPasswordsView {
	return &GetPasswordsView{
		passwordService: passwordService,
		window:          window,
	}
}

func (v *GetPasswordsView) Render() fyne.CanvasObject {
	pwsQueryResults, err := v.passwordService.GetAllPasswords()
	if err != nil {
		dialog.ShowError(err, v.window)
		return widget.NewLabel("Ошибка загрузки паролей")
	}

	passwords := mapper.ToPasswordsResponseList(pwsQueryResults)

	list := widget.NewList(
		func() int {
			return len(passwords)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < len(passwords) {
				pwd := passwords[i]

				box := o.(*fyne.Container)
				box.Objects[0].(*widget.Label).SetText(fmt.Sprintf("URL: %s", pwd.URL))
				box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("Login: %s", pwd.Login))
				box.Objects[2].(*widget.Label).SetText(fmt.Sprintf("Password: %s", pwd.Password))
				box.Objects[3].(*widget.Label).SetText("")
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if int(id) < len(passwords) {
			pwd := passwords[id]

			fields := []struct {
				name  string
				value string
			}{
				{"URL", pwd.URL},
				{"Login", pwd.Login},
				{"Password", pwd.Password},
			}

			buttons := make([]fyne.CanvasObject, 0, len(fields)+1)
			var customDlg dialog.Dialog
			for _, field := range fields {
				f := field
				buttons = append(buttons, widget.NewButton(f.name, func() {
					v.window.Clipboard().SetContent(f.value)
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   f.name + " copied",
						Content: f.value,
					})
				}))
			}

			buttons = append(buttons, widget.NewButton("Delete", func() {
				confirm := dialog.NewConfirm("Delete password", "Are you sure you want to delete this password?", func(confirmed bool) {
					if !confirmed {
						return
					}

					delCmd := &command.DeletePasswordCommand{ID: pwd.ID}
					if err := v.passwordService.DeletePassword(delCmd); err != nil {
						dialog.ShowError(err, v.window)
						return
					}

					if int(id) < len(passwords) {
						passwords = append(passwords[:id], passwords[id+1:]...)
						list.Refresh()
					}

					if customDlg != nil {
						customDlg.Hide()
					}

					dialog.ShowInformation("Deleted", "Password deleted successfully", v.window)
				}, v.window)
				confirm.Show()
			}))

			customDlg = dialog.NewCustom("Click to copy / Delete", "Close", container.NewVBox(buttons...), v.window)
			customDlg.Show()
		}
		list.UnselectAll()
	}

	return container.NewVScroll(list)
}
