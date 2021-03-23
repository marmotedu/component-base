// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jsonutil

import (
	"fmt"
	"strings"

	"github.com/marmotedu/component-base/pkg/json"
)

type JSONRawMessage []byte

func (m JSONRawMessage) Find(key string) JSONRawMessage {
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(m, &objmap)
	if err != nil {
		fmt.Printf("Resolve JSON Key failed, find key =%s, err=%s",
			key, err)
		return nil
	}
	return JSONRawMessage(objmap[key])
}

func (m JSONRawMessage) ToList() []JSONRawMessage {
	var lists []json.RawMessage
	err := json.Unmarshal(m, &lists)
	if err != nil {
		fmt.Printf("Resolve JSON List failed, err=%s",
			err)
		return nil
	}
	var res []JSONRawMessage
	for _, v := range lists {
		res = append(res, JSONRawMessage(v))
	}
	return res
}

func (m JSONRawMessage) ToString() string {
	res := strings.ReplaceAll(string(m[:]), "\"", "")
	return res
}
