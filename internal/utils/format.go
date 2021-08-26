package utils

import (
	"bytes"
	"fmt"
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

func LabelsToString(labels map[string]string) string {
	buffer := new(bytes.Buffer)
	for key, value := range labels {
		buffer.WriteString(
			fmt.Sprintf("%s=%s", key, value))
	}
	return buffer.String()
}
