// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package flag

// OmitEmpty is an interface for flags to report whether their underlying value
// is "empty." If a flag implements OmitEmpty and returns true for a call to Empty(),
// it is assumed that flag may be omitted from the command line.
type OmitEmpty interface {
	Empty() bool
}
