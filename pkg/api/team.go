package api

import (
	"github.com/miprotek/grafana-de/pkg/api/dtos"
	"github.com/miprotek/grafana-de/pkg/bus"
	m "github.com/miprotek/grafana-de/pkg/models"
	"github.com/miprotek/grafana-de/pkg/util"
)

// POST /api/teams
func CreateTeam(c *m.ReqContext, cmd m.CreateTeamCommand) Response {
	cmd.OrgId = c.OrgId
	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrTeamNameTaken {
			return Error(409, "Teamname bereits in benutzung", err)
		}
		return Error(500, "Erstellen des Teams ist fehlgeschlagen", err)
	}

	return JSON(200, &util.DynMap{
		"teamId":  cmd.Result.Id,
		"message": "Team created",
	})
}

// PUT /api/teams/:teamId
func UpdateTeam(c *m.ReqContext, cmd m.UpdateTeamCommand) Response {
	cmd.OrgId = c.OrgId
	cmd.Id = c.ParamsInt64(":teamId")
	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrTeamNameTaken {
			return Error(400, "Teamname bereits in benutzung", err)
		}
		return Error(500, "Aktualisieren des Teams ist fehlgeschlagen", err)
	}

	return Success("Team updated")
}

// DELETE /api/teams/:teamId
func DeleteTeamByID(c *m.ReqContext) Response {
	if err := bus.Dispatch(&m.DeleteTeamCommand{OrgId: c.OrgId, Id: c.ParamsInt64(":teamId")}); err != nil {
		if err == m.ErrTeamNotFound {
			return Error(404, "Löschen des Teams fehlgeschlagen. ID nicht gefunden", nil)
		}
		return Error(500, "Aktualisieren des Teams fehlgeschlagen", err)
	}
	return Success("Team gelöscht")
}

// GET /api/teams/search
func SearchTeams(c *m.ReqContext) Response {
	perPage := c.QueryInt("perpage")
	if perPage <= 0 {
		perPage = 1000
	}
	page := c.QueryInt("page")
	if page < 1 {
		page = 1
	}

	query := m.SearchTeamsQuery{
		OrgId: c.OrgId,
		Query: c.Query("query"),
		Name:  c.Query("name"),
		Page:  page,
		Limit: perPage,
	}

	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Fehler bei der Suche des Teams", err)
	}

	for _, team := range query.Result.Teams {
		team.AvatarUrl = dtos.GetGravatarUrlWithDefault(team.Email, team.Name)
	}

	query.Result.Page = page
	query.Result.PerPage = perPage

	return JSON(200, query.Result)
}

// GET /api/teams/:teamId
func GetTeamByID(c *m.ReqContext) Response {
	query := m.GetTeamByIdQuery{OrgId: c.OrgId, Id: c.ParamsInt64(":teamId")}

	if err := bus.Dispatch(&query); err != nil {
		if err == m.ErrTeamNotFound {
			return Error(404, "Team nicht gefunden", err)
		}

		return Error(500, "Fehler bei der Abfrage des Teams", err)
	}

	query.Result.AvatarUrl = dtos.GetGravatarUrlWithDefault(query.Result.Email, query.Result.Name)
	return JSON(200, &query.Result)
}
