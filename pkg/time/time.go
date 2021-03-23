// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package time implements a new time with specified time format.
package time

import (
	sqldriver "database/sql/driver"
	"fmt"
	"time"
)

const (
	defaultDateTimeFormat = "2006-01-02 15:04:05"
)

// Time format json time field by myself.
type Time struct {
	time.Time
}

// MarshalJSON on Time format Time field with %Y-%m-%d %H:%M:%S.
func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(defaultDateTimeFormat))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (sqldriver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time.
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// ToTime convert string to Time.
func ToTime(str string) (Time, error) {
	var jt Time
	loc, _ := time.LoadLocation("Local")
	value, err := time.ParseInLocation(defaultDateTimeFormat, str, loc)
	if err != nil {
		return jt, err
	}
	return Time{
		Time: value,
	}, nil
}

// Now returns the current time.
func Now() Time {
	return Time{
		Time: time.Now(),
	}
}
