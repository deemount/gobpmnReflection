package utils

// MergeMaps ...
func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

// MergeStringSliceToMap ...
func MergeStringSliceToMap(m map[string][]interface{}, k string, v []interface{}) {
	if m[k] == nil {
		m[k] = make([]interface{}, len(v))
		for i, s := range v {
			m[k][i] = interface{}(s)
		}
	} else {
		m[k] = append(m[k], v...)
	}
}
