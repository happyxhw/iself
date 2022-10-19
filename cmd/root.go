package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/happyxhw/iself/pkg/log"
	"github.com/happyxhw/iself/service"
)

var (
	cfgFile string
	env     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iself",
	Short: "iself",
	Long:  `iself`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config", "配置文件目录")
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "local", "环境: dev, test, pre, prod")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath(cfgFile)
	viper.SetConfigName(strings.ToLower(env))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("read config err: %+v\n", err)
		os.Exit(1)
	}
}

func start() {
	log.InitAppLogger(
		&log.Config{Level: viper.GetString("log.app.level"),
			Encoder: viper.GetString("log.encoder")},
		zap.AddCallerSkip(1), zap.AddCaller())

	service.Serve()
}
