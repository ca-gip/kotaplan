package utils

import (
	"bytes"
	"fmt"
	"github.com/ca-gip/kotaplan/internal/types"
	unit "github.com/docker/go-units"
	v1 "k8s.io/api/core/v1"
)

func ByteToString(bytes int64) string {
	return unit.BytesSize(float64(bytes))
}

func ResourceListToString(list v1.ResourceList) string {
	return fmt.Sprintf("CPU : %dm	MEM: %s", list.Cpu().MilliValue(), ByteToString(list.Memory().Value()))
}

func SpecFromIntToString(ram int64, cpu int64) string {
	return fmt.Sprintf("CPU : %dm	MEM: %s", cpu, ByteToString(ram))
}

func GenerateClaimManifest(stat types.ClusterStat) string {

	var buffer bytes.Buffer

	for _, ns := range stat.NamespacesStat {
		buffer.WriteString(
			fmt.Sprintf(
				"---\napiVersion: ca-gip.github.com/v1\nkind: ResourceQuotaClaim\nmetadata:\n  name: rqc-%s\n  namespace: %s\nspec:\n  memory: %s\n  cpu: %dm\n",
				ns.Name,
				ns.Name,
				ns.Spec.Memory().String(),
				ns.Spec.Cpu().MilliValue()))

	}

	return buffer.String()
}

func LabelsToString(labels map[string]string) string {
	buffer := new(bytes.Buffer)
	for key, value := range labels {
		buffer.WriteString(
			fmt.Sprintf("%s=%s", key, value))
	}
	return buffer.String()
}
