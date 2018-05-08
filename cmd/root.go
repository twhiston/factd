package cmd

import (
	"fmt"
	"os"

	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twhiston/factd/lib/common/logging"
	"github.com/twhiston/factd/lib/factd"
	formatter2 "github.com/twhiston/factd/lib/formatter"
	plugins2 "github.com/twhiston/factd/lib/plugin"
	"github.com/twhiston/factd/lib/plugin/cpu"
	"github.com/twhiston/factd/lib/plugin/disks"
	"github.com/twhiston/factd/lib/plugin/docker"
	"github.com/twhiston/factd/lib/plugin/host"
	"github.com/twhiston/factd/lib/plugin/load"
	"github.com/twhiston/factd/lib/plugin/mem"
	"github.com/twhiston/factd/lib/plugin/net"
	"github.com/twhiston/factd/lib/plugin/user"
)

var cfgFile string

// List of all available plugin.
// To add a new plugin add it to this slice
var factdPlugins = []plugins2.Plugin{
	new(cpu.CPU),
	new(disks.Disks),
	new(host.Host),
	new(mem.Mem),
	new(net.Net),
	new(user.User),
	new(docker.Docker),
	new(load.Load),
	new(plugins2.Version),
}

// internal variable used for plugin resolving
var pluginMap = map[string]plugins2.Plugin{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "factd",
	Short: "Facts Daemon",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.facter.yaml)")

	// These flags are also bound to viper so that they can be set in a config file as well as here
	// Throughout the cli program the values will only be taken from viper and not directly from the flags
	rootCmd.PersistentFlags().String("format", "yaml", "plaintext/json/yaml")
	rootCmd.PersistentFlags().StringSlice("include", []string{}, "what plugin to run")
	rootCmd.PersistentFlags().StringSlice("exclude", []string{}, "what plugin to exclude")

	//make a map out of the plugin list for easier selection and filtering
	for _, v := range factdPlugins {
		pluginMap[v.Name()] = v
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error().Err(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".factd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	//bind root command values
	logging.Fatal(viper.BindPFlag("format", rootCmd.Flags().Lookup("format")))
	logging.Fatal(viper.BindPFlag("include", rootCmd.Flags().Lookup("include")))
	logging.Fatal(viper.BindPFlag("exclude", rootCmd.Flags().Lookup("exclude")))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}

func setupFactD() *factd.Factd {
	conf := factd.NewConfig()
	fname := viper.GetString("format")
	resolveFormatter(conf, fname)
	resolvePlugins(conf)
	return factd.New(*conf)
}

// TODO - make this less crappy
func resolveFormatter(c *factd.Config, cName string) {

	switch cName {
	case "plaintext":
		c.Formatter = formatter2.NewPlainTextFormatter()
	case "json":
		c.Formatter = formatter2.NewJSONFormatter()
	case "yaml":
		c.Formatter = formatter2.NewYAMLFormatter()
	default:
		logging.Fatal(errors.New("cannot find formatter"))
	}
}

func resolvePlugins(c *factd.Config) {

	//Add plugin
	included := viper.GetStringSlice("include")
	if len(included) == 0 || included[0] == "all" {
		//Add all plugin
		for k, v := range pluginMap {
			c.Plugins[k] = v
		}
	} else {
		for _, v := range included {
			if val, ok := pluginMap[v]; ok {
				c.Plugins[val.Name()] = val
			}
		}
	}

	//Remove plugin
	excluded := viper.GetStringSlice("exclude")
	if len(excluded) > 0 {
		for _, v := range excluded {
			if _, ok := c.Plugins[v]; ok {
				delete(c.Plugins, v)
			}
		}
	}

}
