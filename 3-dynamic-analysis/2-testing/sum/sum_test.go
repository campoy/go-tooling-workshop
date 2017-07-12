// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sum_test

import (
	"testing"

	"github.com/campoy/go-tooling-workshop/3-dynamic-analysis/2-testing/sum"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		sum   int
	}{
		{
			name:  "simple test with zero value",
			input: []int{},
			sum:   0,
		},
		{
			name:  "simple test with three value",
			input: []int{1, 2, 3},
			sum:   6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := sum.All(tc.input...), tc.sum; got != want {
				t.Fatalf("unexpected result:\n- want: %d\n-  got: %d\n",
					want, got)
			}
		})
	}
}
