package cmd

import (
	"github.com/ca-gip/kotaplan/internal/services/k8s"
	"github.com/ca-gip/kotaplan/internal/services/render"
	"log"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ResourceQuota manifest",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		params := parseParameters(cmd, args)
		k8sClient, metricsClient := initClients(cmd, args)
		clusterData, err := k8s.GetClusterData(k8sClient, metricsClient, params)
		if err != nil {
			log.Fatalf("Could not gather the required data : %s", err)
		}
		clusterStat := newClusterStat(clusterData, params)
		render.Manifest(clusterStat)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

}
