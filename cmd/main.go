package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/vaguevoid/cloud-cli/internal/api"
	"github.com/vaguevoid/cloud-cli/internal/domain/account"
	"github.com/vaguevoid/cloud-cli/internal/domain/share"
	"github.com/vaguevoid/cloud-cli/internal/lib/httpx"
	"github.com/vaguevoid/cloud-cli/internal/lib/pp"
	"github.com/vaguevoid/cloud-cli/internal/lib/system"
)

//-------------------------------------------------------------------------------------------------

const (
	CommandName              = "void-cloud"
	CommandDescription       = "access to the Void Cloud Platform"
	CommandVersion           = "0.0.1"
	ProductionURL            = "https://play.void.dev/"
	LoginCommandName         = "login"
	LoginCommandDescription  = "tell us who you are"
	DeployCommandName        = "deploy"
	DeployCommandDescription = "share your game with others"
)

//-------------------------------------------------------------------------------------------------

func main() {

	cmd := &cli.Command{
		Name:    CommandName,
		Usage:   CommandDescription,
		Version: CommandVersion,
		Commands: []*cli.Command{
			loginCommand(),
			deployCommand(),
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

//-------------------------------------------------------------------------------------------------

func serverFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "server",
		Usage:   "server endpoint `URL`",
		Sources: cli.EnvVars("SERVER"),
		Value:   ProductionURL,
	}
}

func orgFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "org",
		Usage:   "organization ID",
		Sources: cli.EnvVars("ORG"),
	}
}

func gameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "game",
		Usage:   "game ID",
		Sources: cli.EnvVars("GAME"),
	}
}

func tokenFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "token",
		Usage:   "personal access TOKEN",
		Sources: cli.EnvVars("TOKEN"),
	}
}

//-------------------------------------------------------------------------------------------------

func loginCommand() *cli.Command {

	return &cli.Command{
		Name:               LoginCommandName,
		Usage:              LoginCommandDescription,
		Flags:              []cli.Flag{serverFlag()},
		CustomHelpTemplate: SubcommandHelpTemplate,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			server := cmd.String("server")
			fmt.Println("logging in to", server, "...")
			user, err := account.Login(&account.LoginCommand{
				Server:  server,
				Runtime: system.DefaultRuntime(),
				Keyring: system.DefaultKeyring(server),
			})
			if err != nil {
				return err
			}
			fmt.Println("You are logged in")
			fmt.Println(pp.JSON(user))
			return nil
		},
	}
}

//-------------------------------------------------------------------------------------------------

func deployCommand() *cli.Command {

	return &cli.Command{
		Name:      DeployCommandName,
		Usage:     DeployCommandDescription,
		ArgsUsage: "PATH [LABEL]",
		Flags: []cli.Flag{
			serverFlag(),
			orgFlag(),
			gameFlag(),
			tokenFlag(),
		},
		CustomHelpTemplate: SubcommandHelpTemplate,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			org := cmd.String("org")
			game := cmd.String("game")
			path := cmd.Args().Get(0)
			label := cmd.Args().Get(1)

			if path == "" {
				return fmt.Errorf("missing required argument: PATH")
			}

			api, err := buildAPIClient(cmd)
			if err != nil {
				return err
			}

			fmt.Printf("Deploying %s ...\n", path)
			result, err := share.Deploy(&share.DeployCommand{
				API:   api,
				Org:   org,
				Game:  game,
				Label: label,
				Path:  path,
				OnStarted: func(deployID int64, manifest []share.DeployEntry, incremental []share.DeployEntry) {
					total := len(manifest)
					count := len(incremental)
					if total == count {
						fmt.Printf("deploying ALL %d files\n", total)
					} else {
						fmt.Printf("deploying %d / %d files\n", count, total)
					}
				},
				OnUpload: func(deployID int64, path string) {
					fmt.Printf("deploying %s\n", path)
				},
			})
			if err != nil {
				return err
			}

			fmt.Printf("Deployed to %s\n", result.URL)
			return nil
		},
	}
}

// -------------------------------------------------------------------------------------------------

func buildAPIClient(cmd *cli.Command) (*api.Client, error) {
	server := cmd.String("server")
	token := cmd.String("token")
	if token == "" {
		jwt, _ := system.DefaultKeyring(server).Get(httpx.ParamJWT)
		token = jwt
	}
	return api.NewClient(server, token)
}

//-------------------------------------------------------------------------------------------------

var SubcommandHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.FullName}}{{if .VisibleCommands}} [command [command options]]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Arguments}} [arguments...]{{end}}{{end}}{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`

//-------------------------------------------------------------------------------------------------
