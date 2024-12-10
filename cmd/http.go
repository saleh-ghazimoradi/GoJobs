package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/GoJobs/internal/gateway"

	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http called")
		if err := gateway.Server(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
