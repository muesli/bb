package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GlobalOptions struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

var (
	rootCmd = &cobra.Command{
		Use: "bb",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			viper.Unmarshal(&globalOpts)
		},
	}

	cfgFile    string
	globalOpts = GlobalOptions{}

	username string
	password string
	repoOrga string
	repoSlug string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bb)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "app password")
	rootCmd.PersistentFlags().StringVar(&repoOrga, "repo-orga", "", "repository organisation")
	rootCmd.PersistentFlags().StringVar(&repoSlug, "repo-slug", "", "repository slug")

	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("repoOrga", rootCmd.PersistentFlags().Lookup("repo-orga"))
	viper.BindPFlag("repoSlug", rootCmd.PersistentFlags().Lookup("repo-slug"))
}

func initConfig() {
	if cfgFile == "" {
		configDir := configdir.LocalConfig("bb")
		err := configdir.MakePath(configDir)
		if err != nil {
			panic(err)
		}
		cfgFile = filepath.Join(configDir, "configuration.toml")
		if _, err = os.Stat(cfgFile); os.IsNotExist(err) {
			fh, err := os.Create(cfgFile)
			if err != nil {
				panic(err)
			}
			defer fh.Close()
		}
	}

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
