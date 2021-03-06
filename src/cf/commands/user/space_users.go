package user

import (
	"cf/api"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/codegangsta/cli"
)

type SpaceUsers struct {
	ui        terminal.UI
	spaceRepo api.SpaceRepository
	userRepo  api.UserRepository
	orgReq    requirements.OrganizationRequirement
}

func NewSpaceUsers(ui terminal.UI, spaceRepo api.SpaceRepository, userRepo api.UserRepository) (cmd *SpaceUsers) {
	cmd = new(SpaceUsers)
	cmd.ui = ui
	cmd.spaceRepo = spaceRepo
	cmd.userRepo = userRepo
	return
}

func (cmd *SpaceUsers) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "space-users")
		return
	}

	orgName := c.Args()[0]
	cmd.orgReq = reqFactory.NewOrganizationRequirement(orgName)
	reqs = append(reqs, reqFactory.NewLoginRequirement(), cmd.orgReq)

	return
}

func (cmd *SpaceUsers) Run(c *cli.Context) {
	spaceName := c.Args()[1]
	org := cmd.orgReq.GetOrganization()

	space, apiResponse := cmd.spaceRepo.FindByNameInOrg(spaceName, org)
	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed(apiResponse.Message)
	}

	cmd.ui.Say("Getting users in space %s in org %s",
		terminal.EntityNameColor(space.Name),
		terminal.EntityNameColor(org.Name))

	cmd.userRepo.FindAllInSpaceByRole(space)
	usersByRole, apiResponse := cmd.userRepo.FindAllInSpaceByRole(space)
	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed(apiResponse.Message)
	}

	cmd.ui.Ok()

	for role, users := range usersByRole {
		cmd.ui.Say("")
		cmd.ui.Say("%s", terminal.HeaderColor(role))

		for _, user := range users {
			cmd.ui.Say("  %s", user.Username)
		}
	}
}
