package command

type UpdatePasswordCommand struct {
	ID          uint
	URL         string
	Login       string
	Password    string
	Description string
}
