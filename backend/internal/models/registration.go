package models

import (
	"strings"
	"time"
	"tplatform/internal/apperr"
	"tplatform/pkg/validators"
)

type Registration struct {
	ID             int
	Name           string
	Email          string
	PasswordHashed string
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
	IsActive       bool
	ExpiredAt      time.Time
	ActivatedAt    *time.Time
	Token          string
	VerifyCode     string
}

func (r *Registration) IsValidForActivate() bool {
	return r.ActivatedAt == nil && time.Now().Before(r.ExpiredAt) && r.IsActive
}

func (r *Registration) SendEmailConfirm() {
	// TODO реализовать отправку email для подтверждения email
}

type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegistrationRequest) Validate() error {
	op := "RegistrationRequest.Validate"
	if !validators.IsEmailValid(r.Email) {
		return apperr.BadRequest(nil, "не верный формат email", op)
	}

	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	if r.Name == "" {
		r.Name = r.Email
	} else {
		r.Name = strings.TrimSpace(r.Name)
	}

	if len(r.Password) < 5 || len(r.Password) > 15 {
		return apperr.BadRequest(nil, "неверные учетные данные", op)
	}

	return nil
}
