package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

func initFlags(set *pflag.FlagSet) {
	set.StringVarP(
		&params.URL,
		"url", "u",
		"",
		"[driver://username:password@address/dbname?param=value] select a resource using the URL format",
	)
	set.StringVar(
		&params.DevURL,
		"dev-url",
		"",
		"[driver://username:password@address/dbname?param=value] select a dev database using the URL format",
	)
	set.StringSliceVarP(
		&params.Schema,
		"schema", "s",
		nil,
		"set schema names",
	)
	set.StringSliceVar(
		&params.Exclude,
		"exclude",
		nil,
		"list of glob patterns used to filter resources from applying",
	)

	set.StringVar(&params.Env, "env", "", "set which env from the config file to use")
	set.StringVarP(&params.ConfigURL, "config", "c", "", "select config (project) file using URL format")

	set.StringVar(&params.AtlasCliPath, "atlas", "", "path of the atlas cli, defaults to look in PATH")
}

func verifyFileExists(fpath string) error {
	fi, err := os.Stat(fpath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("atlas cli not found at '%s'", fpath)
	}
	if err != nil {
		return fmt.Errorf("failed to access atlas cli at '%s': %w", fpath, err)
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("atlas cli at '%s' is not a regular file", fpath)
	}
	return nil
}
