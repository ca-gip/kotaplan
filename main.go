package main

import (
	"flag"
	"fmt"
	"github.com/ca-gip/kotaplan/internal/services/aggregate"
	"github.com/ca-gip/kotaplan/internal/services/k8s"
	"github.com/ca-gip/kotaplan/internal/services/render"
	"github.com/ca-gip/kotaplan/internal/types"
	"github.com/ca-gip/kotaplan/internal/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"log"
	"math"
)

func main() {
	kubeconfig := flag.String("kubeconfig", k8s.DefaultKubeconfig(), "Path to a kubeconfig. Only required if out-of-cluster.")
	masterURL := flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	memClaim := flag.Int64("memclaim", 6, "Amount of memory for the default claim in GiB")
	cpuClaim := flag.Int64("cpuclaim", 1000, "Amount of cpu for the default claim in Milli ex 1000 = 1CPU")
	ratioMemNs := flag.Float64("memratio", 1, "Ratio of the maximum amount of Memory that can be claim by a namespace")
	ratioCpuNs := flag.Float64("cpuratio", 1, "Ratio of the maximum amount of CPU that can be claim by a namespace")
	overcommitMem := flag.Float64("memover", 1, "Ratio of the Memory over or under commit ex: 1 is 100 %")
	overcommitCpu := flag.Float64("cpuover", 1, "Ratio of the CPU over or under commit ex: 1 is 100 %")
	manifest := flag.Bool("manifest", false, "Generate YAML manifest with recommended spec")
	margin := flag.Float64("margin", 1.2, "Margin for the recommanded spec")
	label := flag.String("label", "", "Only match certain namespaces with a label")
	showconf := flag.Bool("showconf", false, "Show the running configuration")
	flag.Parse()

	settings := newSettings(cpuClaim, memClaim, ratioMemNs, ratioCpuNs, overcommitCpu, overcommitMem, manifest, margin, label)

	if *showconf {
		render.SettingTable(settings)
	}

	client, metricsClient := k8s.ClientGen(masterURL, kubeconfig)

	clusterData, err := k8s.GetClusterData(client, metricsClient, settings)

	if err != nil {
		log.Fatalf("Could not gather the required data : %s", err)
	}

	clusterStat := newClusterStat(clusterData, settings)

	render.NamespaceTable(clusterStat)
	render.SummaryTable(clusterStat, settings)

	if *manifest {
		fmt.Println(utils.GenerateClaimManifest(clusterStat))
	}

}

func newSettings(cpuClaim *int64, memClaim *int64, ratioMemNs *float64, ratioCpuNs *float64, overcommitCpu *float64, overcommitMem *float64, manifest *bool, margin *float64, label *string) *types.Settings {
	return &types.Settings{
		DefaultClaim: v1.ResourceList{
			v1.ResourceCPU:    *resource.NewMilliQuantity(*cpuClaim, resource.DecimalSI),
			v1.ResourceMemory: *resource.NewQuantity(int64(float64(*memClaim)*math.Pow(2, 30)), resource.BinarySI),
		},
		RatioMemNs:    *ratioMemNs,
		RatioCpuNs:    *ratioCpuNs,
		OvercommitCpu: *overcommitCpu,
		OvercommitMem: *overcommitMem,
		Manifest:      *manifest,
		Margin:        *margin,
		Labels:        utils.LabelsFromString(*label),
	}
}

func newClusterStat(cluster *types.ClusterData, settings *types.Settings) (stats types.ClusterStat) {

	stats.NamespacesCount = aggregate.CountNs(cluster)
	stats.MemAvailable = aggregate.MemNodes(cluster)
	stats.CpuAvailable = aggregate.CpuNodes(cluster)
	stats.NodesCount = aggregate.CountNodes(cluster)

	for _, namespace := range cluster.Namespaces.Items {
		stats.NamespacesStat = append(stats.NamespacesStat, *newNamespaceStat(namespace, cluster, stats, settings))
	}

	return

}

func newNamespaceStat(namespace v1.Namespace, cluster *types.ClusterData, stats types.ClusterStat, settings *types.Settings) *types.NamespaceStat {

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
