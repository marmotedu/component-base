// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sliceutil

// RemoveString remove string from slice if function return true.
func RemoveString(slice []string, remove func(item string) bool) []string {
	for i := 0; i < len(slice); i++ {
		if remove(slice[i]) {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

// FindString return true if target in slice, return false if not.
func FindString(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

// FindInt return true if target in slice, return false if not.
func FindInt(slice []int, target int) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

// FindUint return true if target in slice, return false if not.
func FindUint(slice []uint, target uint) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
