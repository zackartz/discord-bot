package cmds

import (
	"github.com/zackartz/cmdlr2"
	"synergy/cmds/role"
)

var RoleCommand = &cmdlr2.Command{
	Name:        "role",
	Aliases:     []string{"r"},
	Description: "Add roles to the server's role selection menu and more.",
	Usage:       "Use one of the sub commands.",
	Example:     "role add @Role :emoji:",
	Flags:       []string{"role", "emoji"},
	IgnoreCase:  true,
	SubCommands: []*cmdlr2.Command{
		role.CreateMessageRoleCommand,
		role.AddRoleCommand,
		role.RemoveRoleCommand,
	},
	Handler: nil,
}

func init() {
	CommandList = append(CommandList, RoleCommand)
}
