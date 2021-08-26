package cmd

import (
	"fmt"
	"github.com/ca-gip/kotaplan/internal/services/aggregate"
	"github.com/ca-gip/kotaplan/internal/services/k8s"
	"github.com/ca-gip/kotaplan/internal/types"
	"github.com/ca-gip/kotaplan/internal/utils"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"math"
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
	rootCmd.PersistentFlags().Int64P("default-claim-memory", "", 0, "Amount of Memory for the default claim in GiB. Default is 0.")
	rootCmd.PersistentFlags().Int64P("default-claim-cpu", "", 0, "Amount of CPU for the default claim in Milli. ex 1000 = 1CPU. Default is 0.")

	// Max per Namespace
	rootCmd.PersistentFlags().Float64P("ratio-namespace-memory", "", 1, "Ratio of the maximum amount of Memory that can be claim by a namespace. Default is 1 : 100% of the cluster is claimable by a Namespace")
	rootCmd.PersistentFlags().Float64P("ratio-namespace-cpu", "", 1, "Ratio of the maximum amount of CPU that can be claim by a namespace. Default is 1 : 100% of the cluster is claimable by a Namespace")

	// Over commit
	rootCmd.PersistentFlags().Float64P("over-commit-memory", "", 1, "Ratio of the Memory over or under commit. Default is 1 meaning 100 %")
	rootCmd.PersistentFlags().Float64P("over-commit-cpu", "", 1, "Ratio of the CPU over or under commit. Default is 1 meaning 100 %")

	// Margin
	rootCmd.PersistentFlags().Float64P("margin", "", 1.2, "Margin for the recommended spec. Default is 1.2")

	// Label
	rootCmd.PersistentFlags().StringP("labels", "l", "quota=managed", "Match namespace containing a label. Default is quota=managed")
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

func initClients(cmd *cobra.Command, args []string) (client *kubernetes.Clientset, metricsClient *versioned.Clientset) {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	master, _ := cmd.Flags().GetString("master")
	return k8s.ClientGen(&master, &kubeconfig)
}

func parseParameters(cmd *cobra.Command, args []string) *types.Parameters {
	defaultClaimMemory, _ := cmd.Flags().GetInt64("default-claim-memory")
	defaultClaimCpu, _ := cmd.Flags().GetInt64("default-claim-cpu")

	ratioNamespaceMemory, _ := cmd.Flags().GetFloat64("ratio-namespace-memory")
	ratioNamespaceCpu, _ := cmd.Flags().GetFloat64("ratio-namespace-cpu")

	overCommitMemory, _ := cmd.Flags().GetFloat64("over-commit-memory")
	overCommitCpu, _ := cmd.Flags().GetFloat64("over-commit-cpu")

	margin, _ := cmd.Flags().GetFloat64("margin")

	labels, _ := cmd.Flags().GetString("labels")

	return &types.Parameters{
		DefaultClaim: v1.ResourceList{
			v1.ResourceCPU:    *resource.NewMilliQuantity(defaultClaimCpu, resource.DecimalSI),
			v1.ResourceMemory: *resource.NewQuantity(int64(float64(defaultClaimMemory)*math.Pow(2, 30)), resource.BinarySI),
		},
		RatioNsMemory:    ratioNamespaceMemory,
		RatioNsCpu:       ratioNamespaceCpu,
		OverCommitMemory: overCommitMemory,
		OverCommitCpu:    overCommitCpu,
		Margin:           margin,
		Labels:           utils.LabelsFromString(labels),
	}
}

func newClusterStat(cluster *types.ClusterData, settings *types.Parameters) (stats types.ClusterStat) {

	stats.NamespacesCount = aggregate.CountNs(cluster)
	stats.MemAvailable = aggregate.MemNodes(cluster)
	stats.CpuAvailable = aggregate.CpuNodes(cluster)
	stats.NodesCount = aggregate.CountNodes(cluster)

	for _, namespace := range cluster.Namespaces.Items {
		stats.NamespacesStat = append(stats.NamespacesStat, *newNamespaceStat(namespace, cluster, stats, settings))
	}

	return

}

func newNamespaceStat(namespace v1.Namespace, cluster *types.ClusterData, stats types.ClusterStat, settings *types.Parameters) *types.NamespaceStat {

	memReq := aggregate.MemRequestSumByNS(cluster.Pods, namespace)
	memUse := aggregate.MemUsageByNS(cluster.PodsMetric, namespace)
	cpuReq := aggregate.CpuRequestSumByNS(cluster.Pods, namespace)
	cpuUse := aggregate.CpuUsageByNS(cluster.PodsMetric, namespace)
	claimFit, spec := utils.CheckSpec(memReq, cpuReq, settings)
	respectMaxNS := utils.CheckRespectMaxNS(spec, stats, settings)

	return &types.NamespaceStat{
		Name:         namespace.Name,
		PodCount:     aggregate.PodCountByNS(cluster.Pods, namespace),
		MemReq:       memReq,
		MemUse:       memUse,
		MemReqUse:    utils.DivideAsPercent(memUse, memReq),
		CpuReq:       cpuReq,
		CpuUse:       cpuUse,
		CpuReqUse:    utils.DivideAsPercent(cpuUse, cpuReq),
		ClaimFit:     claimFit,
		RespectMaxNS: respectMaxNS,
		Spec:         spec,
	}

}
