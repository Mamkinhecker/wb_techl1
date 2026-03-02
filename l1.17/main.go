package main

func binsearch(arr []int, id int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == id {
			return mid
		}
		if arr[mid] < id {
			left = mid + 1
		} else {
			right = mid - 1
		}

		return -1
	}
}
