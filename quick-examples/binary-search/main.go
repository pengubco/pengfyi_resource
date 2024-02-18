package main

import (
	"fmt"
	"sort"
)

func main() {
	nonDecreasing()
	fmt.Println("=======")
	nonIncreasing()
}

func nonDecreasing() {
	//  []int{1, 3, 4, 4, 6, 7, 9, 9}
	// index: 0, 1, 2, 3, 4, 5, 6, 7
	nums := []int{1, 3, 4, 4, 6, 7, 9, 9}

	fmt.Printf("lowerBound(%d) %d\n", 0, lowerBound(nums, 0))   // 0
	fmt.Printf("lowerBound(%d) %d\n", 1, lowerBound(nums, 1))   // 0
	fmt.Printf("lowerBound(%d) %d\n", 9, lowerBound(nums, 9))   // 6
	fmt.Printf("lowerBound(%d) %d\n", 10, lowerBound(nums, 10)) // 8

	fmt.Printf("upperBound(%d) %d\n", 0, upperBound(nums, 0))   // 0
	fmt.Printf("upperBound(%d) %d\n", 1, upperBound(nums, 1))   // 1
	fmt.Printf("upperBound(%d) %d\n", 9, upperBound(nums, 9))   // 8
	fmt.Printf("upperBound(%d) %d\n", 10, upperBound(nums, 10)) // 8
}

func nonIncreasing() {
	//  []int{9, 9, 7, 6, 4, 4, 3, 1}
	// index: 0, 1, 2, 3, 4, 5, 6, 7
	nums := []int{9, 9, 7, 6, 4, 4, 3, 1}

	fmt.Printf("lowerBoundNonIncreasing(%d) %d\n", 0, lowerBoundNonIncreasing(nums, 0))   // 8
	fmt.Printf("lowerBoundNonIncreasing(%d) %d\n", 1, lowerBoundNonIncreasing(nums, 1))   // 7
	fmt.Printf("lowerBoundNonIncreasing(%d) %d\n", 9, lowerBoundNonIncreasing(nums, 9))   // 0
	fmt.Printf("lowerBoundNonIncreasing(%d) %d\n", 10, lowerBoundNonIncreasing(nums, 10)) // 0

	fmt.Printf("upperBoundNonIncreasing(%d) %d\n", 0, upperBoundNonIncreasing(nums, 0))   // 8
	fmt.Printf("upperBoundNonIncreasing(%d) %d\n", 1, upperBoundNonIncreasing(nums, 1))   // 8
	fmt.Printf("upperBoundNonIncreasing(%d) %d\n", 9, upperBoundNonIncreasing(nums, 9))   // 2
	fmt.Printf("upperBoundNonIncreasing(%d) %d\n", 10, upperBoundNonIncreasing(nums, 10)) // 0
}

// Given nums non-decreasing, returns the smallest index i such that
// target <= nums[i]. Return len(nums) if no such nums[i] exists.
func lowerBound(nums []int, target int) int {
	return sort.Search(len(nums), func(i int) bool {
		return target <= nums[i]
	})
}

// Given nums non-decreasing, returns the smallest index i such that
// target < nums[i]. Return len(nums) if no such nums[i] exists.
func upperBound(nums []int, target int) int {
	return sort.Search(len(nums), func(i int) bool {
		return target < nums[i]
	})
}

// Given nums non-increasing, returns the smallest index i such that
// target >= nums[i]. Return len(nums) if no such nums[i] exists.
func lowerBoundNonIncreasing(nums []int, target int) int {
	return sort.Search(len(nums), func(i int) bool {
		return target >= nums[i]
	})
}

// Given nums non-increasing, returns the smallest index i such that
// target > nums[i]. Return len(nums) if no such nums[i] exists.
func upperBoundNonIncreasing(nums []int, target int) int {
	return sort.Search(len(nums), func(i int) bool {
		return target > nums[i]
	})
}
