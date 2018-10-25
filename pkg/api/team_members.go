package api

import (
	"github.com/miprotek/grafana-de/pkg/api/dtos"
	"github.com/miprotek/grafana-de/pkg/bus"
	m "github.com/miprotek/grafana-de/pkg/models"
	"github.com/miprotek/grafana-de/pkg/setting"
	"github.com/miprotek/grafana-de/pkg/util"
)

// GET /api/teams/:teamId/members
func GetTeamMembers(c *m.ReqContext) Response {
	query := m.GetTeamMembersQuery{OrgId: c.OrgId, TeamId: c.ParamsInt64(":teamId")}

	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Fehler beim abrufen von Teammitgliedern", err)
	}

	for _, member := range query.Result {
		member.AvatarUrl = dtos.GetGravatarUrl(member.Email)
		member.Labels = []string{}

		if setting.IsEnterprise && setting.LdapEnabled && member.External {
			member.Labels = append(member.Labels, "LDAP")
		}
	}

	return JSON(200, query.Result)
}

// POST /api/teams/:teamId/members
func AddTeamMember(c *m.ReqContext, cmd m.AddTeamMemberCommand) Response {
	cmd.TeamId = c.ParamsInt64(":teamId")
	cmd.OrgId = c.OrgId

	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrTeamNotFound {
			return Error(404, "Team nicht gefunden", nil)
		}

		if err == m.ErrTeamMemberAlreadyAdded {
			return Error(400, "Benutzer ist diesem Team bereits hinzugefügt", nil)
		}

		return Error(500, "Fehler beim hinzufügen eines Mitglieds zum Team", err)
	}

	return JSON(200, &util.DynMap{
		"message": "Mitglied zum Team hinzugefügt",
	})
}

// DELETE /api/teams/:teamId/members/:userId
func RemoveTeamMember(c *m.ReqContext) Response {
	if err := bus.Dispatch(&m.RemoveTeamMemberCommand{OrgId: c.OrgId, TeamId: c.ParamsInt64(":teamId"), UserId: c.ParamsInt64(":userId")}); err != nil {
		if err == m.ErrTeamNotFound {
			return Error(404, "Team nicht gefunden", nil)
		}

		if err == m.ErrTeamMemberNotFound {
			return Error(404, "Teammitglied nicht gefunden", nil)
		}

		return Error(500, "Mitglied vom Team entfernen fehlgeschlagen", err)
	}
	return Success("Teammitglied entfernt")
}
