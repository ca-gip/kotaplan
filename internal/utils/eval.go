package utils

import (
	"bytes"
	"fmt"
	"github.com/ca-gip/kotaplan/internal/services/aggregate"
	"github.com/ca-gip/kotaplan/internal/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"math"
)

func recommendedSpec(memReq int64, cpuReq int64, settings *types.Settings) v1.ResourceList {

	if memReq > settings.DefaultClaim.Memory().Value() {
		memReq = int64(math.Round(float64(memReq) * settings.Margin))
	} else {
		memReq = settings.DefaultClaim.Memory().Value()
	}

	if cpuReq > settings.DefaultClaim.Cpu().MilliValue() {
		cpuReq = int64(math.Round(float64(cpuReq) * settings.Margin))
	} else {
		cpuReq = settings.DefaultClaim.Cpu().MilliValue()
	}

	return v1.ResourceList{
		v1.ResourceCPU:    *resource.NewMilliQuantity(cpuReq, resource.DecimalSI),
		v1.ResourceMemory: *resource.NewQuantity(memReq, resource.BinarySI),
	}
}

func defaultClaimFit(memReq int64, cpuReq int64, settings *types.Settings) bool {
	return float64(memReq) < float64(settings.DefaultClaim.Memory().Value())*1-settings.Margin &&
		float64(cpuReq) < float64(settings.DefaultClaim.Cpu().MilliValue())*1-settings.Margin
}

func CheckSpec(memReq int64, cpuReq int64, settings *types.Settings) (claimFit bool, spec v1.ResourceList) {
	claimFit = defaultClaimFit(memReq, cpuReq, settings)
	if claimFit {
		spec = settings.DefaultClaim
	} else {
		spec = recommendedSpec(memReq, cpuReq, settings)
	}

	return
}

func CheckRespectMaxNS(spec v1.ResourceList, stats types.ClusterStat, settings *types.Settings) bool {
	return spec.Cpu().MilliValue() < int64(float64(stats.CpuAvailable)*settings.RatioCpuNs) &&
		spec.Memory().Value() < int64(float64(stats.MemAvailable)*settings.RatioMemNs)
}

func Result(stats types.ClusterStat, settings *types.Settings) string {
	var buffer bytes.Buffer

	commitedMem := int64(math.Round(float64(stats.MemAvailable) * settings.OvercommitMem))
	commitedCpu := int64(math.Round(float64(stats.CpuAvailable) * settings.OvercommitCpu))

	passMem := aggregate.ClaimMemSum(stats) < commitedMem
	passCPU := aggregate.ClaimCpuSum(stats) < commitedCpu
	pass := passCPU && passMem

	if pass {
		buffer.WriteString("OK")
	} else {
		buffer.WriteString("NOT FEASIBLE\n")
	}

	if !passMem {
		missingMem := aggregate.ClaimMemSum(stats) - commitedMem
		buffer.WriteString(fmt.Sprintf("Missing %s of Memory\n", ByteToString(missingMem)))
	}

	if !passCPU {
		missingCpu := aggregate.ClaimCpuSum(stats) - commitedCpu
		buffer.WriteString(fmt.Sprintf("Missing %dm of CPU", missingCpu))
	}

	return buffer.String()

}
