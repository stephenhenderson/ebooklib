package utils

// StringSliceEquals compares two string slices and returns true if they
// contain the same contents or are both nil
func StringSliceEquals(arr1, arr2 []string) bool {
	if arr1 == nil && arr2 == nil {
		return true
	}
	if arr1 != nil && arr2 == nil {
		return false
	}
	if arr1 == nil && arr2 != nil {
		return false
	}
	if len(arr1) != len(arr2) {
		return false
	}
	for i, val := range arr1 {
		if val != arr2[i] {
			return false
		}
	}
	return true
}
