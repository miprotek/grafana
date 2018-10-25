package api

import (
	"strings"

	"github.com/miprotek/grafana-de/pkg/api/dtos"
	"github.com/miprotek/grafana-de/pkg/components/simplejson"
	m "github.com/miprotek/grafana-de/pkg/models"
	"github.com/miprotek/grafana-de/pkg/services/annotations"
	"github.com/miprotek/grafana-de/pkg/services/guardian"
	"github.com/miprotek/grafana-de/pkg/util"
)

func GetAnnotations(c *m.ReqContext) Response {

	query := &annotations.ItemQuery{
		From:        c.QueryInt64("from"),
		To:          c.QueryInt64("to"),
		OrgId:       c.OrgId,
		UserId:      c.QueryInt64("userId"),
		AlertId:     c.QueryInt64("alertId"),
		DashboardId: c.QueryInt64("dashboardId"),
		PanelId:     c.QueryInt64("panelId"),
		Limit:       c.QueryInt64("limit"),
		Tags:        c.QueryStrings("tags"),
		Type:        c.Query("type"),
		MatchAny:    c.QueryBool("matchAny"),
	}

	repo := annotations.GetRepository()

	items, err := repo.Find(query)
	if err != nil {
		return Error(500, "Fehler beim Abrufen von Anmerkungen", err)
	}

	for _, item := range items {
		if item.Email != "" {
			item.AvatarUrl = dtos.GetGravatarUrl(item.Email)
		}
	}

	return JSON(200, items)
}

type CreateAnnotationError struct {
	message string
}

func (e *CreateAnnotationError) Error() string {
	return e.message
}

func PostAnnotation(c *m.ReqContext, cmd dtos.PostAnnotationsCmd) Response {
	if canSave, err := canSaveByDashboardID(c, cmd.DashboardId); err != nil || !canSave {
		return dashboardGuardianResponse(err)
	}

	repo := annotations.GetRepository()

	if cmd.Text == "" {
		err := &CreateAnnotationError{"Textfeld sollte nicht leer sein"}
		return Error(500, "Speichern der Anmerkung fehlgeschlagen", err)
	}

	item := annotations.Item{
		OrgId:       c.OrgId,
		UserId:      c.UserId,
		DashboardId: cmd.DashboardId,
		PanelId:     cmd.PanelId,
		Epoch:       cmd.Time,
		Text:        cmd.Text,
		Data:        cmd.Data,
		Tags:        cmd.Tags,
	}

	if err := repo.Save(&item); err != nil {
		return Error(500, "Speichern der Anmerkung fehlgeschlagen", err)
	}

	startID := item.Id

	// handle regions
	if cmd.IsRegion {
		item.RegionId = startID

		if item.Data == nil {
			item.Data = simplejson.New()
		}

		if err := repo.Update(&item); err != nil {
			return Error(500, "Fehler beim Festlegen der regionId für die Anmerkung", err)
		}

		item.Id = 0
		item.Epoch = cmd.TimeEnd

		if err := repo.Save(&item); err != nil {
			return Error(500, "Fehler beim speichern der Anmerkung für die Endzeit der Region", err)
		}

		return JSON(200, util.DynMap{
			"message": "Anmerkung hinzugefügt",
			"id":      startID,
			"endId":   item.Id,
		})
	}

	return JSON(200, util.DynMap{
		"message": "Anmerkung hinzugefügt",
		"id":      startID,
	})
}

func formatGraphiteAnnotation(what string, data string) string {
	text := what
	if data != "" {
		text = text + "\n" + data
	}
	return text
}

func PostGraphiteAnnotation(c *m.ReqContext, cmd dtos.PostGraphiteAnnotationsCmd) Response {
	repo := annotations.GetRepository()

	if cmd.What == "" {
		err := &CreateAnnotationError{"what field should not be empty"}
		return Error(500, "Fehler beim speichern der Graphite Anmerkung", err)
	}

	text := formatGraphiteAnnotation(cmd.What, cmd.Data)

	// Support tags in prior to Graphite 0.10.0 format (string of tags separated by space)
	var tagsArray []string
	switch tags := cmd.Tags.(type) {
	case string:
		if tags != "" {
			tagsArray = strings.Split(tags, " ")
		} else {
			tagsArray = []string{}
		}
	case []interface{}:
		for _, t := range tags {
			if tagStr, ok := t.(string); ok {
				tagsArray = append(tagsArray, tagStr)
			} else {
				err := &CreateAnnotationError{"tag sollte eine Zeichenfolge sein"}
				return Error(500, "Fehler beim speichern der Graphite Anmerkung", err)
			}
		}
	default:
		err := &CreateAnnotationError{"Nicht unterstütztes tag Format"}
		return Error(500, "Fehler beim speichern der Graphite Anmerkung", err)
	}

	item := annotations.Item{
		OrgId:  c.OrgId,
		UserId: c.UserId,
		Epoch:  cmd.When * 1000,
		Text:   text,
		Tags:   tagsArray,
	}

	if err := repo.Save(&item); err != nil {
		return Error(500, "Fehler beim speichern der Graphite Anmerkung", err)
	}

	return JSON(200, util.DynMap{
		"message": "Graphite Meldung hinzugefügt",
		"id":      item.Id,
	})
}

func UpdateAnnotation(c *m.ReqContext, cmd dtos.UpdateAnnotationsCmd) Response {
	annotationID := c.ParamsInt64(":annotationId")

	repo := annotations.GetRepository()

	if resp := canSave(c, repo, annotationID); resp != nil {
		return resp
	}

	item := annotations.Item{
		OrgId:  c.OrgId,
		UserId: c.UserId,
		Id:     annotationID,
		Epoch:  cmd.Time,
		Text:   cmd.Text,
		Tags:   cmd.Tags,
	}

	if err := repo.Update(&item); err != nil {
		return Error(500, "Aktualisieren der Meldung fehlgeschlagen", err)
	}

	if cmd.IsRegion {
		itemRight := item
		itemRight.RegionId = item.Id
		itemRight.Epoch = cmd.TimeEnd

		// We don't know id of region right event, so set it to 0 and find then using query like
		// ... WHERE region_id = <item.RegionId> AND id != <item.RegionId> ...
		itemRight.Id = 0

		if err := repo.Update(&itemRight); err != nil {
			return Error(500, "Fehler beim aktualisieren der Meldung für die Regionsendzeit", err)
		}
	}

	return Success("Meldung Aktualisiert")
}

func DeleteAnnotations(c *m.ReqContext, cmd dtos.DeleteAnnotationsCmd) Response {
	repo := annotations.GetRepository()

	err := repo.Delete(&annotations.DeleteParams{
		OrgId:       c.OrgId,
		Id:          cmd.AnnotationId,
		RegionId:    cmd.RegionId,
		DashboardId: cmd.DashboardId,
		PanelId:     cmd.PanelId,
	})

	if err != nil {
		return Error(500, "Löschen der Meldung fehlgeschlagen", err)
	}

	return Success("Meldungen gelöscht")
}

func DeleteAnnotationByID(c *m.ReqContext) Response {
	repo := annotations.GetRepository()
	annotationID := c.ParamsInt64(":annotationId")

	if resp := canSave(c, repo, annotationID); resp != nil {
		return resp
	}

	err := repo.Delete(&annotations.DeleteParams{
		OrgId: c.OrgId,
		Id:    annotationID,
	})

	if err != nil {
		return Error(500, "Löschen der Meldung fehlgeschlagen", err)
	}

	return Success("Meldung gelöscht")
}

func DeleteAnnotationRegion(c *m.ReqContext) Response {
	repo := annotations.GetRepository()
	regionID := c.ParamsInt64(":regionId")

	if resp := canSave(c, repo, regionID); resp != nil {
		return resp
	}

	err := repo.Delete(&annotations.DeleteParams{
		OrgId:    c.OrgId,
		RegionId: regionID,
	})

	if err != nil {
		return Error(500, "Der Meldungsbereich konnte nicht gelöscht werden", err)
	}

	return Success("Meldungsbereich gelöscht")
}

func canSaveByDashboardID(c *m.ReqContext, dashboardID int64) (bool, error) {
	if dashboardID == 0 && !c.SignedInUser.HasRole(m.ROLE_EDITOR) {
		return false, nil
	}

	if dashboardID != 0 {
		guard := guardian.New(dashboardID, c.OrgId, c.SignedInUser)
		if canEdit, err := guard.CanEdit(); err != nil || !canEdit {
			return false, err
		}
	}

	return true, nil
}

func canSave(c *m.ReqContext, repo annotations.Repository, annotationID int64) Response {
	items, err := repo.Find(&annotations.ItemQuery{AnnotationId: annotationID, OrgId: c.OrgId})

	if err != nil || len(items) == 0 {
		return Error(500, "Die zu aktualisierende Anmerkung konnte nicht gefunden werden", err)
	}

	dashboardID := items[0].DashboardId

	if canSave, err := canSaveByDashboardID(c, dashboardID); err != nil || !canSave {
		return dashboardGuardianResponse(err)
	}

	return nil
}
