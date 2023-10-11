package kth

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindKthNumberInTwoSortedArrays_large(t *testing.T) {
	n := 1_000_000
	m := 2_000_000
	a := make([]int, n)
	b := make([]int, m)
	c := make([]int, n+m)
	for i := 0; i < n; i++ {
		a[i] = rand.Int()
		c[i] = a[i]
	}
	for i := 0; i < m; i++ {
		b[i] = rand.Int()
		c[n+i] = b[i]
	}

	sort.Ints(a)
	sort.Ints(b)
	sort.Ints(c)

	for k := 1; k <= n+m; k++ {
		assert.Equal(t, c[k-1], FindKthNumberInTwoSortedArrays(a, b, k))
	}
}

func TestFindKthNumberInTwoSortedArrays_1(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{5, 6, 7, 8}
	for k := 1; k <= 8; k++ {
		assert.Equal(t, k, FindKthNumberInTwoSortedArrays(a, b, k))
	}
}

func TestFindKthNumberInTwoSortedArrays_2(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{1, 2, 3, 4}
	for k := 1; k <= 8; k++ {
		assert.Equal(t, (k+1)/2, FindKthNumberInTwoSortedArrays(a, b, k))
	}
}

func TestFindKthNumberInTwoSortedArrays_3(t *testing.T) {
	a := []int{1, 2, 2, 3, 4}
	b := []int{1, 2, 3, 4, 4, 4, 4}
	for k := 1; k <= 2; k++ {
		assert.Equal(t, 1, FindKthNumberInTwoSortedArrays(a, b, k))
	}
	for k := 3; k <= 5; k++ {
		assert.Equal(t, 2, FindKthNumberInTwoSortedArrays(a, b, k))
	}
	for k := 6; k <= 7; k++ {
		assert.Equal(t, 3, FindKthNumberInTwoSortedArrays(a, b, k))
	}
	for k := 8; k <= 12; k++ {
		assert.Equal(t, 4, FindKthNumberInTwoSortedArrays(a, b, k))
	}
}
