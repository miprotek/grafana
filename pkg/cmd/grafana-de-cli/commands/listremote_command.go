package commands

import (
	"github.com/miprotek/grafana-de/pkg/cmd/grafana-de-cli/logger"
	s "github.com/miprotek/grafana-de/pkg/cmd/grafana-de-cli/services"
)

func listremoteCommand(c CommandLine) error {
	plugin, err := s.ListAllPlugins(c.RepoDirectory())

	if err != nil {
		return err
	}

	for _, i := range plugin.Plugins {
		pluginVersion := ""
		if len(i.Versions) > 0 {
			pluginVersion = i.Versions[0].Version
		}

		logger.Infof("id: %v version: %s\n", i.Id, pluginVersion)
	}

	return nil
}
