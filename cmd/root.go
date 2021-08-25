package cmd

import (
	"fmt"
	"github.com/ca-gip/kotaplan/internal/services/k8s"
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kotaplan",
	Short: "Visualize resource consumption and generated ResourceQuota with recommend spec",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Config
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kotaplan.yaml)")

	// kubectl
	rootCmd.PersistentFlags().StringP("kubeconfig", "", k8s.DefaultKubeconfig(), "Path to a kubeconfig")
	rootCmd.PersistentFlags().StringP("master", "", "", "Address of the Kubernetes API server. Overrides any value in kubeconfig.")

	// Default claim
	rootCmd.PersistentFlags().Int64P("default-claim-memory", "", 0, "Amount of Memory for the default claim in GiB.")
	rootCmd.PersistentFlags().Int64P("default-claim-cpu", "", 0, "Amount of CPU for the default claim in Milli. ex 1000 = 1CPU")

	// Max per Namespace
	rootCmd.PersistentFlags().Float64P("ratio-namespace-memory", "", 1, "Ratio of the maximum amount of Memory that can be claim by a namespace. Default is 1 : 100% of the cluster is claimable by a Namespace")
	rootCmd.PersistentFlags().Float64P("ratio-namespace-cpu", "", 1, "Ratio of the maximum amount of CPU that can be claim by a namespace. Default is 1 : 100% of the cluster is claimable by a Namespace")

	// Over commit
	rootCmd.PersistentFlags().Float64P("over-commit-memory", "", 1, "Ratio of the Memory over or under commit. Default is 1 meaning 100 %")
	rootCmd.PersistentFlags().Float64P("over-commit-cpu", "", 1, "Ratio of the CPU over or under commit. Default is 1 meaning 100 %")

	// Margin
	rootCmd.PersistentFlags().Float64P("margin", "", 1.2, "Margin for the recommended spec. Default is 1.2")

	// Label
	rootCmd.PersistentFlags().StringP("label", "l", "quota=managed", "Match namespace containing a label. Default is quota=managed")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kotaplan" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kotaplan")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
