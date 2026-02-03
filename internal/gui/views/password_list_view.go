package views

import (
	"context"
	"fmt"
	command "password-storage/internal/app/command"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/gui/views/dto/mapper"
	req "password-storage/internal/gui/views/dto/request"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pwsQueryResults, err := v.passwordService.GetAllPasswords(ctx)
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
					ctxDel, cancelDel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancelDel()

					if err := v.passwordService.DeletePassword(ctxDel, delCmd); err != nil {
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

			// Edit button
			buttons = append(buttons, widget.NewButton("Edit", func() {
				// prepare form entries prefilled
				urlEntry := widget.NewEntry()
				urlEntry.SetText(pwd.URL)
				loginEntry := widget.NewEntry()
				loginEntry.SetText(pwd.Login)
				passwordEntry := widget.NewPasswordEntry()
				passwordEntry.SetText(pwd.Password)
				descEntry := widget.NewEntry()
				descEntry.SetText(pwd.Description)

				formItems := []*widget.FormItem{
					widget.NewFormItem("URL", urlEntry),
					widget.NewFormItem("Login", loginEntry),
					widget.NewFormItem("Password", passwordEntry),
					widget.NewFormItem("Description", descEntry),
				}

				var form dialog.Dialog
				form = dialog.NewForm("Edit password", "Save", "Cancel", formItems, func(confirmed bool) {
					if !confirmed {
						return
					}

					updateReq := req.UpdatePasswordRequest{
						ID:          pwd.ID,
						URL:         urlEntry.Text,
						Login:       loginEntry.Text,
						Password:    passwordEntry.Text,
						Description: descEntry.Text,
					}
					cmd, err := updateReq.ToUpdatePasswordCommand()
					if err != nil {
						dialog.ShowError(err, v.window)
						return
					}

					// call service
					ctxUpd, cancelUpd := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancelUpd()

					if err := v.passwordService.UpdatePassword(ctxUpd, cmd); err != nil {
						dialog.ShowError(err, v.window)
						return
					}

					// refresh local list by updating the item in slice and refreshing
					if int(id) < len(passwords) {
						passwords[id].URL = cmd.URL
						passwords[id].Login = cmd.Login
						passwords[id].Password = cmd.Password
						passwords[id].Description = cmd.Description
						list.Refresh()
					}

					// close edit and parent dialogs
					form.Hide()
					if customDlg != nil {
						customDlg.Hide()
					}

					dialog.ShowInformation("Updated", "Password updated successfully", v.window)
				}, v.window)

				form.Resize(fyne.NewSize(400, 200))
				form.Show()
			}))

			customDlg = dialog.NewCustom("Click to copy / Delete", "Close", container.NewVBox(buttons...), v.window)
			customDlg.Show()
		}
		list.UnselectAll()
	}

	return container.NewVScroll(list)
}
