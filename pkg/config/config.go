package config

import (
	"errors"
	"fmt"
	"strings"

	logrusformat "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/rkrmr33/onka/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	gConfigLoader configLoader
)

type (
	ConfigLoader interface {
		Load() error
	}

	configLoader struct {
		defaultConfigPath string
		defaultConfigName string

		configPath string
		logLvl     string
		rootCmd    *cobra.Command
		mappings   map[string]string
	}
)

// Init initializes the global config loader
func Init(rootCmd *cobra.Command, defaultConfigPath, defaultConfigName string) {
	gConfigLoader = configLoader{
		defaultConfigPath: defaultConfigPath,
		defaultConfigName: defaultConfigName,

		rootCmd:  rootCmd,
		mappings: map[string]string{},
	}

	rootCmd.Flags().StringVar(&gConfigLoader.logLvl, "log-level", log.GetLevel().String(), "log level")
	rootCmd.Flags().StringVar(&gConfigLoader.configPath, "config-path", defaultConfigPath, "the config path")
}

func Load() error {
	initViper()
	initLogger()

	if gConfigLoader.configPath != gConfigLoader.defaultConfigPath {
		viper.SetConfigFile(gConfigLoader.configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Debug("config file not found")
			return nil
		}
		return err
	}

	log.Debugf("config file loaded from: %s", viper.ConfigFileUsed())

	bindFlags()

	return nil
}

func BindUnder(prefix string, set *pflag.FlagSet) {
	set.VisitAll(func(f *pflag.Flag) {
		var configName string
		name := kebabToLowerCamel(f.Name)
		if prefix == "" {
			configName = name
		} else {
			configName = fmt.Sprintf("%s.%s", prefix, name)
		}
		BindConfig(configName, f)
	})
	gConfigLoader.rootCmd.Flags().AddFlagSet(set)
}

func BindConfig(configName string, flag *pflag.Flag) {
	gConfigLoader.mappings[flag.Name] = configName
	util.Must(viper.BindPFlag(configName, flag))
}

func kebabToLowerCamel(s string) string {
	words := strings.Split(s, "-")
	for i := range words {
		if i == 0 {
			continue
		}
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func initViper() {
	viper.SetConfigName(gConfigLoader.defaultConfigName)
	viper.AddConfigPath(".")
	viper.AddConfigPath(gConfigLoader.defaultConfigPath)
}

func initLogger() error {
	l, err := log.ParseLevel(gConfigLoader.logLvl)
	if err != nil {
		return err
	}

	if l == log.DebugLevel {
		log.StandardLogger().Formatter = &logrusformat.Formatter{
			ChildFormatter: &log.TextFormatter{FullTimestamp: true},
			Line:           true,
		}
	}

	log.SetLevel(l)
	return nil
}

func bindFlags() {
	gConfigLoader.rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		cfgName, ok := gConfigLoader.mappings[f.Name]
		if !ok {
			// Flag unbound to a viper key
			return
		}

		if !f.Changed && viper.IsSet(cfgName) {
			val := viper.Get(cfgName)
			util.Must(gConfigLoader.rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)))
		}
	})
}
