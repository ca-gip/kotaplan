package aggregate

import (
	underscore "github.com/ahl5esoft/golang-underscore"
	"github.com/ca-gip/kotaplan/internal/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func ClaimMemSum(cluster types.ClusterStat) (result int64) {
	underscore.
		Chain(cluster.NamespacesStat).
		Map(func(ns types.NamespaceStat, _ int) int64 { return ns.Spec.Memory().Value() }).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)
	return
}

func ClaimCpuSum(cluster types.ClusterStat) (result int64) {
	underscore.
		Chain(cluster.NamespacesStat).
		Map(func(ns types.NamespaceStat, _ int) int64 { return ns.Spec.Cpu().MilliValue() }).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)
	return
}

func PodsSum(cluster types.ClusterStat) (result int) {
	underscore.
		Chain(cluster.NamespacesStat).
		Map(func(ns types.NamespaceStat, _ int) int { return ns.PodCount }).
		Aggregate(0, func(acc int, cur int, _ int) int { return acc + cur }).
		Value(&result)
	return
}

func PodCountByNS(pods *v1.PodList, ns v1.Namespace) int {
	return underscore.Chain(pods.Items).WhereBy(map[string]interface{}{
		"Namespace": ns.Name,
	}).Count()
}

func MemRequestSumByNS(pods *v1.PodList, ns v1.Namespace) (result int64) {
	underscore.
		Chain(pods.Items).
		WhereBy(map[string]interface{}{"Namespace": ns.Name}).
		Map(func(pod v1.Pod, _ int) []v1.Container { return pod.Spec.Containers }).
		Map(func(containers []v1.Container, _ int) (totalMemReq int64) {
			underscore.
				Chain(containers).
				Aggregate(int64(0), func(acc int64, cur v1.Container, _ int) int64 { return acc + cur.Resources.Requests.Memory().Value() }).
				Value(&totalMemReq)
			return
		}).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)
	return
}

func CpuRequestSumByNS(pods *v1.PodList, ns v1.Namespace) (result int64) {
	underscore.
		Chain(pods.Items).
		WhereBy(map[string]interface{}{"Namespace": ns.Name}).
		Map(func(pod v1.Pod, _ int) []v1.Container { return pod.Spec.Containers }).
		Map(func(containers []v1.Container, _ int) (totalCpuReq int64) {
			underscore.
				Chain(containers).
				Aggregate(int64(0), func(acc int64, cur v1.Container, _ int) int64 { return acc + cur.Resources.Requests.Cpu().MilliValue() }).
				Value(&totalCpuReq)
			return
		}).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)
	return
}

func MemUsageByNS(podsMetric *v1beta1.PodMetricsList, ns v1.Namespace) (result int64) {
	underscore.
		Chain(podsMetric.Items).
		WhereBy(map[string]interface{}{"Namespace": ns.Name}).
		Map(func(podMetric v1beta1.PodMetrics, _ int) []v1beta1.ContainerMetrics { return podMetric.Containers }).
		Map(func(containersMetric []v1beta1.ContainerMetrics, _ int) (totalMemUsage int64) {
			underscore.
				Chain(containersMetric).
				Aggregate(int64(0), func(acc int64, cur v1beta1.ContainerMetrics, _ int) int64 { return acc + cur.Usage.Memory().Value() }).
				Value(&totalMemUsage)
			return
		}).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)
	return
}

func CpuUsageByNS(podsMetric *v1beta1.PodMetricsList, ns v1.Namespace) (result int64) {
	underscore.
		Chain(podsMetric.Items).
		WhereBy(map[string]interface{}{"Namespace": ns.Name}).
		Map(func(podMetric v1beta1.PodMetrics, _ int) []v1beta1.ContainerMetrics { return podMetric.Containers }).
		Map(func(containersMetric []v1beta1.ContainerMetrics, _ int) (totalCpuUsage int64) {
			underscore.
				Chain(containersMetric).
				Aggregate(int64(0), func(acc int64, cur v1beta1.ContainerMetrics, _ int) int64 { return acc + cur.Usage.Cpu().MilliValue() }).
				Value(&totalCpuUsage)
			return
		}).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)
	return
}

func MemNodes(data *types.ClusterData) (result int64) {
	underscore.
		Chain(data.Nodes).
		Map(func(node v1.Node, _ int) int64 { return node.Status.Allocatable.Memory().Value() }).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)

	return
}

func CpuNodes(data *types.ClusterData) (result int64) {
	underscore.
		Chain(data.Nodes).
		Map(func(node v1.Node, _ int) int64 { return node.Status.Allocatable.Cpu().MilliValue() }).
		Aggregate(int64(0), func(acc int64, cur int64, _ int) int64 { return acc + cur }).
		Value(&result)
	return
}

func CountNodes(data *types.ClusterData) int {
	return underscore.
		Chain(data.Nodes).
		Count()
}

func CountNs(data *types.ClusterData) int {
	return underscore.Chain(data.Namespaces.Items).Count()
}
