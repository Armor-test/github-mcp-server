package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/github/github-mcp-server/internal/ghmcp"
	"github.com/github/github-mcp-server/pkg/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// These variables are set by the build process using ldflags.

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate("{{.Short}}\n{{.Version}}\n")

	// Add global flags that will be shared by all commands
	rootCmd.PersistentFlags().StringSlice("toolsets", github.DefaultTools, "An optional comma separated list of groups of tools to allow, defaults to enabling all")
	rootCmd.PersistentFlags().Bool("dynamic-toolsets", false, "Enable dynamic toolsets")
	rootCmd.PersistentFlags().Bool("read-only", false, "Restrict the server to read-only operations")
	rootCmd.PersistentFlags().String("log-file", "", "Path to log file")
	rootCmd.PersistentFlags().Bool("enable-command-logging", false, "When enabled, the server will log all command requests and responses to the log file")
	rootCmd.PersistentFlags().Bool("export-translations", false, "Save translations to a JSON file")
	rootCmd.PersistentFlags().String("gh-host", "", "Specify the GitHub hostname (for GitHub Enterprise etc.)")

	// Bind flag to viper
	_ = viper.BindPFlag("toolsets", rootCmd.PersistentFlags().Lookup("toolsets"))
	_ = viper.BindPFlag("dynamic_toolsets", rootCmd.PersistentFlags().Lookup("dynamic-toolsets"))
	_ = viper.BindPFlag("read-only", rootCmd.PersistentFlags().Lookup("read-only"))
	_ = viper.BindPFlag("log-file", rootCmd.PersistentFlags().Lookup("log-file"))
	_ = viper.BindPFlag("enable-command-logging", rootCmd.PersistentFlags().Lookup("enable-command-logging"))
	_ = viper.BindPFlag("export-translations", rootCmd.PersistentFlags().Lookup("export-translations"))
	_ = viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("gh-host"))

	// Add subcommands
	rootCmd.AddCommand(stdioCmd)
}

func initConfig() {
	// Initialize Viper configuration
	viper.SetEnvPrefix("github")
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}type ApprovalWorkflow struct {
    RequiresMFA      bool
    ApprovalLevels   int
    BypassConditions map[string]bool
}

func main() {
    http.HandleFunc("/api/approve-transaction", handleTransactionApproval)
    http.HandleFunc("/api/user/mfa", handleMFAValidation)
    http.HandleFunc("/api/admin/access", handleAdminAccess)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
