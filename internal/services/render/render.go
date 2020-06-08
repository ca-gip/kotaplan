package render

import (
	"github.com/ca-gip/kotaplan/internal/services/aggregate"
	"github.com/ca-gip/kotaplan/internal/types"
	"github.com/ca-gip/kotaplan/internal/utils"
	"github.com/jedib0t/go-pretty/table"
	"math"
	"os"
)

func SettingTable(settings *types.Settings) {
	t := table.NewWriter()
	t.SetTitle("Running Settings")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"name", "value"})
	t.AppendRow(table.Row{"Default Claim", utils.ResourceListToString(settings.DefaultClaim)})
	t.AppendRow(table.Row{"memratio", utils.ResourceListToString(settings.DefaultClaim)})
	t.Render()
}

func NamespaceTable(cluster types.ClusterStat) {

	t := table.NewWriter()
	t.SetTitle("Namespaces Details")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"namespace", "pods", "mem req", "current mem use", "mem req usage", "cpu req", "current cpu use", "cpu req usage", "fit default", "respect max allocation ns", "spec"})

	for _, ns := range cluster.NamespacesStat {
		t.AppendRow(table.Row{
			ns.Name,
			ns.PodCount,
			utils.ByteToString(ns.MemReq),
			utils.ByteToString(ns.MemUse),
			ns.MemReqUse,
			ns.CpuReq,
			ns.CpuUse,
			ns.CpuReqUse,
			ns.ClaimFit,
			ns.RespectMaxNS,
			utils.ResourceListToString(ns.Spec),
		})
	}

	t.AppendFooter(table.Row{
		cluster.NamespacesCount,
		aggregate.PodsSum(cluster),
		"", "", "", "", "", "", "", "",
		utils.SpecFromIntToString(aggregate.ClaimMemSum(cluster), aggregate.ClaimCpuSum(cluster)),
	})

	t.Render()

	return

}

func SummaryTable(cluster types.ClusterStat, settings *types.Settings) {
	t := table.NewWriter()
	t.Style().Options.SeparateColumns = false
	t.SetTitle("Summary")
	t.SetOutputMirror(os.Stdout)
	t.AppendRow(table.Row{"Number of nodes", cluster.NodesCount})
	t.AppendRow(table.Row{"Available resources (real)", utils.SpecFromIntToString(cluster.MemAvailable, cluster.CpuAvailable)})
	t.AppendRow(table.Row{"Available resources (commit)", utils.SpecFromIntToString(int64(math.Round(float64(cluster.MemAvailable)*settings.OvercommitMem)), int64(float64(cluster.CpuAvailable)*settings.OvercommitCpu))})
	t.AppendRow(table.Row{"Max per NS", utils.SpecFromIntToString(int64(float64(cluster.MemAvailable)*settings.RatioMemNs), int64(float64(cluster.CpuAvailable)*settings.RatioCpuNs))})
	t.AppendFooter(table.Row{"Result", utils.Result(cluster, settings)})
	t.Render()
}
