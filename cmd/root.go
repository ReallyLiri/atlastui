package cmd

import (
	"fmt"
	"github.com/reallyliri/atlaspect/inspect"
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

var params inspect.InspectParams

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
		if params.AtlasCliPath != "" {
			err := verifyFileExists(params.AtlasCliPath)
			if err != nil {
				return err
			}
		} else {
			cliPath, err := exec.LookPath(atlasCli)
			if err != nil {
				return fmt.Errorf("atlas cli not found in PATH")
			}
			params.AtlasCliPath = cliPath
		}
		schemas, err := inspect.Inspect(cmd.Context(), &params)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", schemas)
		return nil
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
