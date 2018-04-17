package api

import (
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

func StarDashboard(c *m.ReqContext) Response {
	if !c.IsSignedIn {
		return Error(412, "Sie müssen sich anmelden um Dashboards zu markieren", nil)
	}

	cmd := m.StarDashboardCommand{UserId: c.UserId, DashboardId: c.ParamsInt64(":id")}

	if cmd.DashboardId <= 0 {
		return Error(400, "Dashboard id fehlt", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Fehler beim anzeigen des Dashboards", err)
	}

	return Success("Dashboard markiert!")
}

func UnstarDashboard(c *m.ReqContext) Response {

	cmd := m.UnstarDashboardCommand{UserId: c.UserId, DashboardId: c.ParamsInt64(":id")}

	if cmd.DashboardId <= 0 {
		return Error(400, "Dashboard id fehlt", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "Kann Dashboard nicht unmarkieren", err)
	}

	return Success("Dashboard unmarkiert")
}
