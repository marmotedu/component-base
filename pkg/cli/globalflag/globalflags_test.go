// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package globalflag

import (
	"flag"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/pflag"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
)

func TestAddGlobalFlags(t *testing.T) {
	namedFlagSets := &cliflag.NamedFlagSets{}
	nfs := namedFlagSets.FlagSet("global")
	AddGlobalFlags(nfs, "test-cmd")

	actualFlag := []string{}
	nfs.VisitAll(func(flag *pflag.Flag) {
		actualFlag = append(actualFlag, flag.Name)
	})

	// Get all flags from flags.CommandLine, except flag `test.*`.
	wantedFlag := []string{"help"}
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.VisitAll(func(flag *pflag.Flag) {
		if !strings.Contains(flag.Name, "test.") {
			wantedFlag = append(wantedFlag, normalize(flag.Name))
		}
	})
	sort.Strings(wantedFlag)

	if !reflect.DeepEqual(wantedFlag, actualFlag) {
		t.Errorf("[Default]: expected %+v, got %+v", wantedFlag, actualFlag)
	}

	tests := []struct {
		expectedFlag  []string
		matchExpected bool
	}{
		{
			// Happy case
			expectedFlag:  []string{"help"},
			matchExpected: false,
		},
		{
			// Missing flag
			expectedFlag:  []string{"logtostderr", "log-dir"},
			matchExpected: true,
		},
		{
			// Empty flag
			expectedFlag:  []string{},
			matchExpected: true,
		},
		{
			// Invalid flag
			expectedFlag:  []string{"foo"},
			matchExpected: true,
		},
	}

	for i, test := range tests {
		if reflect.DeepEqual(test.expectedFlag, actualFlag) == test.matchExpected {
			t.Errorf("[%d]: expected %+v, got %+v", i, test.expectedFlag, actualFlag)
		}
	}
}
