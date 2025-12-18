package request

import (
	"errors"
	"regexp"
	"strings"

	command "password-storage/internal/app/command"
)

const (
	MinLoginLength    = 3
	MaxLoginLength    = 100
	MinPasswordLength = 8
	MaxDescLength     = 500
)

var (
	urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

	ErrURLRequired      = errors.New("URL is required")
	ErrInvalidURL       = errors.New("URL must be valid (http:// or https://)")
	ErrLoginRequired    = errors.New("Login is required")
	ErrLoginTooShort    = errors.New("Login must contain at least 3 characters")
	ErrLoginTooLong     = errors.New("Login must not exceed 100 characters")
	ErrPasswordRequired = errors.New("Password is required")
	ErrPasswordTooWeak  = errors.New("Password must contain at least 8 characters")
	ErrDescTooLong      = errors.New("Description must not exceed 500 characters")
)

type AddPasswordRequest struct {
	URL         string `json:"url"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

func (req *AddPasswordRequest) Sanitize() {
	req.URL = strings.TrimSpace(req.URL)
	req.Login = strings.TrimSpace(req.Login)
	req.Password = strings.TrimSpace(req.Password)
	req.Description = strings.TrimSpace(req.Description)
}

func (req *AddPasswordRequest) Validate() error {

	req.Sanitize()

	if req.URL == "" {
		return ErrURLRequired
	}
	if !urlRegex.MatchString(req.URL) {
		return ErrInvalidURL
	}

	if req.Login == "" {
		return ErrLoginRequired
	}
	if len(req.Login) < MinLoginLength {
		return ErrLoginTooShort
	}
	if len(req.Login) > MaxLoginLength {
		return ErrLoginTooLong
	}

	if req.Password == "" {
		return ErrPasswordRequired
	}
	if len(req.Password) < MinPasswordLength {
		return ErrPasswordTooWeak
	}

	if len(req.Description) > MaxDescLength {
		return ErrDescTooLong
	}

	return nil
}

func (req *AddPasswordRequest) ToAddPasswordCommand() (*command.AddPasswordCommand, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &command.AddPasswordCommand{
		URL:         req.URL,
		Login:       req.Login,
		Password:    req.Password,
		Description: req.Description,
	}, nil
}
