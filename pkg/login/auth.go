package login

import (
	"errors"

	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

var (
	ErrInvalidCredentials   = errors.New("Ungültiger Benutzername oder Passwort")
	ErrTooManyLoginAttempts = errors.New("Zuviele falsche Anmeldeversuche. Login für Benutzer temporär geblockt")
)

type LoginUserQuery struct {
	Username  string
	Password  string
	User      *m.User
	IpAddress string
}

func Init() {
	bus.AddHandler("auth", AuthenticateUser)
	loadLdapConfig()
}

func AuthenticateUser(query *LoginUserQuery) error {
	if err := validateLoginAttempts(query.Username); err != nil {
		return err
	}

	err := loginUsingGrafanaDB(query)
	if err == nil || (err != m.ErrUserNotFound && err != ErrInvalidCredentials) {
		return err
	}

	ldapEnabled, ldapErr := loginUsingLdap(query)
	if ldapEnabled {
		if ldapErr == nil || ldapErr != ErrInvalidCredentials {
			return ldapErr
		}

		err = ldapErr
	}

	if err == ErrInvalidCredentials {
		saveInvalidLoginAttempt(query)
	}

	if err == m.ErrUserNotFound {
		return ErrInvalidCredentials
	}

	return err
}
