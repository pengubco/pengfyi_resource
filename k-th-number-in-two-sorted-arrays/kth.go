package kth

// FindKthNumberInTwoSortedArrays returns the k-th smallest number in two
// non empty, sorted arrays, a and b.
func FindKthNumberInTwoSortedArrays(a, b []int, k int) int {
	n, m := len(a), len(b)
	if n > m {
		n, m = m, n
		a, b = b, a
	}

	// Handle special cases for code readability.
	// Case 1. use 0 number from a
	if k <= m && a[0] >= b[k-1] {
		return b[k-1]
	}
	// Case 2. use 0 numbers from b
	if k <= n && b[0] >= a[k-1] {
		return a[k-1]
	}

	// Case 3. use some numbers from a and b.
	low, high := 1, min(k-1, n)
	for low < high {
		// use i numbers in a.
		i := (low + high) / 2
		j := k - i
		if j < 0 || (j < m && a[i-1] > b[j]) {
			high = i - 1
			continue
		}
		if j > m || (i < n && j > 0 && b[j-1] > a[i]) {
			low = i + 1
			continue
		}
		low = i
		break
	}
	ans := a[low-1]
	if k-low-1 >= 0 {
		ans = max(a[low-1], b[k-low-1])
	}
	return ans
}
