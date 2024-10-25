package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/flexwie/go-common/logger"
	"github.com/flexwie/owntracks-api/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "owntracks",
	Long:    "Minimal replication of the OwnTracks reporter connected to a SQL storage",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
		
                     __                        __                                   .__ 
  ______  _  _______/  |_____________    ____ |  | __  ______         _____  ______ |__|
 /  _ \ \/ \/ /    \   __\_  __ \__  \ _/ ___\|  |/ / /  ___/  ______ \__  \ \____ \|  |
(  <_> )     /   |  \  |  |  | \// __ \\  \___|    <  \___ \  /_____/  / __ \|  |_> >  |
 \____/ \/\_/|___|  /__|  |__|  (____  /\___  >__|_ \/____  >         (____  /   __/|__|
                  \/                 \/     \/     \/     \/               \/|__|       

		`)

		for key, value := range viper.GetViper().AllSettings() {
			fmt.Printf("%s: %v\n", key, value)
		}
		fmt.Println()

		if viper.GetBool("verbose") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		fx.New(
			logger.WithLoggerFactory(log.GetLevel()),
			logger.WithFxLogger,
			internal.WithBusinessLogic(),
		).Run()
	},
}

func init() {
	rootCmd.PersistentFlags().String("addr", ":8080", "address")
	viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))

	rootCmd.PersistentFlags().String("db-host", "localhost", "hostname of db")
	viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))
	rootCmd.PersistentFlags().String("db-user", "postgres", "user for db")
	viper.BindPFlag("db-user", rootCmd.PersistentFlags().Lookup("db-user"))
	rootCmd.PersistentFlags().String("db-password", "", "db user password")
	viper.BindPFlag("db-password", rootCmd.PersistentFlags().Lookup("db-password"))
	rootCmd.PersistentFlags().String("db-name", "", "name of the database")
	viper.BindPFlag("db-name", rootCmd.PersistentFlags().Lookup("db-name"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "set log level to debug")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
