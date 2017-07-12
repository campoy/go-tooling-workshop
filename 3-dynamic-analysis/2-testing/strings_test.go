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

package strings_test

import (
	"strings"
	"testing"
)

func TestIndex_Repeat(t *testing.T) {
	if i := strings.Index("hello, world", "hello"); i != 0 {
		t.Errorf(`"hello, world" should contain "hello" at position 0, not %d`, i)
	}
	if i := strings.Index("hello, world", "bye"); i != -1 {
		t.Errorf(`"hello, world" should not contain "bye"`)
	}
	if i := strings.Index("hello, world", "world"); i != 7 {
		t.Errorf(`"hello, world" should contain "world" at position 7, not %d`, i)
	}
}

func TestIndex_Table(t *testing.T) {
	tt := []struct {
		text string
		sub  string
		idx  int
	}{
		{"hello, world", "hello", 0},
		{"hello, world", "bye", -1},
		{"hello, world", "world", 7},
	}
	for _, tc := range tt {
		if idx := strings.Index(tc.text, tc.sub); idx != tc.idx {
			if tc.idx >= 0 {
				t.Errorf("%s should contain %s at position %d, not %d", tc.text, tc.sub, tc.idx, idx)
			} else {
				t.Errorf("%s should not contain %s", tc.text, tc.sub)
			}
		}
	}
}

func TestIndex_Subtest(t *testing.T) {
	tt := []struct {
		name string
		text string
		sub  string
		idx  int
	}{
		{"first character", "hello, world", "hello", 0},
		{"not found", "hello, world", "bye", -1},
		{"last character", "hello, world", "world", 7},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if idx := strings.Index(tc.text, tc.sub); idx != tc.idx {
				if tc.idx >= 0 {
					t.Fatalf("%s should contain %s at position %d, not %d", tc.text, tc.sub, tc.idx, idx)
				}
				t.Fatalf("%s should not contain %s", tc.text, tc.sub)
			}
		})
	}
}
