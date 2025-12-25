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

// MCP Protocol definitions
type MCPContext struct {
    Version     string
    Environment string
    Toolsets    []string
    Config      *MCPConfig
}

type MCPConfig struct {
    IsReadOnly    bool
    IsDynamic     bool
    EnableLogging bool
    Host          string
}

// MCP response interface
type MCPResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Error   string      `json:"error,omitempty"`
}

// These variables are set by the build process using ldflags.
var version = "version"
var commit = "commit"
var date = "date"

// Global MCP context
var mcpContext *MCPContext

var (
    rootCmd = &cobra.Command{
        Use:     "server",
        Short:   "GitHub MCP Server",
        Long:    `A GitHub MCP server that handles various tools and resources.`,
        Version: fmt.Sprintf("Version: %s\nCommit: %s\nBuild Date: %s", version, commit, date),
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
            // Initialize MCP context
            mcpContext = &MCPContext{
                Version:     version,
                Environment: os.Getenv("GITHUB_ENVIRONMENT"),
                Config:      &MCPConfig{},
            }
        },
    }

    stdioCmd = &cobra.Command{
        Use:   "stdio",
        Short: "Start stdio server",
        Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
        RunE: func(_ *cobra.Command, _ []string) error {
            token := viper.GetString("personal_access_token")
            if token == "" {
                return createMCPError("GITHUB_PERSONAL_ACCESS_TOKEN not set")
            }

            var enabledToolsets []string
            if err := viper.UnmarshalKey("toolsets", &enabledToolsets); err != nil {
                return createMCPError(fmt.Sprintf("failed to unmarshal toolsets: %v", err))
            }

            // Update MCP context
            mcpContext.Toolsets = enabledToolsets
            mcpContext.Config = &MCPConfig{
                IsReadOnly:    viper.GetBool("read-only"),
                IsDynamic:     viper.GetBool("dynamic_toolsets"),
                EnableLogging: viper.GetBool("enable-command-logging"),
                Host:          viper.GetString("host"),
            }

            stdioServerConfig := ghmcp.StdioServerConfig{
                Version:              version,
                Host:                 mcpContext.Config.Host,
                Token:                token,
                EnabledToolsets:      mcpContext.Toolsets,
                DynamicToolsets:      mcpContext.Config.IsDynamic,
                ReadOnly:             mcpContext.Config.IsReadOnly,
                ExportTranslations:   viper.GetBool("export-translations"),
                EnableCommandLogging: mcpContext.Config.EnableLogging,
                LogFilePath:          viper.GetString("log-file"),
            }

            return handleMCPResponse(ghmcp.RunStdioServer(stdioServerConfig))
        },
    }
)

// MCP helper functions
func createMCPError(message string) error {
    return &MCPResponse{
        Success: false,
        Error:   message,
    }
}

func handleMCPResponse(err error) error {
    if err != nil {
        return &MCPResponse{
            Success: false,
            Error:   err.Error(),
        }
    }
    return &MCPResponse{
        Success: true,
        Data:    mcpContext,
    }
}

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

    // Bind flags to viper
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
    viper.SetEnvPrefix("github")
    viper.AutomaticEnv()
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        response, ok := err.(*MCPResponse)
        if ok {
            fmt.Fprintf(os.Stderr, "Error: %s\n", response.Error)
        } else {
            fmt.Fprintf(os.Stderr, "%v\n", err)
        }
        os.Exit(1)
    }
}
