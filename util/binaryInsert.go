package util

func InsertSorted(array []int64, val int64) []int64 {
	if len(array) == 0 {
		return []int64{val}
	}

	left, right := 0, len(array)-1
	for left <= right {
		mid := (left + right) / 2
		if array[mid] == val {
			return append(array[:mid],
				append([]int64{val}, array[mid:]...)...)
		} else if array[mid] < val {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return append(array[:left], append([]int64{val},
		array[left:]...)...)
}
