/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wizcas/ame/api"
	"golang.design/x/clipboard"
)

var (
	copySecret bool
)

// getSecretCmd represents the get command
var getSecretCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runGetSecret(args)
	},
}

func init() {
	secretCmd.AddCommand(getSecretCmd)
	getSecretCmd.Flags().BoolVarP(&copySecret, "copy", "c", false, `To copy the value to the clipboard.`)
}

func runGetSecret(args []string) {
	arn := args[0]
	field := "password"

	if len(arn) == 0 {
		log.Fatal("need to provide secret ARN")
	}

	if len(args) > 1 && len(args[1]) > 0 {
		field = args[1]
	}

	var value = api.GetSecretByIDField(arn, field, isURIEncoded)

	if copySecret {
		err := clipboard.Init()
		if err != nil {
			log.Fatal("failed to initialized the clipboard")
		}
		clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf("%v", value)))
		log.Info("✅ The secret value has been copied to the system clipboard.")
		return
	}

	log.Infof("%s", value)
}
