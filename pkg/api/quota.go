package api

import (
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

func GetOrgQuotas(c *m.ReqContext) Response {
	if !setting.Quota.Enabled {
		return Error(404, "Kontingente sind nicht aktiviert", nil)
	}
	query := m.GetOrgQuotasQuery{OrgId: c.ParamsInt64(":orgId")}

	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Kontingente konnten nicht abgerufen werden", err)
	}

	return JSON(200, query.Result)
}

func UpdateOrgQuota(c *m.ReqContext, cmd m.UpdateOrgQuotaCmd) Response {
	if !setting.Quota.Enabled {
		return Error(404, "Kontingente sind nicht aktiviert", nil)
	}
	cmd.OrgId = c.ParamsInt64(":orgId")
	cmd.Target = c.Params(":target")

	if _, ok := setting.Quota.Org.ToMap()[cmd.Target]; !ok {
		return Error(404, "Ungültiges Kontingentziel", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Organisationskontingente konnten nicht aktualisert werden", err)
	}
	return Success("Organisations Kontingent aktualisiert")
}

func GetUserQuotas(c *m.ReqContext) Response {
	if !setting.Quota.Enabled {
		return Error(404, "Kontingente sind nicht aktiviert", nil)
	}
	query := m.GetUserQuotasQuery{UserId: c.ParamsInt64(":id")}

	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Kontingente konnten nicht abgerufen werden", err)
	}

	return JSON(200, query.Result)
}

func UpdateUserQuota(c *m.ReqContext, cmd m.UpdateUserQuotaCmd) Response {
	if !setting.Quota.Enabled {
		return Error(404, "Kontingente sind nicht aktiviert", nil)
	}
	cmd.UserId = c.ParamsInt64(":id")
	cmd.Target = c.Params(":target")

	if _, ok := setting.Quota.User.ToMap()[cmd.Target]; !ok {
		return Error(404, "Ungültiges kontingent Ziel", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Organisationskontingente konnten nicht aktualisert werden", err)
	}
	return Success("Organisationskontingent aktualisiert")
}
