package pool

import "strings"

var (
	structPool = "pool"
)

// IsPool checks if the field is a pool
func IsPool(field string) bool {
	return strings.ToLower(field) == structPool
}
