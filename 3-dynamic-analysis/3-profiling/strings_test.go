package strings_test

import (
	"strings"
	"testing"
)

func BenchmarkFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Fields("hello, dear friend")
	}
}

var res []string

func BenchmarkFields_Escape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res = strings.Fields("hello, dear friend")
	}
}
