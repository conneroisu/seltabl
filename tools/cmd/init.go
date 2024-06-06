package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/seltabl/tools/internal"
	"github.com/conneroisu/seltabl/tools/internal/config"
	"github.com/spf13/viper"
)

var cfgFile string

// init is the entry point for the application
func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.seltabl.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ctx := context.Background()
	internal.AddRoutes(ctx, rootCmd, &config.Config{})
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("seltabl.yaml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
