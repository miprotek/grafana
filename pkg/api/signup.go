package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/events"
	"github.com/grafana/grafana/pkg/metrics"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

// GET /api/user/signup/options
func GetSignUpOptions(c *m.ReqContext) Response {
	return JSON(200, util.DynMap{
		"verifyEmailEnabled": setting.VerifyEmailEnabled,
		"autoAssignOrg":      setting.AutoAssignOrg,
	})
}

// POST /api/user/signup
func SignUp(c *m.ReqContext, form dtos.SignUpForm) Response {
	if !setting.AllowUserSignUp {
		return Error(401, "Benutzerregistrierungs ist deaktiviert", nil)
	}

	existing := m.GetUserByLoginQuery{LoginOrEmail: form.Email}
	if err := bus.Dispatch(&existing); err == nil {
		return Error(422, "Ein Benutzer mit der gleichen E-Mailadresse existiert bereits", nil)
	}

	cmd := m.CreateTempUserCommand{}
	cmd.OrgId = -1
	cmd.Email = form.Email
	cmd.Status = m.TmpUserSignUpStarted
	cmd.InvitedByUserId = c.UserId
	cmd.Code = util.GetRandomString(20)
	cmd.RemoteAddr = c.Req.RemoteAddr

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Fehler beim erstellen der Registrierung", err)
	}

	bus.Publish(&events.SignUpStarted{
		Email: form.Email,
		Code:  cmd.Code,
	})

	metrics.M_Api_User_SignUpStarted.Inc()

	return JSON(200, util.DynMap{"status": "SignUpCreated"})
}

func SignUpStep2(c *m.ReqContext, form dtos.SignUpStep2Form) Response {
	if !setting.AllowUserSignUp {
		return Error(401, "Benutzerregistrierung ist deaktivert", nil)
	}

	createUserCmd := m.CreateUserCommand{
		Email:    form.Email,
		Login:    form.Username,
		Name:     form.Name,
		Password: form.Password,
		OrgName:  form.OrgName,
	}

	// verify email
	if setting.VerifyEmailEnabled {
		if ok, rsp := verifyUserSignUpEmail(form.Email, form.Code); !ok {
			return rsp
		}
		createUserCmd.EmailVerified = true
	}

	// check if user exists
	existing := m.GetUserByLoginQuery{LoginOrEmail: form.Email}
	if err := bus.Dispatch(&existing); err == nil {
		return Error(401, "Ein Benutzer mit der gleichen E-Mailadresse existiert bereits", nil)
	}

	// dispatch create command
	if err := bus.Dispatch(&createUserCmd); err != nil {
		return Error(500, "Fehler beim erstellen des Benutzers", err)
	}

	// publish signup event
	user := &createUserCmd.Result
	bus.Publish(&events.SignUpCompleted{
		Email: user.Email,
		Name:  user.NameOrFallback(),
	})

	// mark temp user as completed
	if ok, rsp := updateTempUserStatus(form.Code, m.TmpUserCompleted); !ok {
		return rsp
	}

	// check for pending invites
	invitesQuery := m.GetTempUsersQuery{Email: form.Email, Status: m.TmpUserInvitePending}
	if err := bus.Dispatch(&invitesQuery); err != nil {
		return Error(500, "Failed to query database for invites", err)
	}

	apiResponse := util.DynMap{"message": "Benutzerregistrierung war erfolgreich", "code": "redirect-to-landing-page"}
	for _, invite := range invitesQuery.Result {
		if ok, rsp := applyUserInvite(user, invite, false); !ok {
			return rsp
		}
		apiResponse["code"] = "redirect-to-select-org"
	}

	loginUserWithUser(user, c)
	metrics.M_Api_User_SignUpCompleted.Inc()

	return JSON(200, apiResponse)
}

func verifyUserSignUpEmail(email string, code string) (bool, Response) {
	query := m.GetTempUserByCodeQuery{Code: code}

	if err := bus.Dispatch(&query); err != nil {
		if err == m.ErrTempUserNotFound {
			return false, Error(404, "Ungültiger E-Mailbestätigungscode", nil)
		}
		return false, Error(500, "Der temporäre Benutzer konnte nicht gelesen werden", err)
	}

	tempUser := query.Result
	if tempUser.Email != email {
		return false, Error(404, "E-Mailbestätigungscode stimmt nicht mit der E-Mail überein", nil)
	}

	return true, nil
}
