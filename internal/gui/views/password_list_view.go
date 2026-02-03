package views

import (
	"context"
	"fmt"
	command "password-storage/internal/app/command"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/gui/events"
	"password-storage/internal/gui/views/dto/mapper"
	req "password-storage/internal/gui/views/dto/request"
	response "password-storage/internal/gui/views/dto/response"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type GetPasswordsView struct {
	passwordService interfaces.PasswordService
	window          fyne.Window
	list            *widget.List
	passwords       []*response.PasswordsResponse // filtered (displayed)
	allPasswords    []*response.PasswordsResponse // full set
	searchEntry     *widget.Entry
	unsub           func()
}

func NewGetPasswordsView(passwordService interfaces.PasswordService, window fyne.Window) *GetPasswordsView {
	return &GetPasswordsView{
		passwordService: passwordService,
		window:          window,
		list:            nil,
		passwords:       nil,
	}
}

func (v *GetPasswordsView) loadPasswords() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pwsQueryResults, err := v.passwordService.GetAllPasswords(ctx)
	if err != nil {
		return err
	}

	v.allPasswords = mapper.ToPasswordsResponseList(pwsQueryResults)
	// apply current filter (or empty) to populate v.passwords
	v.applyFilter("")
	return nil
}

func (v *GetPasswordsView) Refresh() error {
	if err := v.loadPasswords(); err != nil {
		return err
	}
	if v.list != nil {
		v.list.Refresh()
	}
	return nil
}

// applyFilter filters v.allPasswords by URL containing the query (case-insensitive)
func (v *GetPasswordsView) applyFilter(query string) {
	if query == "" {
		v.passwords = v.allPasswords
		return
	}
	q := strings.ToLower(query)
	filtered := make([]*response.PasswordsResponse, 0, len(v.allPasswords))
	for _, p := range v.allPasswords {
		if strings.Contains(strings.ToLower(p.URL), q) {
			filtered = append(filtered, p)
		}
	}
	v.passwords = filtered
}

func (v *GetPasswordsView) Render() fyne.CanvasObject {
	if err := v.loadPasswords(); err != nil {
		dialog.ShowError(err, v.window)
		return widget.NewLabel("Ошибка загрузки паролей")
	}

	// create search entry
	if v.searchEntry == nil {
		v.searchEntry = widget.NewEntry()
		v.searchEntry.SetPlaceHolder("Search by URL")
		v.searchEntry.OnChanged = func(text string) {
			// filter and refresh on UI thread
			fyne.Do(func() {
				v.applyFilter(text)
				if v.list != nil {
					v.list.Refresh()
				}
			})
		}
	}

	v.list = widget.NewList(
		func() int {
			return len(v.passwords)
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
			if i < len(v.passwords) {
				pwd := v.passwords[i]

				box := o.(*fyne.Container)
				box.Objects[0].(*widget.Label).SetText(fmt.Sprintf("URL: %s", pwd.URL))
				box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("Login: %s", pwd.Login))
				box.Objects[2].(*widget.Label).SetText(fmt.Sprintf("Password: %s", pwd.Password))
				box.Objects[3].(*widget.Label).SetText("")
			}
		},
	)

	v.list.OnSelected = func(id widget.ListItemID) {
		if int(id) < len(v.passwords) {
			pwd := v.passwords[id]

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

					// refresh entire list to keep consistency
					if err := v.Refresh(); err != nil {
						dialog.ShowError(err, v.window)
						return
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

					// refresh entire list
					if err := v.Refresh(); err != nil {
						dialog.ShowError(err, v.window)
						return
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
		v.list.UnselectAll()
	}

	// subscribe to password added events (perform UI updates via fyne.Do)
	if v.unsub == nil {
		ch, unsub := events.Subscribe("password.added")
		v.unsub = unsub

		go func() {
			for range ch {
				fyne.Do(func() {
					if err := v.Refresh(); err != nil {
						dialog.ShowError(err, v.window)
					}
				})
			}
		}()

		// ensure we unsubscribe when the window is closed to avoid leaks
		v.window.SetOnClosed(func() {
			if v.unsub != nil {
				v.unsub()
				v.unsub = nil
			}
		})
	}

	return container.NewBorder(v.searchEntry, nil, nil, nil, container.NewVScroll(v.list))
}
