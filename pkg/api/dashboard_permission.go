package api

import (
	"time"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/guardian"
)

func GetDashboardPermissionList(c *m.ReqContext) Response {
	dashID := c.ParamsInt64(":dashboardId")

	_, rsp := getDashboardHelper(c.OrgId, "", dashID, "")
	if rsp != nil {
		return rsp
	}

	g := guardian.New(dashID, c.OrgId, c.SignedInUser)

	if canAdmin, err := g.CanAdmin(); err != nil || !canAdmin {
		return dashboardGuardianResponse(err)
	}

	acl, err := g.GetAcl()
	if err != nil {
		return Error(500, "Fehler beim abrufen der Dashboard-Berechtigungen", err)
	}

	for _, perm := range acl {
		if perm.Slug != "" {
			perm.Url = m.GetDashboardFolderUrl(perm.IsFolder, perm.Uid, perm.Slug)
		}
	}

	return JSON(200, acl)
}

func UpdateDashboardPermissions(c *m.ReqContext, apiCmd dtos.UpdateDashboardAclCommand) Response {
	dashID := c.ParamsInt64(":dashboardId")

	_, rsp := getDashboardHelper(c.OrgId, "", dashID, "")
	if rsp != nil {
		return rsp
	}

	g := guardian.New(dashID, c.OrgId, c.SignedInUser)
	if canAdmin, err := g.CanAdmin(); err != nil || !canAdmin {
		return dashboardGuardianResponse(err)
	}

	cmd := m.UpdateDashboardAclCommand{}
	cmd.DashboardId = dashID

	for _, item := range apiCmd.Items {
		cmd.Items = append(cmd.Items, &m.DashboardAcl{
			OrgId:       c.OrgId,
			DashboardId: dashID,
			UserId:      item.UserId,
			TeamId:      item.TeamId,
			Role:        item.Role,
			Permission:  item.Permission,
			Created:     time.Now(),
			Updated:     time.Now(),
		})
	}

	if okToUpdate, err := g.CheckPermissionBeforeUpdate(m.PERMISSION_ADMIN, cmd.Items); err != nil || !okToUpdate {
		if err != nil {
			if err == guardian.ErrGuardianPermissionExists ||
				err == guardian.ErrGuardianOverride {
				return Error(400, err.Error(), err)
			}

			return Error(500, "Fehler beim überprüfen der Dashboard-Berechtigungen", err)
		}

		return Error(403, "Die eigene Administratorberechtigung für einen Ordner kann nicht entfernt werden", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrDashboardAclInfoMissing || err == m.ErrDashboardPermissionDashboardEmpty {
			return Error(409, err.Error(), err)
		}
		return Error(500, "Fehler beim erstellen der Berechtigung", err)
	}

	return Success("Dashboard-Berechtigungen aktualisiert")
}
