package main

import (
	"fmt"
	"os"

	"eventsguard/internal/app"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/server"

	"eventsguard/internal/infrastructure/migrations"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	core_commands "eventsguard/internal/core/adapters/api/cli"
)

var exitAfterCommand = true

func NewAllCommands(
	userCommands []cli.Command,
	migrationCommands []cli.Command,
) []cli.Command {
	// Combina les comandes d'usuari i migració en una sola llista
	return append(userCommands, migrationCommands...)
}

func CLIApp(
	lifecycle fx.Lifecycle,
	cfg *config.AppConfig,
	startServer func(),
	userCommands core_commands.CoreCommands,
	migrationCommands migrations.MigrationsCommands,
) *cli.App {
	cliApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "test",
				Usage: "Test command",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("CLI is working fine")
					return nil
				},
			},
			{
				Name:  "runserver",
				Usage: "Run the server",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("Starting server...")
					exitAfterCommand = false
					startServer()
					return nil
				},
			},
		},
	}

	allCommands := append(userCommands.Commands, migrationCommands.Commands...)

	commandPtrs := make([]*cli.Command, len(allCommands))
	for i := range allCommands {
		commandPtrs[i] = &allCommands[i]
	}

	cliApp.Commands = append(cliApp.Commands, commandPtrs...)

	return cliApp
}

func main() {
	appModule := fx.Options(app.Module)
	fxApp := fx.New(
		appModule,
		fx.Provide(
			CLIApp,
			func(
				lifecycle fx.Lifecycle,
				cfg *config.AppConfig,
				s server.WebServer,
			) func() {
				s.RegisterHooks(lifecycle)
				return s.Start
			},
		),
		fx.Invoke(func(cliApp *cli.App) {
			if err := cliApp.Run(os.Args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if exitAfterCommand {
				os.Exit(0)
			}
		}),
	)

	fxApp.Run()
}
