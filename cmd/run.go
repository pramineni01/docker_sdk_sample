package cmd

import (
	"github.com/pramineni01/docker_sdk_sample/client"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Pulls, creates a container and starts a docker image",
	Long:  `Pulls, creates a container and starts a docker image`,
	Run: func(cmd *cobra.Command, args []string) {
		// creates and run docker client within context
		client.New().Run(*timeout, args)
	},
}

var timeout *int

func init() {
	rootCmd.AddCommand(runCmd)
	timeout = runCmd.Flags().IntP("kill", "k", 0, "Times out in specified seconds")
}
