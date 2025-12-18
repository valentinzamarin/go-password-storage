package entities

import "errors"

type Password struct {
	id          uint
	url         string
	login       string
	password    string
	description string
}

func (p *Password) validate() error {
	if p.url == "" {
		return errors.New("url is required")
	}
	if p.login == "" {
		return errors.New("login is required")
	}
	if p.password == "" {
		return errors.New("password is required")
	}
	return nil
}

func NewPassword(url, login, password, description string) (*Password, error) {
	p := &Password{
		url:         url,
		login:       login,
		password:    password,
		description: description,
	}

	if err := p.validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Password) GetID() uint {
	return p.id
}
func (p *Password) GetURL() string {
	return p.url
}

func (p *Password) GetLogin() string {
	return p.login
}

func (p *Password) GetPassword() string {
	return p.password
}

func (p *Password) GetDescription() string {
	return p.description
}

func (p *Password) SetID(id uint) {
	p.id = id
}

func (p *Password) SetPassword(newPassword string) {
	p.password = newPassword
}
