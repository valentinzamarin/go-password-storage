package main

import (
	"log"
	"net"

	"password-storage/internal/app/services"
	"password-storage/internal/encrypt"

	"password-storage/internal/gui"
	"password-storage/internal/infrastructure/sqlite"
	"password-storage/internal/infrastructure/sqlite/auth"
	passwords "password-storage/internal/infrastructure/sqlite/password"
)

func main() {

	listener, err := net.Listen("tcp", "127.0.0.1:54231")
	if err != nil {
		log.Fatalf("app is already running %v", err)
	}
	defer listener.Close()

	basePath := "./notebook.db"

	db, err := sqlite.NewConnection(basePath)
	if err != nil {
		log.Fatalf("db conn error: %v", err)
	}

	sqlite.Migrate(db)

	/* additional functional */
	encrypt := encrypt.NewPasswordEncrypt()

	/* db repos */
	authRepo := auth.NewAuthRepo(db)
	passwordRepo := passwords.NewGormPasswordRepository(db)

	/* app services */
	authService := services.NewAuthService(authRepo, encrypt)
	passwordService := services.NewPasswordService(passwordRepo, encrypt)

	uiApp := gui.NewApp(passwordService, authService)

	uiApp.Run()
}
