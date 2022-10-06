package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tmlbl/rem/config"
	"github.com/tmlbl/rem/provision"
	"github.com/tmlbl/rem/storage"
)

const appName = "rem"

// rem.yaml in current directory
func defaultConfigPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "rem.yaml")
}

func main() {
	root := cobra.Command{
		Use: fmt.Sprintf("%s [cmd] [args...]", appName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(defaultConfigPath())
			if err != nil {
				return err
			}

			provisioner, err := provision.GetProvisioner(cfg.Platform)
			if err != nil {
				return err
			}

			state, err := provisioner.Build(&cfg.Base)
			if err != nil {
				return err
			}

			db := storage.Default()

			fmt.Println(state)

			return db.SaveState(defaultConfigPath(), state)
		},
	}

	err := root.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
