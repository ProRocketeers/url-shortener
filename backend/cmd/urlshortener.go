package cmd

import (
	"fmt"
	"os"

	"github.com/ProRocketeers/url-shortener/internal/infrastructure"
	"github.com/spf13/cobra"
)

var (
	config  infrastructure.ServerConfig
	rootCmd = &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Started with config: %v", config)
			return nil
		},
	}
)

func Execute() {
	cfg, err := infrastructure.ParseServerConfig()
	if err != nil {
		fmt.Printf("error parsing server config: %v", err)
		os.Exit(1)
	}
	config = cfg
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error while executing root command: %v", err)
		os.Exit(1)
	}
}
