package request

import (
	"errors"
	"regexp"
	"strings"

	command "password-storage/internal/app/command"
)

const (
	MinLoginLengthUpdate    = 3
	MaxLoginLengthUpdate    = 100
	MinPasswordLengthUpdate = 8
	MaxDescLengthUpdate     = 500
)

var (
	urlRegexUpdate = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

	ErrURLRequiredUpdate      = errors.New("URL is required")
	ErrInvalidURLUpdate       = errors.New("URL must be valid (http:// or https://)")
	ErrLoginRequiredUpdate    = errors.New("Login is required")
	ErrLoginTooShortUpdate    = errors.New("Login must contain at least 3 characters")
	ErrLoginTooLongUpdate     = errors.New("Login must not exceed 100 characters")
	ErrPasswordRequiredUpdate = errors.New("Password is required")
	ErrPasswordTooWeakUpdate  = errors.New("Password must contain at least 8 characters")
	ErrDescTooLongUpdate      = errors.New("Description must not exceed 500 characters")
)

type UpdatePasswordRequest struct {
	ID          uint   `json:"id"`
	URL         string `json:"url"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

func (req *UpdatePasswordRequest) Sanitize() {
	req.URL = strings.TrimSpace(req.URL)
	req.Login = strings.TrimSpace(req.Login)
	req.Password = strings.TrimSpace(req.Password)
	req.Description = strings.TrimSpace(req.Description)
}

func (req *UpdatePasswordRequest) Validate() error {
	req.Sanitize()

	if req.URL == "" {
		return ErrURLRequiredUpdate
	}
	if !urlRegexUpdate.MatchString(req.URL) {
		return ErrInvalidURLUpdate
	}

	if req.Login == "" {
		return ErrLoginRequiredUpdate
	}
	if len(req.Login) < MinLoginLengthUpdate {
		return ErrLoginTooShortUpdate
	}
	if len(req.Login) > MaxLoginLengthUpdate {
		return ErrLoginTooLongUpdate
	}

	if req.Password == "" {
		return ErrPasswordRequiredUpdate
	}
	if len(req.Password) < MinPasswordLengthUpdate {
		return ErrPasswordTooWeakUpdate
	}

	if len(req.Description) > MaxDescLengthUpdate {
		return ErrDescTooLongUpdate
	}

	return nil
}

func (req *UpdatePasswordRequest) ToUpdatePasswordCommand() (*command.UpdatePasswordCommand, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &command.UpdatePasswordCommand{
		ID:          req.ID,
		URL:         req.URL,
		Login:       req.Login,
		Password:    req.Password,
		Description: req.Description,
	}, nil
}
