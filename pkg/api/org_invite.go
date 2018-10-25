package api

import (
	"fmt"

	"github.com/miprotek/grafana-de/pkg/api/dtos"
	"github.com/miprotek/grafana-de/pkg/bus"
	"github.com/miprotek/grafana-de/pkg/events"
	"github.com/miprotek/grafana-de/pkg/metrics"
	m "github.com/miprotek/grafana-de/pkg/models"
	"github.com/miprotek/grafana-de/pkg/setting"
	"github.com/miprotek/grafana-de/pkg/util"
)

func GetPendingOrgInvites(c *m.ReqContext) Response {
	query := m.GetTempUsersQuery{OrgId: c.OrgId, Status: m.TmpUserInvitePending}

	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Empfang der Einladung von db fehlgeschlagen", err)
	}

	for _, invite := range query.Result {
		invite.Url = setting.ToAbsUrl("invite/" + invite.Code)
	}

	return JSON(200, query.Result)
}

func AddOrgInvite(c *m.ReqContext, inviteDto dtos.AddInviteForm) Response {
	if !inviteDto.Role.IsValid() {
		return Error(400, "Ungültige Rolle angegeben", nil)
	}

	// first try get existing user
	userQuery := m.GetUserByLoginQuery{LoginOrEmail: inviteDto.LoginOrEmail}
	if err := bus.Dispatch(&userQuery); err != nil {
		if err != m.ErrUserNotFound {
			return Error(500, "Datenbankabfrage für vorhandene Benutzerprüfung fehlgeschlagen", err)
		}

		if setting.DisableLoginForm {
			return Error(401, "Benutzer konnte nicht gefunden werden", nil)
		}
	} else {
		return inviteExistingUserToOrg(c, userQuery.Result, &inviteDto)
	}

	cmd := m.CreateTempUserCommand{}
	cmd.OrgId = c.OrgId
	cmd.Email = inviteDto.LoginOrEmail
	cmd.Name = inviteDto.Name
	cmd.Status = m.TmpUserInvitePending
	cmd.InvitedByUserId = c.UserId
	cmd.Code = util.GetRandomString(30)
	cmd.Role = inviteDto.Role
	cmd.RemoteAddr = c.Req.RemoteAddr

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Fehler beim speichern der Einladung in die Datenbank", err)
	}

	// send invite email
	if inviteDto.SendEmail && util.IsEmail(inviteDto.LoginOrEmail) {
		emailCmd := m.SendEmailCommand{
			To:       []string{inviteDto.LoginOrEmail},
			Template: "new_user_invite.html",
			Data: map[string]interface{}{
				"Name":      util.StringsFallback2(cmd.Name, cmd.Email),
				"OrgName":   c.OrgName,
				"Email":     c.Email,
				"LinkUrl":   setting.ToAbsUrl("invite/" + cmd.Code),
				"InvitedBy": util.StringsFallback3(c.Name, c.Email, c.Login),
			},
		}

		if err := bus.Dispatch(&emailCmd); err != nil {
			if err == m.ErrSmtpNotEnabled {
				return Error(412, err.Error(), err)
			}
			return Error(500, "Senden der E-Mail Einladung fehlgeschlagen", err)
		}

		emailSentCmd := m.UpdateTempUserWithEmailSentCommand{Code: cmd.Result.Code}
		if err := bus.Dispatch(&emailSentCmd); err != nil {
			return Error(500, "Fehler beim aktualisieren der Einladung mit den gesendeten E-Mail Informationen", err)
		}

		return Success(fmt.Sprintf("Einladen gesendet an %s", inviteDto.LoginOrEmail))
	}

	return Success(fmt.Sprintf("Einladung erstellt für %s", inviteDto.LoginOrEmail))
}

func inviteExistingUserToOrg(c *m.ReqContext, user *m.User, inviteDto *dtos.AddInviteForm) Response {
	// user exists, add org role
	createOrgUserCmd := m.AddOrgUserCommand{OrgId: c.OrgId, UserId: user.Id, Role: inviteDto.Role}
	if err := bus.Dispatch(&createOrgUserCmd); err != nil {
		if err == m.ErrOrgUserAlreadyAdded {
			return Error(412, fmt.Sprintf("Benutzer %s ist bereits zur Organisation hinzugefügt", inviteDto.LoginOrEmail), err)
		}
		return Error(500, "Fehler beim Versuch, einen Organisationsbenutzer zu erstellen", err)
	}

	if inviteDto.SendEmail && util.IsEmail(user.Email) {
		emailCmd := m.SendEmailCommand{
			To:       []string{user.Email},
			Template: "invited_to_org.html",
			Data: map[string]interface{}{
				"Name":      user.NameOrFallback(),
				"OrgName":   c.OrgName,
				"InvitedBy": util.StringsFallback3(c.Name, c.Email, c.Login),
			},
		}

		if err := bus.Dispatch(&emailCmd); err != nil {
			return Error(500, "E-Mail invited_to_org konnte nicht gesendet werden", err)
		}
	}

	return Success(fmt.Sprintf("Existierender Grafana Benutzer %s hinzugefügt zu Organisation %s", user.NameOrFallback(), c.OrgName))
}

func RevokeInvite(c *m.ReqContext) Response {
	if ok, rsp := updateTempUserStatus(c.Params(":code"), m.TmpUserRevoked); !ok {
		return rsp
	}

	return Success("Einladung widerrufen")
}

func GetInviteInfoByCode(c *m.ReqContext) Response {
	query := m.GetTempUserByCodeQuery{Code: c.Params(":code")}

	if err := bus.Dispatch(&query); err != nil {
		if err == m.ErrTempUserNotFound {
			return Error(404, "Einladung nicht gefunden", nil)
		}
		return Error(500, "Einladung konnte nicht abgerufen werden", err)
	}

	invite := query.Result

	return JSON(200, dtos.InviteInfo{
		Email:     invite.Email,
		Name:      invite.Name,
		Username:  invite.Email,
		InvitedBy: util.StringsFallback3(invite.InvitedByName, invite.InvitedByLogin, invite.InvitedByEmail),
	})
}

func CompleteInvite(c *m.ReqContext, completeInvite dtos.CompleteInviteForm) Response {
	query := m.GetTempUserByCodeQuery{Code: completeInvite.InviteCode}

	if err := bus.Dispatch(&query); err != nil {
		if err == m.ErrTempUserNotFound {
			return Error(404, "Einladung nicht gefunden", nil)
		}
		return Error(500, "Einladung konnte nicht abgerufen werden", err)
	}

	invite := query.Result
	if invite.Status != m.TmpUserInvitePending {
		return Error(412, fmt.Sprintf("Einladung kann nicht im Status %s verwendet werden", invite.Status), nil)
	}

	cmd := m.CreateUserCommand{
		Email:        completeInvite.Email,
		Name:         completeInvite.Name,
		Login:        completeInvite.Username,
		Password:     completeInvite.Password,
		SkipOrgSetup: true,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Benutzererstellung fehlgeschlagen", err)
	}

	user := &cmd.Result

	bus.Publish(&events.SignUpCompleted{
		Name:  user.NameOrFallback(),
		Email: user.Email,
	})

	if ok, rsp := applyUserInvite(user, invite, true); !ok {
		return rsp
	}

	loginUserWithUser(user, c)

	metrics.M_Api_User_SignUpCompleted.Inc()
	metrics.M_Api_User_SignUpInvite.Inc()

	return Success("User created and logged in")
}

func updateTempUserStatus(code string, status m.TempUserStatus) (bool, Response) {
	// update temp user status
	updateTmpUserCmd := m.UpdateTempUserStatusCommand{Code: code, Status: status}
	if err := bus.Dispatch(&updateTmpUserCmd); err != nil {
		return false, Error(500, "Fehler beim aktualisieren des Einladungsstatus", err)
	}

	return true, nil
}

func applyUserInvite(user *m.User, invite *m.TempUserDTO, setActive bool) (bool, Response) {
	// add to org
	addOrgUserCmd := m.AddOrgUserCommand{OrgId: invite.OrgId, UserId: user.Id, Role: invite.Role}
	if err := bus.Dispatch(&addOrgUserCmd); err != nil {
		if err != m.ErrOrgUserAlreadyAdded {
			return false, Error(500, "Fehler beim Versuch, einen Organisationsbenutzer zu erstellen", err)
		}
	}

	// update temp user status
	if ok, rsp := updateTempUserStatus(invite.Code, m.TmpUserCompleted); !ok {
		return false, rsp
	}

	if setActive {
		// set org to active
		if err := bus.Dispatch(&m.SetUsingOrgCommand{OrgId: invite.OrgId, UserId: user.Id}); err != nil {
			return false, Error(500, "Organisation als aktiv festzulegen ist fehlgeschlagen", err)
		}
	}

	return true, nil
}
