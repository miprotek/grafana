package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/components/apikeygen"
	m "github.com/grafana/grafana/pkg/models"
)

func GetAPIKeys(c *m.ReqContext) Response {
	query := m.GetApiKeysQuery{OrgId: c.OrgId}

	if err := bus.Dispatch(&query); err != nil {
		return Error(500, "Konnte keine API-Schlüssel auflisten", err)
	}

	result := make([]*m.ApiKeyDTO, len(query.Result))
	for i, t := range query.Result {
		result[i] = &m.ApiKeyDTO{
			Id:   t.Id,
			Name: t.Name,
			Role: t.Role,
		}
	}

	return JSON(200, result)
}

func DeleteAPIKey(c *m.ReqContext) Response {
	id := c.ParamsInt64(":id")

	cmd := &m.DeleteApiKeyCommand{Id: id, OrgId: c.OrgId}

	err := bus.Dispatch(cmd)
	if err != nil {
		return Error(500, "Löschen der API-Schlüssel fehlgeschlagen", err)
	}

	return Success("API-Schlüssel gelöscht")
}

func AddAPIKey(c *m.ReqContext, cmd m.AddApiKeyCommand) Response {
	if !cmd.Role.IsValid() {
		return Error(400, "Falsche Rolle spezifiziert", nil)
	}

	cmd.OrgId = c.OrgId

	newKeyInfo := apikeygen.New(cmd.OrgId, cmd.Name)
	cmd.Key = newKeyInfo.HashedKey

	if err := bus.Dispatch(&cmd); err != nil {
		return Error(500, "API-Schlüssel hinzufügen fehlgeschlagen", err)
	}

	result := &dtos.NewApiKeyResult{
		Name: cmd.Result.Name,
		Key:  newKeyInfo.ClientSecret}

	return JSON(200, result)
}
