package gui

import (
	"fmt"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/app/query"
	"password-storage/internal/gui/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

type App struct {
	fyneApp         fyne.App
	window          fyne.Window
	passwordService interfaces.PasswordService
	authService     interfaces.AuthService
}

func NewApp(
	passwordService interfaces.PasswordService,
	authService interfaces.AuthService,
) *App {
	a := app.New()
	w := a.NewWindow("Password Storage")

	return &App{
		fyneApp:         a,
		window:          w,
		passwordService: passwordService,
		authService:     authService,
	}
}

func (a *App) Run() {
	a.window.Resize(fyne.NewSize(720, 600))
	a.showLoginOrSetupDialog()
	a.window.ShowAndRun()
}

func (a *App) showLoginOrSetupDialog() {
	q := &query.IsMasterPasswordSetQuery{}
	result, err := a.authService.IsMasterPasswordSet(q)
	if err != nil {
		dialog.ShowError(fmt.Errorf("fatal database error: %w", err), a.window)
		a.window.Close()
		return
	}

	if result.IsSet {
		loginView := views.NewLoginView(a.authService, a.window, a.loadMainView)
		loginView.Show()
	} else {
		setupView := views.NewSetupView(a.authService, a.window, a.loadMainView)
		setupView.Show()
	}
}

func (a *App) loadMainView() {
	passwordList := views.NewGetPasswordsView(a.passwordService, a.window)
	addPasswordView := views.NewAddPasswordView(a.passwordService, a.window)

	tabs := container.NewAppTabs(
		container.NewTabItem("Passwords", passwordList.Render()),
		container.NewTabItem("Add", addPasswordView.Render()),
	)

	windowView := container.NewBorder(nil, nil, nil, nil, tabs)
	a.window.SetContent(windowView)
}
