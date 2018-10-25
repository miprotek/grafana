package provisioning

import (
	"context"
	"fmt"
	"path"

	"github.com/miprotek/grafana-de/pkg/registry"
	"github.com/miprotek/grafana-de/pkg/services/provisioning/dashboards"
	"github.com/miprotek/grafana-de/pkg/services/provisioning/datasources"
	"github.com/miprotek/grafana-de/pkg/setting"
)

func init() {
	registry.RegisterService(&ProvisioningService{})
}

type ProvisioningService struct {
	Cfg *setting.Cfg `inject:""`
}

func (ps *ProvisioningService) Init() error {
	datasourcePath := path.Join(ps.Cfg.ProvisioningPath, "datasources")
	if err := datasources.Provision(datasourcePath); err != nil {
		return fmt.Errorf("Datasource provisioning error: %v", err)
	}

	return nil
}

func (ps *ProvisioningService) Run(ctx context.Context) error {
	dashboardPath := path.Join(ps.Cfg.ProvisioningPath, "dashboards")
	dashProvisioner := dashboards.NewDashboardProvisioner(dashboardPath)

	if err := dashProvisioner.Provision(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}
