package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var kubeconfigPath string
var config *api.Config

type AliasConfig struct {
	Aliases map[string]string `yaml:"aliases"`
}

var aliasFile = filepath.Join(getHomeDir(), ".kubex", "aliases.yaml")

var rootCmd = &cobra.Command{
	Use:   "kubex",
	Short: "A simple Kubernetes context manager",
	Long:  `kubex is a stripped-down version of kubectx for managing Kubernetes contexts.`,
}

func main() {

	rootCmd.PersistentFlags().StringVar(&kubeconfigPath, "kubeconfig", os.Getenv("KUBECONFIG"), "Path to kubeconfig file")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(completionCmd)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if kubeconfigPath == "" {
			if home := os.Getenv("HOME"); home != "" {
				kubeconfigPath = home + "/.kube/config"
			} else if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
				kubeconfigPath = userProfile + "\\.kube\\config"
			} else {
				fmt.Println("Error: Unable to determine kubeconfig path.")
				os.Exit(1)
			}
		}
		var err error
		config, err = clientcmd.LoadFromFile(kubeconfigPath)
		if err != nil {
			return fmt.Errorf("error loading kubeconfig: %v", err)
		}
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper function to get the home directory (cross-platform)
func getHomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		return userProfile
	}
	fmt.Println("Error: Unable to determine home directory.")
	os.Exit(1)
	return ""
}

// Load aliases from the alias file
func loadAliases() (AliasConfig, error) {
	var config AliasConfig
	config.Aliases = make(map[string]string)

	if _, err := os.Stat(aliasFile); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(aliasFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

// Save aliases to the alias file
func saveAliases(config AliasConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	dir := filepath.Dir(aliasFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(aliasFile, data, 0644)
}

// List Command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contexts",
	Run: func(cmd *cobra.Command, args []string) {
		listContexts()
	},
}

func listContexts() {
	fmt.Println("Available contexts:")
	for context := range config.Contexts {
		if context == config.CurrentContext {
			fmt.Printf("- %s (active)\n", context)
		} else {
			fmt.Printf("- %s\n", context)
		}
	}
}

// Use Command
var useCmd = &cobra.Command{
	Use:   "use [context]",
	Short: "Switch to a specific context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switchContext(args[0])
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var suggestions []string

		// Add contexts
		for context := range config.Contexts {
			suggestions = append(suggestions, context)
		}

		// Add aliases
		aliasConfig, _ := loadAliases()
		for alias := range aliasConfig.Aliases {
			suggestions = append(suggestions, alias)
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
}

func switchContext(contextName string) {
	// Load aliases
	aliasConfig, err := loadAliases()
	if err != nil {
		fmt.Printf("Error loading aliases: %v\n", err)
		return
	}

	// Resolve alias
	if resolvedContext, exists := aliasConfig.Aliases[contextName]; exists {
		contextName = resolvedContext
	}

	if _, exists := config.Contexts[contextName]; !exists {
		fmt.Printf("Error: Context '%s' not found.\n", contextName)
		return
	}

	config.CurrentContext = contextName
	err = clientcmd.WriteToFile(*config, kubeconfigPath)
	if err != nil {
		fmt.Printf("Error updating kubeconfig: %v\n", err)
		return
	}

	fmt.Printf("Switched to context: %s\n", contextName)
}

// Current Command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current context",
	Run: func(cmd *cobra.Command, args []string) {
		showCurrentContext()
	},
}

func showCurrentContext() {
	fmt.Printf("Current context: %s\n", config.CurrentContext)
}

// Alias Command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage context aliases",
}

var aliasAddCmd = &cobra.Command{
	Use:   "add <alias> <context>",
	Short: "Add an alias for a context",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		context := args[1]

		config, err := loadAliases()
		if err != nil {
			fmt.Printf("Error loading aliases: %v\n", err)
			return
		}

		if _, exists := config.Aliases[alias]; exists {
			fmt.Printf("Alias '%s' already exists. Overwrite? [y/N]: ", alias)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Operation canceled.")
				return
			}
		}

		config.Aliases[alias] = context
		if err := saveAliases(config); err != nil {
			fmt.Printf("Error saving aliases: %v\n", err)
			return
		}

		fmt.Printf("Added alias: %s -> %s\n", alias, context)
	},
}

var aliasRemoveCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove an alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]

		config, err := loadAliases()
		if err != nil {
			fmt.Printf("Error loading aliases: %v\n", err)
			return
		}

		if _, exists := config.Aliases[alias]; !exists {
			fmt.Printf("Error: Alias '%s' not found.\n", alias)
			return
		}

		delete(config.Aliases, alias)
		if err := saveAliases(config); err != nil {
			fmt.Printf("Error saving aliases: %v\n", err)
			return
		}

		fmt.Printf("Removed alias: %s\n", alias)
	},
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadAliases()
		if err != nil {
			fmt.Printf("Error loading aliases: %v\n", err)
			return
		}

		if len(config.Aliases) == 0 {
			fmt.Println("No aliases defined.")
			return
		}

		fmt.Println("Aliases:")
		for alias, context := range config.Aliases {
			fmt.Printf("- %s -> %s\n", alias, context)
		}
	},
}

func init() {
	aliasCmd.AddCommand(aliasAddCmd)
	aliasCmd.AddCommand(aliasRemoveCmd)
	aliasCmd.AddCommand(aliasListCmd)
	rootCmd.AddCommand(aliasCmd)
}

// Completion Command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate autocompletion script for the specified shell",
	Long: `To enable autocompletion, run the following command:
For Bash: source <(kubex completion bash)
For Zsh: source <(kubex completion zsh)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		default:
			fmt.Println("Unsupported shell. Use 'bash', 'zsh', or 'fish'.")
		}
	},
}
