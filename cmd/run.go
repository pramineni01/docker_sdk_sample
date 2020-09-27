/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"docker_sdk_sample/client"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		fmt.Println("Flag passed: ", *killInSecs)
		fmt.Println("Args: ", args)

		var ctx context.Context
		if *killInSecs > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(*killInSecs)*time.Second)
			defer cancel()
		}

		// creates and run docker client within context
		client.New().Run(ctx, args...)
	},
}

var killInSecs *int

func init() {
	rootCmd.AddCommand(runCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	killInSecs = runCmd.Flags().IntP("kill", "k", -1, "kill after specified seconds")
}
