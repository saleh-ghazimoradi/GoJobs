package cmd

import (
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"github.com/saleh-ghazimoradi/GoJobs/logger"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "GoJobs",
	Short: "A brief description of your application",
}

func Execute() {
	err := os.Setenv("TZ", time.UTC.String())
	if err != nil {
		panic(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	err := config.LoadingConfig()
	if err != nil {
		logger.Logger.Error("there went something wrong while loading config file")
	}
}
