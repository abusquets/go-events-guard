package cli

import (
	core_service_ports "eventsguard/internal/core/domain/ports/services"

	"github.com/urfave/cli/v2"
)

type CoreCommands struct {
	Commands []cli.Command
}

func NewCoreCommands(userService core_service_ports.UserService) CoreCommands {
	commands := []cli.Command{}
	commands = append(commands, userCommands(userService)...)
	commands = append(commands, secretCommands()...)
	return CoreCommands{
		Commands: commands,
	}
}
