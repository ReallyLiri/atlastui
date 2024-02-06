package cmd

import (
	"context"
	"fmt"
	"github.com/reallyliri/atlaspect/inspect"
	"github.com/reallyliri/atlaspect/tui"
	"github.com/samber/lo"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	version  = "0.0.1"
	toolName = "atlaspect"
	atlasCli = "atlas"
)

var examples = []string{
	`-u "mysql://user:pass@localhost:3306/dbname"`,
	`-u "mariadb://user:pass@localhost:3306/" --schema=schemaA,schemaB -s schemaC`,
	`--url "postgres://user:pass@host:port/dbname?sslmode=disable"`,
	`-u "sqlite://file:ex1.db?_fk=1"`,
}

var params inspect.Params

var rootCmd = &cobra.Command{
	Use:     toolName,
	Version: version,
	Short:   "Beautiful terminal UI for inspecting your database schemas",
	Long: fmt.Sprintf(`Beautiful terminal UI for inspecting your database schemas.
This is a complimentary CLI tool for atlas, an "atlas schema inspect" on steroids if you will.
%s connects to the given database and visualize its schema.`, toolName),
	Example:      strings.Join(lo.Map(examples, func(example string, _ int) string { return fmt.Sprintf("  %s %s", toolName, example) }), "\n"),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetchData(cmd.Context())
		if err != nil {
			return err
		}
		return tui.Run(cmd.Context(), *data)
	},
}

func Execute() {
	rootCmd.HelpFunc()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initFlags(rootCmd.Flags())
}

func fetchData(ctx context.Context) (*inspect.Data, error) {
	if params.AtlasCliPath != "" {
		err := verifyFileExists(params.AtlasCliPath)
		if err != nil {
			return nil, err
		}
	} else {
		cliPath, err := exec.LookPath(atlasCli)
		if err != nil {
			return nil, fmt.Errorf("atlas cli not found in PATH")
		}
		params.AtlasCliPath = cliPath
	}

	if params.FromFilePath != "" {
		err := verifyFileExists(params.FromFilePath)
		if err != nil {
			return nil, err
		}
		data, err := inspect.LoadFromFile(params.FromFilePath)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	data, err := inspect.Inspect(ctx, &params)
	if err != nil {
		return nil, err
	}
	return data, nil
}
