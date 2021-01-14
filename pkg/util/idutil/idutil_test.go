// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package idutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUUID(t *testing.T) {
	fmt.Println(GetUUID(""))
}

func TestGetUUID36(t *testing.T) {
	fmt.Println(GetUUID36(""))
}

func TestGetManyUuid(t *testing.T) {
	for i := 0; i < 10000; i++ {
		testID := GetUUID("")
		if len(testID) != 12 {
			t.Errorf("GetUUID failed, expected uuid length 12, got: %d", len(testID))
		}
	}
}

func TestRandString(t *testing.T) {
	str := randString(Alphabet62, 50)
	assert.Equal(t, 50, len(str))
	t.Log(str)

	str = randString(Alphabet62, 255)
	assert.Equal(t, 255, len(str))
	t.Log(str)
}
