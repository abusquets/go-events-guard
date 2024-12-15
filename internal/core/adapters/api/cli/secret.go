package cli

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

func secretCommands() []cli.Command {
	commands := []cli.Command{}

	createSecret := &cli.Command{
		Name:  "create-secret",
		Usage: "Creates a new secret",
		Action: func(cCtx *cli.Context) error {

			secret := make([]byte, 32)
			_, err := rand.Read(secret)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Generated Secret:")
			fmt.Println(hex.EncodeToString(secret))

			return nil
		},
	}

	commands = append(commands, *createSecret)

	return commands
}

func NewSecretCommands() []cli.Command {
	return secretCommands()
}
