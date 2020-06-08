package types

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type ClusterData struct {
	Nodes      []v1.Node
	Namespaces *v1.NamespaceList
	Pods       *v1.PodList
	PodsMetric *v1beta1.PodMetricsList
}

type NamespaceStat struct {
	Name         string
	PodCount     int
	MemReq       int64
	MemUse       int64
	MemReqUse    string
	CpuReq       int64
	CpuUse       int64
	CpuReqUse    string
	ClaimFit     bool
	RespectMaxNS bool
	Spec         v1.ResourceList
}

type ClusterStat struct {
	MemAvailable    int64
	CpuAvailable    int64
	NodesCount      int
	NamespacesCount int
	NamespacesStat  []NamespaceStat
}

type Settings struct {
	DefaultClaim  v1.ResourceList
	RatioMemNs    float64
	RatioCpuNs    float64
	OvercommitCpu float64
	OvercommitMem float64
	Manifest      bool
	Margin        float64
	Labels        map[string]string
}
