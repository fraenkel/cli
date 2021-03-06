package application

import (
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"github.com/codegangsta/cli"
	"strings"
)

type ListApps struct {
	ui             terminal.UI
	config         *configuration.Configuration
	appSummaryRepo api.AppSummaryRepository
}

func NewListApps(ui terminal.UI, config *configuration.Configuration, appSummaryRepo api.AppSummaryRepository) (cmd ListApps) {
	cmd.ui = ui
	cmd.config = config
	cmd.appSummaryRepo = appSummaryRepo
	return
}

func (cmd ListApps) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		reqFactory.NewTargetedSpaceRequirement(),
	}
	return
}

func (cmd ListApps) Run(c *cli.Context) {
	cmd.ui.Say("Getting apps in %s...",
		terminal.EntityNameColor(cmd.config.Space.Name))

	apps, apiResponse := cmd.appSummaryRepo.GetSummariesInCurrentSpace()

	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed(apiResponse.Message)
		return
	}

	cmd.ui.Ok()

	table := [][]string{
		[]string{"name", "state", "instances", "memory", "disk", "urls"},
	}

	for _, app := range apps {
		table = append(table, []string{
			app.Name,
			coloredAppState(app),
			coloredAppInstaces(app),
			byteSize(app.Memory * MEGABYTE),
			byteSize(app.DiskQuota * MEGABYTE),
			strings.Join(app.Urls, ", "),
		})
	}

	cmd.ui.DisplayTable(table)
}
