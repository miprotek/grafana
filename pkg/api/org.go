package api

import (
	"github.com/miprotek/grafana-de/pkg/api/dtos"
	"github.com/miprotek/grafana-de/pkg/bus"
	"github.com/miprotek/grafana-de/pkg/metrics"
	m "github.com/miprotek/grafana-de/pkg/models"
	"github.com/miprotek/grafana-de/pkg/setting"
	"github.com/miprotek/grafana-de/pkg/util"
)

// GET /api/org
func GetOrgCurrent(c *m.ReqContext) Response {
	return getOrgHelper(c.OrgId)
}

// GET /api/orgs/:orgId
func GetOrgByID(c *m.ReqContext) Response {
	return getOrgHelper(c.ParamsInt64(":orgId"))
}

// Get /api/orgs/name/:name
func GetOrgByName(c *m.ReqContext) Response {
	query := m.GetOrgByNameQuery{Name: c.Params(":name")}
	if err := bus.Dispatch(&query); err != nil {
		if err == m.ErrOrgNotFound {
			return Error(404, "Organisation nicht gefunden", err)
		}

		return Error(500, "Fehler beim erhalten der Organisation", err)
	}
	org := query.Result
	result := m.OrgDetailsDTO{
		Id:   org.Id,
		Name: org.Name,
		Address: m.Address{
			Address1: org.Address1,
			Address2: org.Address2,
			City:     org.City,
			ZipCode:  org.ZipCode,
			State:    org.State,
			Country:  org.Country,
		},
	}

	return JSON(200, &result)
}

func getOrgHelper(orgID int64) Response {
	query := m.GetOrgByIdQuery{Id: orgID}

	if err := bus.Dispatch(&query); err != nil {
		if err == m.ErrOrgNotFound {
			return Error(404, "Organisation nicht gefunden", err)
		}

		return Error(500, "Fehler beim erhalten der Organisation", err)
	}

	org := query.Result
	result := m.OrgDetailsDTO{
		Id:   org.Id,
		Name: org.Name,
		Address: m.Address{
			Address1: org.Address1,
			Address2: org.Address2,
			City:     org.City,
			ZipCode:  org.ZipCode,
			State:    org.State,
			Country:  org.Country,
		},
	}

	return JSON(200, &result)
}

// POST /api/orgs
func CreateOrg(c *m.ReqContext, cmd m.CreateOrgCommand) Response {
	if !c.IsSignedIn || (!setting.AllowUserOrgCreate && !c.IsGrafanaAdmin) {
		return Error(403, "Zugriff verweigert", nil)
	}

	cmd.UserId = c.UserId
	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrOrgNameTaken {
			return Error(409, "Organisationsname bereits in benutzung", err)
		}
		return Error(500, "Erstellen der Organisation fehlgeschlagen", err)
	}

	metrics.M_Api_Org_Create.Inc()

	return JSON(200, &util.DynMap{
		"orgId":   cmd.Result.Id,
		"message": "Organisation erstellt",
	})
}

// PUT /api/org
func UpdateOrgCurrent(c *m.ReqContext, form dtos.UpdateOrgForm) Response {
	return updateOrgHelper(form, c.OrgId)
}

// PUT /api/orgs/:orgId
func UpdateOrg(c *m.ReqContext, form dtos.UpdateOrgForm) Response {
	return updateOrgHelper(form, c.ParamsInt64(":orgId"))
}

func updateOrgHelper(form dtos.UpdateOrgForm, orgID int64) Response {
	cmd := m.UpdateOrgCommand{Name: form.Name, OrgId: orgID}
	if err := bus.Dispatch(&cmd); err != nil {
		if err == m.ErrOrgNameTaken {
			return Error(400, "Organisationsname bereits in benutzung", err)
		}
		return Error(500, "Aktualisierung der Organisation fehlgeschlagen", err)
	}

	return Success("Organisation aktualisiert")
}

// PUT /api/org/address
func UpdateOrgAddressCurrent(c *m.ReqContext, form dtos.UpdateOrgAddressForm) Response {
	return updateOrgAddressHelper(form, c.OrgId)
}

// PUT /api/orgs/:orgId/address
func UpdateOrgAddress(c *m.ReqContext, form dtos.UpdateOrgAddressForm) Response {
	return updateOrgAddressHelper(form, c.ParamsInt64(":orgId"))
}

func updateOrgAddressHelper(form dtos.UpdateOrgAddressForm, orgID int64) Response {
	cmd := m.UpdateOrgAddressCommand{
		OrgId: orgID,
		Address: m.Address{
			Address1: form.Address1,
			Address2: form.Address2,
			City:     form.City,
			State:    form.State,
			ZipCode:  form.ZipCode,
			Country:  form.Country,
		},
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Fehler beim aktualisieren der Organisationsadresse", err)
	}

	return Success("Addresse aktualisiert")
}

// GET /api/orgs/:orgId
func DeleteOrgByID(c *m.ReqContext) Response {
	if err := bus.Dispatch(&m.DeleteOrgCommand{Id: c.ParamsInt64(":orgId")}); err != nil {
		if err == m.ErrOrgNotFound {
			return Error(404, "Löschen der Organisation fehlgeschlagen. ID nicht gefunden", nil)
		}
		return Error(500, "Aktualisieren der Organisation fehlgeschlagen", err)
	}
	return Success("Organisation gelöscht")
}

func SearchOrgs(c *m.ReqContext) Response {
	query := m.SearchOrgsQuery{
		Query: c.Query("query"),
		Name:  c.Query("name"),
		Page:  0,
		Limit: 1000,
	}

	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Suche nach Organisationen fehlgeschlagen", err)
	}

	return JSON(200, query.Result)
}
