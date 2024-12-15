package cli

import (
	"fmt"
	"os"

	core_service_ports "eventsguard/internal/core/domain/ports/services"
	"eventsguard/internal/core/dtos"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func userCommands(userService core_service_ports.UserService) []cli.Command {
	commands := []cli.Command{}

	addUserFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "email",
			Value:    "",
			Usage:    "Email (username)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "first_name",
			Value:    "",
			Usage:    "First name of the user",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "last_name",
			Value: "",
			Usage: "Last name of the user",
		},
	}

	addUser := &cli.Command{
		Name:  "adduser",
		Usage: "Add a new admin user",
		Flags: addUserFlags,
		Action: func(cCtx *cli.Context) error {
			email := cCtx.String("email")
			firstName := cCtx.String("first_name")
			lastName := cCtx.String("last_name")

			fmt.Print("Password: ")
			password, err := readPassword()
			if err != nil {
				return err
			}

			userData := dtos.CreateUserInput{
				Email:     email,
				FirstName: firstName,
				Password:  password,
				IsActive:  true,
			}
			if lastName != "" {
				userData.LastName = &lastName
			}

			user, appErr := userService.CreateUser(cCtx.Context, userData)
			if appErr != nil {
				return appErr
			} else {
				fmt.Printf("User %s created\n", user.ID)
			}
			id := user.ID.String()

			isAdmin := true
			inData := dtos.UpdatePartialAdminUserInput{
				IsAdmin: &isAdmin,
			}
			_, appErr = userService.UpdatePartialUser(cCtx.Context, id, inData)

			if appErr != nil {
				return appErr
			}

			return nil
		},
	}

	setAdminFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "email",
			Value:    "",
			Usage:    "Email (username)",
			Required: true,
		},
	}

	setAdmin := &cli.Command{
		Name:  "admin",
		Usage: "Add Admin persmission",
		Flags: setAdminFlags,
		Action: func(cCtx *cli.Context) error {
			email := cCtx.String("email")

			user, appErr := userService.GetUserByEmail(cCtx.Context, email)
			if appErr != nil {
				return appErr
			}
			id := user.ID.String()

			isAdmin := true
			inData := dtos.UpdatePartialAdminUserInput{
				IsAdmin: &isAdmin,
			}
			_, appErr = userService.UpdatePartialUser(cCtx.Context, id, inData)

			if appErr != nil {
				return appErr
			}
			fmt.Println("Admin permission added to user: ", email)
			return nil
		},
	}

	commands = append(commands, *addUser, *setAdmin)

	return commands
}

func readPassword() (string, error) {
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("error llegint la contrasenya: %v", err)
	}
	fmt.Println()
	return string(bytePassword), nil
}

func NewUserCommands(userService core_service_ports.UserService) []cli.Command {
	return userCommands(userService)
}
