/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wizcas/ame/api"
)

// syncSecretCmd represents the sync command
var syncSecretCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSyncSecret(args)
	},
}

func init() {
	secretCmd.AddCommand(syncSecretCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type DbSecretsConfig struct {
	Outfile string
	List    []DbSecretItem
}
type DbSecretItem struct {
	Name     string
	Env      string
	Arn      string
	Username string
	Password string
}

func (item DbSecretItem) GetVarName() string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", item.Name, item.Env))
}

func runSyncSecret(args []string) {
	m := viper.Get("DbSecrets")
	data, err := json.Marshal(m)
	if err != nil {
		log.WithField("reason", "DB secrets config").Fatal(err)
	}
	var config DbSecretsConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.WithField("reason", "DB secrets config").Fatal(err)
	}

	f, err := os.Create(path.Join(localFolder, config.Outfile))
	if err != nil {
		log.WithField("reason", "creating local DB secrets").Fatal(err)
	}
	defer f.Close()

	for _, item := range config.List {
		varName := item.GetVarName()

		log.WithField("for", varName).Info("fetching secrets...")

		if len(item.Username) == 0 {
			item.Username = "username"
		}
		if len(item.Password) == 0 {
			item.Password = "password"
		}

		username := api.GetSecretByIDField(item.Arn, item.Username, true)
		password := api.GetSecretByIDField(item.Arn, item.Password, true)

		exports := fmt.Sprintf(`export %[1]s_USERNAME=%[2]s
export %[1]s_PASSWORD=%[3]s
`, varName, username, password)
		f.WriteString(exports)
	}
}
