package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/metrics"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/util"
)

func AdminCreateUser(c *m.ReqContext, form dtos.AdminCreateUserForm) {
	cmd := m.CreateUserCommand{
		Login:    form.Login,
		Email:    form.Email,
		Password: form.Password,
		Name:     form.Name,
	}

	if len(cmd.Login) == 0 {
		cmd.Login = cmd.Email
		if len(cmd.Login) == 0 {
			c.JsonApiErr(400, "Validierungsfehler, Sie müssen entweder Benutzername oder E-Mail angeben", nil)
			return
		}
	}

	if len(cmd.Password) < 4 {
		c.JsonApiErr(400, "Passwort fehlt oder ist zu kurz", nil)
		return
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Benutzererstellung ist fehlgeschlagen", err)
		return
	}

	metrics.M_Api_Admin_User_Create.Inc()

	user := cmd.Result

	result := m.UserIdDTO{
		Message: "Benutzer erstellt",
		Id:      user.Id,
	}

	c.JSON(200, result)
}

func AdminUpdateUserPassword(c *m.ReqContext, form dtos.AdminUpdateUserPasswordForm) {
	userID := c.ParamsInt64(":id")

	if len(form.Password) < 4 {
		c.JsonApiErr(400, "Neues Passwort zu kurz", nil)
		return
	}

	userQuery := m.GetUserByIdQuery{Id: userID}

	if err := bus.Dispatch(&userQuery); err != nil {
		c.JsonApiErr(500, "Benutzer konnte nicht aus der Datenbank gelesen werden", err)
		return
	}

	passwordHashed := util.EncodePassword(form.Password, userQuery.Result.Salt)

	cmd := m.ChangeUserPasswordCommand{
		UserId:      userID,
		NewPassword: passwordHashed,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Das Benutzerpasswort konnte nicht aktualisiert werden", err)
		return
	}

	c.JsonOK("Benutzerpasswort wurde aktualisiert")
}

func AdminUpdateUserPermissions(c *m.ReqContext, form dtos.AdminUpdateUserPermissionsForm) {
	userID := c.ParamsInt64(":id")

	cmd := m.UpdateUserPermissionsCommand{
		UserId:         userID,
		IsGrafanaAdmin: form.IsGrafanaAdmin,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Fehler beim Aktualisieren der Benutzerberechtigung", err)
		return
	}

	c.JsonOK("Benutzerberechtigungen wurden geupdatet")
}

func AdminDeleteUser(c *m.ReqContext) {
	userID := c.ParamsInt64(":id")

	cmd := m.DeleteUserCommand{UserId: userID}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Löschen des Benutzers ist fehlgeschlagen", err)
		return
	}

	c.JsonOK("Benutzer erfolgreich gelöscht")
}
