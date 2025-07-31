package k8sutil

import v2 "k8s.io/api/core/v1"

func ConvertResourceListToStringMap(list v2.ResourceList) map[string]int64 {
	capacity := make(map[string]int64, len(list))
	for k, v := range list {
		capacity[string(k)] = v.Value()
	}
	return capacity
}

func ConvertPVAccessModesToStrings(modes []v2.PersistentVolumeAccessMode) []string {
	modesSlice := make([]string, len(modes))
	for i, mode := range modes {
		modesSlice[i] = string(mode)
	}
	return modesSlice
}

func ConvertResourceStatusesToStringMap(statuses map[v2.ResourceName]v2.ClaimResourceStatus) map[string]string {
	m := make(map[string]string, len(statuses))
	for k, v := range statuses {
		m[string(k)] = string(v)
	}
	return m
}
