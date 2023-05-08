package util

// Constants for all supported sorts
const (
	NAME        = "name"
	DESCRIPTION = "description"
)

// IsSupportedSort returns true if the sort is supported
func IsSupportedSort(sort string) bool {
	switch sort {
	case NAME, DESCRIPTION:
		return true
	}
	return false
}
